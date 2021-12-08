package service

import (
	"HelloWorld/common"
	"HelloWorld/common/logger"
	"HelloWorld/pkg/pb"
	g "HelloWorld/service/grpc"
	"context"

	"github.com/gorilla/mux"
)

// HelloWorldService main struct to hold the microservice server
// Consider change to RunTime Server
type HelloWorldService struct {
	pb.UnimplementedGreeterServer
	configOpts common.ConfigOpts
	grpc       *g.Runtime
	http       *mux.Router
	cancelFunc []func()
}

// NewHelloWorldService create a new server for the microservice protobuffer
func NewHelloWorldService(configOpts common.ConfigOpts) HelloWorldService {

	svc := HelloWorldService{}

	svc.grpc = g.NewServer(configOpts.GrpcOpts)
	pb.RegisterGreeterServer(svc.grpc.GetRuntime(), svc)

	return svc
}


func (svc HelloWorldService) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	r := pb.HelloReply{}
	return &r, nil
}
	

func (svc HelloWorldService) SayHelloAgain(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	r := pb.HelloReply{}
	return &r, nil
}
	

func (svc HelloWorldService) SayHelloClientStream(svr pb.Greeter_SayHelloClientStreamServer) error {
	return nil
}

func (svc HelloWorldService) SayHelloServerStream(in *pb.HelloRequest, svr pb.Greeter_SayHelloServerStreamServer) error {
	return nil
}

func (svc HelloWorldService) SayHelloStream(svr pb.Greeter_SayHelloStreamServer) error {
	return nil
}

// Start ...
func (svc HelloWorldService) Start() {
	err := svc.grpc.Serve()
	if err != nil {
		logger.Log.Errorf("%s", err)
	}
}
