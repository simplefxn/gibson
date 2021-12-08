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
	"context"
	"fmt"
	"log"
	"ripe/common/logger"
	"ripe/pkg/pb"
	"ripe/service/sse"
	"sync"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup
		logger.RegisterLog()
		// Wait channel , we close when no more goroutines running
		waitCh := make(chan struct{})
		evCh := make(chan *pb.RIS_Message)
		sendCh := make(chan *pb.RIS_Message)

		uri := "https://ris-live.ripe.net/v1/stream/?format=sse&client=ripe-client"

		wg.Add(1)
		go func() {
			sse.Notify(uri, evCh)
			wg.Done()
			close(waitCh)
		}()

		// Send goroutine
		wg.Add(1)
		go func() {
			conn, err := grpc.Dial(fmt.Sprintf("%s:%d", opts.GrpcOpts.Host, opts.GrpcOpts.Port), grpc.WithInsecure())
			if err != nil {
				log.Fatalf("did not connect: %s", err)
			}
			defer conn.Close()

			client := pb.NewRipeClient(conn)

			stream, err := client.Stream(context.Background())
			if err != nil {
				log.Fatalf("Error when calling Ping: %s", err)
			}

			for {
				select {
				case ev := <-sendCh:
					err := stream.Send(ev)
					if err != nil {
						logger.Log.Errorf(err.Error())
						break
					}
				case <-waitCh:
					logger.Log.Info("Sender exiting too")
					break
				}
			}
		}()

		// This is our block statement
		for {
			select {
			case ev := <-evCh:
				sendCh <- ev
			case <-waitCh:
				logger.Log.Info("Exiting")
				break
			}
		}
	},
}

func init() {
	runCmd.AddCommand(clientCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clientCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clientCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
