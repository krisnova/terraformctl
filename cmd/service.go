// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/kris-nova/kubicorn/cutil/logger"
	"github.com/kris-nova/terraformctl/service"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"strconv"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run terraformctl as a service",
	Long: fmt.Sprintf(`%s

Use this command to start a gRPC server, and kick of the terraform controller concurrently.`, banner()),
	Run: func(cmd *cobra.Command, args []string) {
		err := RunServer(so)
		if err != nil {
			logger.Critical("Error while running the terraformctl service: %v", err)
		}
	},
}

// ServiceOptions defines the configuration to be used for the RunServer function. This is how
// a user will configure their gRPC server for terraformctl.
type ServiceOptions struct {
	Port int
}

// so is an unexported ServiceOptions variable that will be set and managed by command line input from the user.
var so = &ServiceOptions{}

func init() {
	RootCmd.AddCommand(serverCmd)

	// Port defaults to 4392 because that is the elevation of Mt. Rainier (in meters).
	// which is the next mountain I will be climbing after I write this code.
	serverCmd.Flags().IntVarP(&so.Port, "port", "p", 4392, "Set the port number for the gRPC server to listen on.")
}

// RunServer will start a gRPC server using ServiceOptions to configure the server.
func RunServer(options *ServiceOptions) error {
	port := fmt.Sprintf(":%s", strconv.Itoa(options.Port))

	// gRPC listens over TCP and accepts a TCP listener that can be built from the go standard library.
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("Unable to open TCP socket to listen on: %v", err)
	}
	// Here we initilise a new gRPC server, and register an implementation of the gRPC interface generated from the
	// protobuf file in service/terraformctl.proto.
	server := grpc.NewServer()
	terraformctl.RegisterTerraformCTLAPIServer(server, &terraformctl.TerraformCTLAPIServerImplementation{})
	reflection.Register(server)

	// Now we can run the server.
	logger.Debug("Starting terraformctl gRPC server..")
	if err := server.Serve(listener); err != nil {
		return fmt.Errorf("Failed to start serving gRPC service: %v", err)
	}
	return nil
}
