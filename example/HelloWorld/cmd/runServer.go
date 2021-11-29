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
	"errors"
	"fmt"
	"log"
	"net"

	pb "HelloWorld/pkg/pb"
	g "HelloWorld/server/grpc"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// runServerCmd represents the run command
var runServerCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
			return err
		}
		var opts []grpc.ServerOption
		if tls {
			if certFile == "" {
				return errors.New("No certFile found")
			}
			if keyFile == "" {
				return errors.New("No keyFile found")
			}
			creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
			if err != nil {
				log.Fatalf("Failed to generate credentials %v", err)
				return err
			}
			opts = []grpc.ServerOption{grpc.Creds(creds)}
		}
		grpcServer := grpc.NewServer(opts...)
		helloWorldServer := g.Server{}
		pb.RegisterGreeterServer(grpcServer, helloWorldServer)
		grpcServer.Serve(lis)

		return nil
	},
}

func init() {
	serverCmd.AddCommand(runServerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runServerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runServerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
