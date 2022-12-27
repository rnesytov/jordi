package tui

import (
	"github.com/fullstorydev/grpcurl"
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type TUIGRPCInvocationHandler struct {
}

var _ grpcurl.InvocationEventHandler = (*TUIGRPCInvocationHandler)(nil)

func (h *TUIGRPCInvocationHandler) OnResolveMethod(*desc.MethodDescriptor)        {}
func (h *TUIGRPCInvocationHandler) OnSendHeaders(metadata.MD)                     {}
func (h *TUIGRPCInvocationHandler) OnReceiveHeaders(metadata.MD)                  {}
func (h *TUIGRPCInvocationHandler) OnReceiveResponse(proto.Message)               {}
func (h *TUIGRPCInvocationHandler) OnReceiveTrailers(*status.Status, metadata.MD) {}
