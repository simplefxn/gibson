package grpc

import (
	"HelloWorld/common"
	"HelloWorld/common/logger"
	"crypto/tls"
	"fmt"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

// Runtime ...
type Runtime struct {
	config     common.GrpcOptions
	runtime    *grpc.Server
	grpcOpts   *[]grpc.ServerOption
	cancelFunc []func()
}

// NewServer ...
func NewServer(configOpts common.GrpcOptions) *Runtime {
	r := Runtime{}
	r.config = configOpts
	r.grpcOpts = serverConfig(configOpts)

	r.runtime = grpc.NewServer(*r.grpcOpts...)

	return &r
}

// GetRuntime ...
func (g Runtime) GetRuntime() *grpc.Server {
	return g.runtime
}

// Shutdown ...
func (g Runtime) Shutdown() {
	g.runtime.GracefulStop()
}

// Serve ...
func (g Runtime) Serve() error {
	logger.Log.Infof(" Listening on %s:%d", g.config.Host, g.config.Port)
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", g.config.Host, g.config.Port))
	if err != nil {
		logger.Log.Fatalf("failed to listen: %v", err)
	}

	return errors.WithStack(g.runtime.Serve(lis))
}

// serverConfig ...
func serverConfig(configOpts common.GrpcOptions) *[]grpc.ServerOption {
	grpcOpts := []grpc.ServerOption{}

	// TODO: Wire a flag
	if configOpts.Tracing {
		grpc.EnableTracing = true
	}

	if configOpts.Cert != "" && configOpts.Key != "" {
		config, err := serverConfigTLS(configOpts)
		if err != nil {
			logger.Log.Fatalf("Failed to log gRPC cert/key/ca: %v", err)
		}

		// create the creds server options
		creds := credentials.NewTLS(config)
		grpcOpts = append(grpcOpts, grpc.Creds(creds))
	}

	if configOpts.InitialConnWindowSize != 0 {
		logger.Log.Infof("Setting grpc server initial conn window size to %d", int32(configOpts.InitialConnWindowSize))
		grpcOpts = append(grpcOpts, grpc.InitialConnWindowSize(int32(configOpts.InitialConnWindowSize)))
	}

	if configOpts.InitialWindowSize != 0 {
		logger.Log.Infof("Setting grpc server initial window size to %d", int32(configOpts.InitialWindowSize))
		grpcOpts = append(grpcOpts, grpc.InitialWindowSize(int32(configOpts.InitialWindowSize)))
	}

	ep := keepalive.EnforcementPolicy{
		MinTime:             configOpts.KeepAliveEnforcementPolicyMinTime,
		PermitWithoutStream: configOpts.KeepAliveEnforcementPolicyPermitWithoutStream,
	}
	grpcOpts = append(grpcOpts, grpc.KeepaliveEnforcementPolicy(ep))

	if configOpts.MaxConnectionAge != 0 {
		ka := keepalive.ServerParameters{
			MaxConnectionAge: configOpts.MaxConnectionAge,
		}
		if configOpts.MaxConnectionAgeGrace != 0 {
			ka.MaxConnectionAgeGrace = configOpts.MaxConnectionAgeGrace
		}
		grpcOpts = append(grpcOpts, grpc.KeepaliveParams(ka))
	}
	var streamInterceptors []grpc.StreamServerInterceptor
	var unaryInterceptor []grpc.UnaryServerInterceptor

	streamInterceptors = append(streamInterceptors, grpc_ctxtags.StreamServerInterceptor())
	unaryInterceptor = append(unaryInterceptor, grpc_ctxtags.UnaryServerInterceptor())

	if configOpts.Zap {
		grpc_zap.ReplaceGrpcLoggerV2(logger.GetLogger())
		opts := []grpc_zap.Option{}

		streamInterceptors = append(streamInterceptors, grpc_zap.StreamServerInterceptor(logger.GetLogger(), opts...))
		unaryInterceptor = append(unaryInterceptor, grpc_zap.UnaryServerInterceptor(logger.GetLogger(), opts...))
	}

	if configOpts.Prometheus {
		streamInterceptors = append(streamInterceptors, grpc_prometheus.StreamServerInterceptor)
		unaryInterceptor = append(unaryInterceptor, grpc_prometheus.UnaryServerInterceptor)
	}

	if configOpts.Recovery {
		streamInterceptors = append(streamInterceptors, grpc_recovery.StreamServerInterceptor())
		unaryInterceptor = append(unaryInterceptor, grpc_recovery.UnaryServerInterceptor())
	}

	if configOpts.Opentracing {
		streamInterceptors = append(streamInterceptors, grpc_opentracing.StreamServerInterceptor())
		unaryInterceptor = append(unaryInterceptor, grpc_opentracing.UnaryServerInterceptor())
	}

	grpcOpts = append(grpcOpts, grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(streamInterceptors...)))

	grpcOpts = append(grpcOpts, grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(unaryInterceptor...)))

	//gOpts = append(configOpts, interceptors()...)
	return &grpcOpts
}

// serverConfigTLS ...
func serverConfigTLS(grpcOpts common.GrpcOptions) (*tls.Config, error) {
	//func ServerConfig(cert, key, ca, crl, serverCA string, minTLSVersion uint16) (*tls.Config, error) {
	config := newTLSConfig(tls.VersionTLS12)

	var certificates *[]tls.Certificate
	var err error

	if grpcOpts.CA != "" {
		certificates, err = combineAndLoadTLSCertificates(grpcOpts.ServerCA, grpcOpts.Cert, grpcOpts.Key)
	} else {
		certificates, err = loadTLSCertificate(grpcOpts.Cert, grpcOpts.Key)
	}

	if err != nil {
		return nil, err
	}
	config.Certificates = *certificates

	// if specified, load ca to validate client,
	// and enforce clients present valid certs.
	if grpcOpts.CA != "" {
		certificatePool, err := loadx509CertPool(grpcOpts.CA)

		if err != nil {
			return nil, err
		}

		config.ClientCAs = certificatePool
		config.ClientAuth = tls.RequireAndVerifyClientCert
	}

	if grpcOpts.CRL != "" {
		crlFunc, err := verifyPeerCertificateAgainstCRL(grpcOpts.CRL)
		if err != nil {
			return nil, err
		}
		config.VerifyPeerCertificate = crlFunc
	}

	return config, nil
}
