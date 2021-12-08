/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"
	"time"

	common "RipeIngest/common"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var cfgFile string
var opts common.ConfigOpts

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "RipeIngest",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&opts.GrpcOpts.Host, "grpc.host", "localhost", "host to connect or bind the socket")

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.RipeIngest.yaml)")

	// GRPCPort is the port to listen on for gRPC. If not set or zero, don't listen.
	rootCmd.PersistentFlags().IntVar(&opts.GrpcOpts.Port, "grpc.port", 8123, "Port to listen on for gRPC calls")

	// GRPCCert is the cert to use if TLS is enabled
	rootCmd.PersistentFlags().StringVar(&opts.GrpcOpts.Cert, "grpc.cert", "", "server certificate to use for gRPC connections, requires grpc_key, enables TLS")

	// GRPCKey is the key to use if TLS is enabled
	rootCmd.PersistentFlags().StringVar(&opts.GrpcOpts.Key, "grpc.key", "", "server private key to use for gRPC connections, requires grpc_cert, enables TLS")

	// GRPCCA is the CA to use if TLS is enabled
	rootCmd.PersistentFlags().StringVar(&opts.GrpcOpts.CA, "grpc.ca", "", "server CA to use for gRPC connections, requires TLS, and enforces client certificate check")

	// GRPCCRL is the CRL (Certificate Revocation List) to use if TLS is enabled
	rootCmd.PersistentFlags().StringVar(&opts.GrpcOpts.CRL, "grpc_crl", "", "path to a certificate revocation list in PEM format, client certificates will be further verified against this file during TLS handshake")

	// GRPCServerCA if specified will combine server cert and server CA
	rootCmd.PersistentFlags().StringVar(&opts.GrpcOpts.ServerCA, "grpc.server.ca", "", "path to server CA in PEM format, which will be combine with server cert, return full certificate chain to clients")

	// GRPCMaxConnectionAge is the maximum age of a client connection, before GoAway is sent.
	// This is useful for L4 loadbalancing to ensure rebalancing after scaling.
	rootCmd.PersistentFlags().DurationVar(&opts.GrpcOpts.MaxConnectionAge, "grpc.max_connection_age", time.Duration(0), "Maximum age of a client connection before GoAway is sent.")

	// GRPCMaxConnectionAgeGrace is an additional grace period after GRPCMaxConnectionAge, after which
	// connections are forcibly closed.
	rootCmd.PersistentFlags().DurationVar(&opts.GrpcOpts.MaxConnectionAgeGrace, "grpc.max_connection_age_grace", time.Duration(0), "Additional grace period after grpc_max_connection_age, after which connections are forcibly closed.")

	// GRPCInitialConnWindowSize ServerOption that sets window size for a connection.
	// The lower bound for window size is 64K and any value smaller than that will be ignored.
	rootCmd.PersistentFlags().IntVar(&opts.GrpcOpts.InitialConnWindowSize, "grpc.server_initial_conn_window_size", 0, "gRPC server initial connection window size")

	// GRPCInitialWindowSize ServerOption that sets window size for stream.
	// The lower bound for window size is 64K and any value smaller than that will be ignored.
	rootCmd.PersistentFlags().IntVar(&opts.GrpcOpts.InitialWindowSize, "grpc.server_initial_window_size", 0, "gRPC server initial window size")

	// EnforcementPolicy MinTime that sets the keepalive enforcement policy on the server.
	// This is the minimum amount of time a client should wait before sending a keepalive ping.
	rootCmd.PersistentFlags().DurationVar(&opts.GrpcOpts.KeepAliveEnforcementPolicyMinTime, "grpc.server_keepalive_enforcement_policy_min_time", 10*time.Second, "gRPC server minimum keepalive time")

	// EnforcementPolicy PermitWithoutStream - If true, server allows keepalive pings
	// even when there are no active streams (RPCs). If false, and client sends ping when
	// there are no active streams, server will send GOAWAY and close the connection.
	rootCmd.PersistentFlags().BoolVar(&opts.GrpcOpts.KeepAliveEnforcementPolicyPermitWithoutStream, "grpc.server_keepalive_enforcement_policy_permit_without_stream", false, "gRPC server permit client keepalive pings even when there are no active streams (RPCs)")

	rootCmd.PersistentFlags().BoolVar(&opts.EnableHTTPServer, "enable_http_server", true, "enable the HTTP server side of the microservice")
	rootCmd.PersistentFlags().BoolVar(&opts.GrpcOpts.Tracing, "grpc.enable_tracing", true, "enable grpc tracing")

	rootCmd.PersistentFlags().BoolVar(&opts.GrpcOpts.Prometheus, "grpc.prometheus", true, "enable grpc prometheus middleware")
	rootCmd.PersistentFlags().BoolVar(&opts.GrpcOpts.Recovery, "grpc.recovery", true, "enable grpc recovery middleware")
	rootCmd.PersistentFlags().BoolVar(&opts.GrpcOpts.Opentracing, "grpc.opentracing", true, "enable grpc opentracing middleware")
	rootCmd.PersistentFlags().BoolVar(&opts.GrpcOpts.Zap, "grpc.zap", true, "enable grpc zap logging middleware")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".RipeIngest" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".RipeIngest")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
