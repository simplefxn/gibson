package grpc

import (
	pb "HelloWorld/pkg/pb"
	"context"
)

type Server struct {
	pb.UnimplementedGreeterServer
}

func (server Server) SayHello(context.Context, *pb.HelloRequest) (*pb.HelloReply, error) {
	reply := pb.HelloReply{}
	return &reply, nil
}

func (server Server) SayHelloAgain(context.Context, *pb.HelloRequest) (*pb.HelloReply, error) {
	reply := pb.HelloReply{}
	return &reply, nil
}
