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
	"github.com/kris-nova/terraformctl/controller"
	"github.com/kris-nova/terraformctl/parser"
	"github.com/spf13/cobra"
	"os"
)

// reconcileCmd represents the reconcile command
var reconcileCmd = &cobra.Command{
	Use:   "run",
	Short: "This will attempt to run a single reconcile for a given configuration.",
	Long: `This is a great command for running a one-off reconcile. A reconcile is the
same procedure that the controller will call in a loop.
So this command is useful for quick and dirty changes
or better yet as an easy way for devs to run their code.

The paradigm behind this subcommand is that it uses the
same internal libraries that the control loop will use.
So getting things working well here, is a good step in
getting things working well in the control loop.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := RunReconile(reco)
		if err != nil {
			logger.Critical(err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	},
}

type ReconcileOptions struct {
	RootOptions
	Filename string
	Name     string
}

var reco = &ReconcileOptions{}

func init() {
	RootCmd.AddCommand(reconcileCmd)
	reconcileCmd.Flags().StringVarP(&reco.Filename, "filename", "f", ".", "The filename to use to look for a Terraform configuration in. Accepts files or directories.")
	reconcileCmd.Flags().StringVarP(&reco.Name, "name", "n", namer.RandomName(), "The unique identifier name of the configuration. If no name is provided, terraformctl will generate one for you.")

}

func RunReconile(options *ReconcileOptions) error {

	// Init a terraform configuration
	tf, err := parser.NewTerraformConfigurationFromPath(options.Name, options.Filename)
	if err != nil {
		return fmt.Errorf("Unable to read terraform configuration with error: %v", err)
	}

	// Create and run against the Terraform SDK
	runner := controller.NewTerraformRunner(tf)
	err = runner.Reconcile()
	if err != nil {
		return fmt.Errorf("Error while attempting to reconcile: %v", err)
	}

	return nil
}
