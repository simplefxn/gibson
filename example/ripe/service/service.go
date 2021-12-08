package service

import (
	"io"
	"ripe/common"
	"ripe/common/logger"
	"ripe/pkg/pb"
	g "ripe/service/grpc"

	"github.com/gorilla/mux"
)

// ripeService main struct to hold the microservice server
// Consider change to RunTime Server
type ripeService struct {
	pb.UnimplementedRipeServer
	configOpts common.ConfigOpts
	grpc       *g.Runtime
	http       *mux.Router
	cancelFunc []func()
}

// NewripeService create a new server for the microservice protobuffer
func NewripeService(configOpts common.ConfigOpts) ripeService {

	svc := ripeService{}

	svc.grpc = g.NewServer(configOpts.GrpcOpts)
	pb.RegisterRipeServer(svc.grpc.GetRuntime(), svc)

	return svc
}

// Stream ...
func (svc ripeService) Stream(srv pb.Ripe_StreamServer) error {
	for {
		req, err := srv.Recv()
		if err == io.EOF {
			// Close the connection and return the response to the client
			return srv.SendAndClose(&pb.Empty{})
		}

		// Handle any possible errors while streaming requests
		if err != nil {
			logger.Log.Fatalf("Error when reading client request stream: %v", err)
		}

		logger.Log.Infof("%v", req)
	}
}

// Start ...
func (svc ripeService) Start() {
	err := svc.grpc.Serve()
	if err != nil {
		logger.Log.Errorf("%s", err)
	}
}
