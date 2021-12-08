package service

import (
	"RipeIngest/common"
	"RipeIngest/common/logger"
	"RipeIngest/pkg/pb"
	g "RipeIngest/service/grpc"
	"io"

	"github.com/gorilla/mux"
)

// RipeIngestService main struct to hold the microservice server
// Consider change to RunTime Server
type RipeIngestService struct {
	pb.UnimplementedRipeServer
	configOpts common.ConfigOpts
	grpc       *g.Runtime
	http       *mux.Router
	cancelFunc []func()
}

// NewRipeIngestService create a new server for the microservice protobuffer
func NewRipeIngestService(configOpts common.ConfigOpts) RipeIngestService {

	svc := RipeIngestService{}

	svc.grpc = g.NewServer(configOpts.GrpcOpts)
	pb.RegisterRipeServer(svc.grpc.GetRuntime(), svc)

	return svc
}

// Stream ...
func (svc RipeIngestService) Stream(srv pb.Ripe_StreamServer) error {
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
func (svc RipeIngestService) Start() {
	err := svc.grpc.Serve()
	if err != nil {
		logger.Log.Errorf("%s", err)
	}
}
