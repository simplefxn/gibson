package common

import "time"

// GrpcOptions ...
type GrpcOptions struct {
	Host                                          string
	Port                                          int
	Cert                                          string
	Key                                           string
	CA                                            string
	CRL                                           string
	ServerCA                                      string
	MaxConnectionAge                              time.Duration
	MaxConnectionAgeGrace                         time.Duration
	InitialConnWindowSize                         int
	InitialWindowSize                             int
	KeepAliveEnforcementPolicyMinTime             time.Duration
	KeepAliveEnforcementPolicyPermitWithoutStream bool
	Tracing                                       bool
	Prometheus                                    bool
	Opentracing                                   bool
	Recovery                                      bool
	Zap                                           bool
}

// ConfigOpts ...
type ConfigOpts struct {
	GrpcOpts         GrpcOptions
	EnableHTTPServer bool
}
