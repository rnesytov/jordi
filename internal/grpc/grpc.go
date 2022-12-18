package grpc

import (
	"bytes"
	"context"
	"crypto/tls"
	"time"

	"github.com/fullstorydev/grpcurl"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/grpcreflect"
	"github.com/profx5/jordi/internal/version"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	reflectpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
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
	cc, err := grpcurl.BlockingDial(ctx, network, target, creds, opts...)
	if err != nil {
		return nil, err
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

func (g *GRPCWrapper) GetInOutDescription(method string) (string, string, error) {
	dsc, err := g.descSource.FindSymbol(method)
	if err != nil {
		return "", "", err
	}
	methodDsc, ok := dsc.(*desc.MethodDescriptor)
	if !ok {
		panic("not a method")
	}
	in, err := grpcurl.GetDescriptorText(methodDsc.GetInputType(), nil)
	if err != nil {
		panic(err)
	}
	out, err := grpcurl.GetDescriptorText(methodDsc.GetOutputType(), nil)
	if err != nil {
		panic(err)
	}
	return in, out, nil
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
