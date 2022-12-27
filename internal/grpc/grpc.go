package grpc

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/fullstorydev/grpcurl"
	"github.com/golang/protobuf/jsonpb"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/grpcreflect"
	"github.com/profx5/jordi/internal/version"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	reflectpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type (
	Opts struct {
		ConnectTimeout time.Duration
		KeepaliveTime  time.Duration
		MaxMsgSize     int
		Insecure       bool
	}
	GRPCWrapper struct {
		cc         *grpc.ClientConn
		refClient  *grpcreflect.Client
		descSource grpcurl.DescriptorSource
	}
)

func DefaultOpts() Opts {
	return Opts{
		ConnectTimeout: 10 * time.Second,
	}
}

func New(ctx context.Context, target string, connOpts Opts) (*GRPCWrapper, error) {
	var opts []grpc.DialOption
	if connOpts.KeepaliveTime > 0 {
		opts = append(opts, grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:    connOpts.KeepaliveTime,
			Timeout: connOpts.KeepaliveTime,
		}))
	}
	if connOpts.MaxMsgSize > 0 {
		opts = append(opts, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(connOpts.MaxMsgSize)))
	}
	var creds credentials.TransportCredentials
	if connOpts.Insecure {
		creds = insecure.NewCredentials()
	} else {
		creds = credentials.NewTLS(&tls.Config{})
	}

	userAgent := "jordi/" + version.Version
	opts = append(opts, grpc.WithUserAgent(userAgent))

	network := "tcp"

	type res struct {
		cc  *grpc.ClientConn
		err error
	}
	ch := make(chan res)
	go func() {
		cc, err := grpcurl.BlockingDial(ctx, network, target, creds, opts...)
		if err != nil {
			ch <- res{err: err}
			return
		}
		ch <- res{cc: cc}
	}()
	var cc *grpc.ClientConn
	timeoutCtx, cancel := context.WithTimeout(ctx, connOpts.ConnectTimeout)
	defer cancel()
	select {
	case <-timeoutCtx.Done():
		return nil, fmt.Errorf("connection timed out")
	case r := <-ch:
		if r.err != nil {
			return nil, r.err
		}
		cc = r.cc
	}
	refClient := grpcreflect.NewClientV1Alpha(ctx, reflectpb.NewServerReflectionClient(cc))
	descSource := grpcurl.DescriptorSourceFromServer(ctx, refClient)
	return &GRPCWrapper{
		cc:         cc,
		refClient:  refClient,
		descSource: descSource,
	}, nil
}

func (g *GRPCWrapper) ListServices() ([]string, error) {
	return grpcurl.ListServices(g.descSource)
}

func (g *GRPCWrapper) ListMethods(service string) ([]string, error) {
	return grpcurl.ListMethods(g.descSource, service)
}

type MessageWrapper struct {
	Msg protoreflect.Message
}

func (mw MessageWrapper) ProtoReflect() protoreflect.Message {
	return mw.Msg
}

func (g *GRPCWrapper) GetInOutDescription(method string) (string, string, error) {
	dsc, err := g.descSource.FindSymbol(method)
	if err != nil {
		return "", "", err
	}
	methodDsc, ok := dsc.(*desc.MethodDescriptor)
	if !ok {
		panic("not a method")
	}
	inType := methodDsc.GetInputType()
	in, err := grpcurl.GetDescriptorText(inType, g.descSource)
	if err != nil {
		return "", "", err
	}
	protoMsg := grpcurl.MakeTemplate(inType)
	example, err := (&jsonpb.Marshaler{
		EmitDefaults: true,
		Indent:       "  ",
	}).MarshalToString(protoMsg)
	if err != nil {
		example = "{}"
	}

	return in, example, nil
}

func (g *GRPCWrapper) Invoke(method string, request string) (string, error) {
	options := grpcurl.FormatOptions{
		EmitJSONDefaultFields: false,
		IncludeTextSeparator:  false,
		AllowUnknownFields:    false,
	}
	rf, formatter, err := grpcurl.RequestParserAndFormatter(grpcurl.FormatJSON, g.descSource, bytes.NewReader([]byte(request)), options)
	if err != nil {
		return "", err
	}
	buf := bytes.NewBuffer(nil)
	h := &grpcurl.DefaultEventHandler{
		Out:            buf,
		Formatter:      formatter,
		VerbosityLevel: 1,
	}

	err = grpcurl.InvokeRPC(context.Background(), g.descSource, g.cc, method, []string{}, h, rf.Next)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (g *GRPCWrapper) Close() {
	g.refClient.Reset()
	g.cc.Close()
}
