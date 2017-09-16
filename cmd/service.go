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
	"github.com/kris-nova/terraformctl/controller"
	"github.com/kris-nova/terraformctl/service"
	"github.com/kris-nova/terraformctl/storage"
	"github.com/kris-nova/terraformctl/storage/blobPersist"
	"github.com/kris-nova/terraformctl/storage/memoryCache"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
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
	Port                 int
	CacheSelection       string
	PersistenceSelection string
}

// so is an unexported ServiceOptions variable that will be set and managed by command line input from the user.
var so = &ServiceOptions{}

func init() {
	RootCmd.AddCommand(serverCmd)

	// Port defaults to 4392 because that is the elevation of Mt. Rainier (in meters).
	// which is the next mountain I will be climbing after I write this code.
	serverCmd.Flags().IntVarP(&so.Port, "port", "p", 4392, "Set the port number for the gRPC server to listen on.")

	// CacheSelection defines which caching implementation to use.
	serverCmd.Flags().StringVarP(&so.CacheSelection, "cache", "c", "memory", "Set the cache selection string to use. Currently supports [memory].")

	// PersistenceSelection defines which caching implementation to use.
	serverCmd.Flags().StringVarP(&so.PersistenceSelection, "persistence", "P", "blob", "Set the persistence selection string to use. Currently supports [blob].")

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

	// Register caching and persistence layers
	cacher, err := getCacher(options.CacheSelection)
	if err != nil {
		return fmt.Errorf("Unable to get cacher: %v", err)
	}
	persistence, err := getPersistent(options.PersistenceSelection)
	if err != nil {
		return fmt.Errorf("Unable to get persistence: %v", err)
	}
	loop := controller.NewTerraformControlLoop(&controller.TerraformControlLoopOptions{}, cacher, persistence)
	errch := loop.Run()

	// Now we can run the server.
	logger.Debug("Starting terraformctl gRPC server..")
	if err := server.Serve(listener); err != nil {
		return fmt.Errorf("Failed to start serving gRPC service: %v", err)
	}

	// Look for errors and exit
	for {
		err := <-errch
		if err != nil {
			logger.Warning("Error from control loop: %v", err)
		}
		if err == nil {
			logger.Info("Exiting terraformctl control loop..")
			os.Exit(1)
		}
	}
	return nil
}

func getCacher(selection string) (storage.TerraformCtlCacher, error) {
	logger.Debug("Using cacher [%s]", selection)
	switch selection {
	case "memory":
		return memoryCache.NewMemoryCache(), nil
	default:
		return nil, fmt.Errorf("Invalid cache selection [%s]", selection)
	}
}

func getPersistent(selection string) (storage.TerraformCtlPersister, error) {
	logger.Debug("Using persistence [%s]", selection)
	switch selection {
	case "blob":
		return blobPersist.NewBlobPersist(), nil
	default:
		return nil, fmt.Errorf("Invalid persitency selection [%s]", selection)
	}
}
