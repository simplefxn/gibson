// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// RipeClient is the client API for Ripe service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RipeClient interface {
	Stream(ctx context.Context, opts ...grpc.CallOption) (Ripe_StreamClient, error)
}

type ripeClient struct {
	cc grpc.ClientConnInterface
}

func NewRipeClient(cc grpc.ClientConnInterface) RipeClient {
	return &ripeClient{cc}
}

func (c *ripeClient) Stream(ctx context.Context, opts ...grpc.CallOption) (Ripe_StreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &Ripe_ServiceDesc.Streams[0], "/RipeIngest.Ripe/Stream", opts...)
	if err != nil {
		return nil, err
	}
	x := &ripeStreamClient{stream}
	return x, nil
}

type Ripe_StreamClient interface {
	Send(*RIS_Message) error
	CloseAndRecv() (*Empty, error)
	grpc.ClientStream
}

type ripeStreamClient struct {
	grpc.ClientStream
}

func (x *ripeStreamClient) Send(m *RIS_Message) error {
	return x.ClientStream.SendMsg(m)
}

func (x *ripeStreamClient) CloseAndRecv() (*Empty, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Empty)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// RipeServer is the server API for Ripe service.
// All implementations must embed UnimplementedRipeServer
// for forward compatibility
type RipeServer interface {
	Stream(Ripe_StreamServer) error
	mustEmbedUnimplementedRipeServer()
}

// UnimplementedRipeServer must be embedded to have forward compatible implementations.
type UnimplementedRipeServer struct {
}

func (UnimplementedRipeServer) Stream(Ripe_StreamServer) error {
	return status.Errorf(codes.Unimplemented, "method Stream not implemented")
}
func (UnimplementedRipeServer) mustEmbedUnimplementedRipeServer() {}

// UnsafeRipeServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RipeServer will
// result in compilation errors.
type UnsafeRipeServer interface {
	mustEmbedUnimplementedRipeServer()
}

func RegisterRipeServer(s grpc.ServiceRegistrar, srv RipeServer) {
	s.RegisterService(&Ripe_ServiceDesc, srv)
}

func _Ripe_Stream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(RipeServer).Stream(&ripeStreamServer{stream})
}

type Ripe_StreamServer interface {
	SendAndClose(*Empty) error
	Recv() (*RIS_Message, error)
	grpc.ServerStream
}

type ripeStreamServer struct {
	grpc.ServerStream
}

func (x *ripeStreamServer) SendAndClose(m *Empty) error {
	return x.ServerStream.SendMsg(m)
}

func (x *ripeStreamServer) Recv() (*RIS_Message, error) {
	m := new(RIS_Message)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Ripe_ServiceDesc is the grpc.ServiceDesc for Ripe service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Ripe_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "RipeIngest.Ripe",
	HandlerType: (*RipeServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Stream",
			Handler:       _Ripe_Stream_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "ripe.proto",
}
