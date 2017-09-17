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
	"github.com/kris-nova/kubicorn/cutil/namer"
	"github.com/kris-nova/terraformctl/parser"
	"github.com/kris-nova/terraformctl/service"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"os"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set a new configuration definition in terraformctl",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		seto.RootOptions = *ro
		err := RunSet(seto)
		if err != nil {
			logger.Critical(err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	},
}

func init() {
	RootCmd.AddCommand(setCmd)
	setCmd.Flags().StringVarP(&seto.Filename, "filename", "f", ".", "The filename to use to look for a Terraform configuration in. Accepts files or directories.")
	setCmd.Flags().StringVarP(&seto.Name, "name", "n", namer.RandomName(), "The unique identifier name of the configuration. If no name is provided, terraformctl will generate one for you.")

}

type SetOptions struct {
	RootOptions
	Filename string
	Name     string
}

var seto = &SetOptions{}

func RunSet(options *SetOptions) error {
	logger.Info("Using name [%s]", options.Name)
	client := terraformctl.NewTerraformCtlClient(options.Hostname, options.Port)
	err := client.Connect()
	if err != nil {
		return fmt.Errorf("Unable to connect to gRPC server [%s:%d] with error: %v", options.Hostname, options.Port, err)
	}
	tf, err := parser.NewTerraformConfigurationFromPath(options.Name, options.Filename)
	if err != nil {
		return fmt.Errorf("Unable to read terraform configuration with error: %v", err)
	}
	tfString, err := tf.String()
	if err != nil {
		return fmt.Errorf("Unable to convert terraform configuration to string with error: %v", err)
	}
	_, err = client.Client.Set(context.TODO(), &terraformctl.SetRequest{
		Name:   options.Name,
		Config: tfString,
	})
	if err != nil {
		return fmt.Errorf("Unable to set terraform configuration on server with error: %v", err)
	}
	logger.Info("Successfully set configuration [%s]", options.Name)
	return nil
}
