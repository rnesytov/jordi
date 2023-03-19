package grpc

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/fullstorydev/grpcurl"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/grpcreflect"
	"github.com/profx5/jordi/internal/version"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	reflectpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type (
	Opts struct {
		ConnectTimeout time.Duration
		KeepaliveTime  time.Duration
		MaxMsgSize     int
		Insecure       bool
	}
	Wrapper struct {
		cc         *grpc.ClientConn
		refClient  *grpcreflect.Client
		descSource grpcurl.DescriptorSource
		reqCancel  func()
		Target     string
	}
	TypeAndError[T any] struct {
		Result T
		Err    error
	}
	InDesc struct {
		Desc    string
		Example string
		Err     error
	}
	gRPCEventType int
	Event         struct {
		Type    gRPCEventType
		Payload interface{}
		Err     error
	}
	gRPCEventHandler struct {
		c chan<- Event
	}
)

const (
	GRPCEventTypeUnknown gRPCEventType = iota
	MethodResolved       gRPCEventType = iota
	HeadersSent          gRPCEventType = iota
	HeadersReceived      gRPCEventType = iota
	ResponseReceived     gRPCEventType = iota
	ReceivedTrailers     gRPCEventType = iota
	EventError           gRPCEventType = iota
)

var (
	ProtoJSONMarshaler = &jsonpb.Marshaler{
		EmitDefaults: true,
		Indent:       "  ",
	}
)

func DefaultOpts() Opts {
	return Opts{
		ConnectTimeout: 10 * time.Second, // TODO: make this configurable
	}
}

func New(ctx context.Context, target string, connOpts Opts) (*Wrapper, error) {
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
		clientConn *grpc.ClientConn
		err        error
	}
	resultChan := make(chan res)
	go func() {
		cc, err := grpcurl.BlockingDial(ctx, network, target, creds, opts...)
		if err != nil {
			resultChan <- res{err: err}
			return
		}
		resultChan <- res{clientConn: cc}
	}()
	var clientConn *grpc.ClientConn
	timeoutCtx, cancel := context.WithTimeout(ctx, connOpts.ConnectTimeout)
	defer cancel()
	select {
	case <-timeoutCtx.Done():
		return nil, fmt.Errorf("connection timed out")
	case r := <-resultChan:
		if r.err != nil {
			return nil, r.err
		}
		clientConn = r.clientConn
	}
	refClient := grpcreflect.NewClientV1Alpha(ctx, reflectpb.NewServerReflectionClient(clientConn))
	descSource := grpcurl.DescriptorSourceFromServer(ctx, refClient)
	return &Wrapper{
		cc:         clientConn,
		refClient:  refClient,
		descSource: descSource,
		reqCancel:  nil,
		Target:     target,
	}, nil
}

func (g *Wrapper) ListServices() <-chan TypeAndError[[]string] {
	resultChan := make(chan TypeAndError[[]string])
	go func() {
		defer close(resultChan)
		services, err := grpcurl.ListServices(g.descSource)
		resultChan <- TypeAndError[[]string]{Result: services, Err: err}
	}()
	return resultChan
}

func (g *Wrapper) ListMethods(service string) <-chan TypeAndError[[]string] {
	resultChan := make(chan TypeAndError[[]string])
	go func() {
		defer close(resultChan)
		methods, err := grpcurl.ListMethods(g.descSource, service)
		resultChan <- TypeAndError[[]string]{Result: methods, Err: err}
	}()
	return resultChan
}

type MessageWrapper struct {
	Msg protoreflect.Message
}

func (mw MessageWrapper) ProtoReflect() protoreflect.Message {
	return mw.Msg
}

func (g *Wrapper) getInDescription(method string) (string, string, error) {
	dsc, err := g.descSource.FindSymbol(method)
	if err != nil {
		return "", "", err
	}
	methodDsc, ok := dsc.(*desc.MethodDescriptor)
	if !ok {
		return "", "", fmt.Errorf("not a method")
	}
	inType := methodDsc.GetInputType()
	inDescText, err := grpcurl.GetDescriptorText(inType, g.descSource)
	if err != nil {
		return "", "", err
	}
	protoMsg := grpcurl.MakeTemplate(inType)
	example, err := ProtoJSONMarshaler.MarshalToString(protoMsg)
	if err != nil {
		example = "{}"
	}

	return inDescText, example, nil
}

func (g *Wrapper) GetInputDescription(method string) <-chan InDesc {
	resultChan := make(chan InDesc)
	go func() {
		defer close(resultChan)
		desc, example, err := g.getInDescription(method)
		resultChan <- InDesc{Desc: desc, Example: example, Err: err}
	}()
	return resultChan
}

func (g *Wrapper) Invoke(method string, request string) (<-chan Event, error) {
	options := grpcurl.FormatOptions{
		EmitJSONDefaultFields: false,
		IncludeTextSeparator:  false,
		AllowUnknownFields:    false,
	}
	requestFormatter, _, err := grpcurl.RequestParserAndFormatter(grpcurl.FormatJSON, g.descSource, bytes.NewReader([]byte(request)), options)
	if err != nil {
		return nil, err
	}
	resultChan := make(chan Event, 10)

	ctx, cancel := context.WithCancel(context.Background())
	g.reqCancel = cancel
	h := &gRPCEventHandler{c: resultChan}
	go func() {
		err := grpcurl.InvokeRPC(ctx, g.descSource, g.cc, method, []string{}, h, requestFormatter.Next)
		if err != nil {
			resultChan <- Event{Type: EventError, Err: err}
			cancel()
			close(resultChan)
		}
	}()
	return resultChan, nil
}

func (g *Wrapper) CancelInvoke() {
	if g.reqCancel != nil {
		g.reqCancel()
		g.reqCancel = nil
	}
}

func (g *Wrapper) Close() {
	g.refClient.Reset()
	g.cc.Close()
}

func (h *gRPCEventHandler) OnResolveMethod(_ *desc.MethodDescriptor) {
	h.c <- Event{Type: MethodResolved}
}
func (h *gRPCEventHandler) OnSendHeaders(metadata.MD) {
	h.c <- Event{Type: HeadersSent}
}
func (h *gRPCEventHandler) OnReceiveHeaders(metadata.MD) {
	h.c <- Event{Type: HeadersReceived}
}
func (h *gRPCEventHandler) OnReceiveResponse(m proto.Message) {
	responseJSON, err := ProtoJSONMarshaler.MarshalToString(m)

	h.c <- Event{Type: ResponseReceived, Payload: responseJSON, Err: err}
}
func (h *gRPCEventHandler) OnReceiveTrailers(s *status.Status, _ metadata.MD) {
	h.c <- Event{Type: ReceivedTrailers, Payload: s.Code().String()}
	close(h.c)
}
