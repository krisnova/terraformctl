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
	"github.com/kris-nova/terraformctl/parser"
	"github.com/kris-nova/terraformctl/service"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"os"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Will get a configuration from the gRPC server.",
	Long:  `This command will get a configuration from the gRPC server and echo it to STDOUT.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := RunGet(geto)
		if err != nil {
			logger.Critical(err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	},
}

type GetOptions struct {
	RootOptions
	Name string
}

var geto = &GetOptions{}

func init() {
	RootCmd.AddCommand(getCmd)
	getCmd.Flags().StringVarP(&geto.Name, "name", "n", "", "The name to use while attempting to look up a known configuration. This is a required flag.")
}

func RunGet(options *GetOptions) error {
	if options.Name == "" {
		return fmt.Errorf("Empty value for the requried parameter 'name'")
	}
	client := terraformctl.NewTerraformCtlClient(options.Hostname, options.Port)
	err := client.Connect()
	if err != nil {
		return fmt.Errorf("Unable to connect to gRPC server [%s:%d] with error: %v", options.Hostname, options.Port, err)
	}
	response, err := client.Client.Get(context.TODO(), &terraformctl.GetRequest{
		Name: options.Name,
	}, nil)
	if err != nil {
		return fmt.Errorf("Unable to get configuration [%s] from gRPC server with error: %v", options.Name, err)
	}
	tfconfig, err := parser.NewTerraformConfigurationFromEncodedString(response.String())
	if err != nil {
		return fmt.Errorf("Unable to create Terraform configuration from encoded string with error: %v", err)
	}
	printStr, err := tfconfig.String()
	if err != nil {
		return fmt.Errorf("Unable to generate string from Terraform configuration: %v", err)
	}
	fmt.Printf("%s\n", printStr)
	return nil
}
