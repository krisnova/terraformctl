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
	"os"

	"github.com/kris-nova/kubicorn/cutil/logger"
	"github.com/spf13/cobra"
	"strconv"
)

var cfgFile string

const (
	// Port defaults to 4392 because that is the elevation of Mt. Rainier (in meters).
	// which is the next mountain I will be climbing after I write this code.
	terraformctlDefaultPort = 4392
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "terraformctl",
	Short: "Run terraform as cloud native infrastructure as software",
	Long: fmt.Sprintf(`%s

A long time ago, in a galaxy far far away infrastructure was managed in many other ways.
This tool is an example of how we might want to start managing infrastructure in a cloud native way...
That is..

Infrastructure as software!`, banner()),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type RootOptions struct {
	Hostname string
	Port     int
}

var ro = &RootOptions{}

func init() {
	RootCmd.PersistentFlags().IntVarP(&logger.Level, "verbose", "v", 4, "The logger level (0 to 4) to use with 4 being the most verbose.")
	RootCmd.PersistentFlags().StringVarP(&ro.Hostname, "hostname", "H", strEnvDef("TERRAFORMCTL_HOSTNAME", "localhost"), "The hostname to use to connect to a listening terraformctl gRPC server. Will respect the $TERRAFORMCTL_HOSTNAME environmental variable more than anything else.")
	RootCmd.PersistentFlags().IntVarP(&ro.Port, "port", "p", intEnvDef("TERRAFORMCTL_PORT", terraformctlDefaultPort), "The port to use to connect to a listening terraformctl gRPC server. Will respect the $TERRAFORMCTL_HOSTNAME environmental variable more than anything else.")
}

// strEnvDef allows a user to override the default value in a cobra flag
func strEnvDef(env string, def string) string {
	val := os.Getenv(env)
	if val == "" {
		return def
	}
	return val
}

// strEnvDef allows a user to override the default value in a cobra flag
func intEnvDef(env string, def int) int {
	val := os.Getenv(env)
	if val == "" {
		return def
	}
	ival, err := strconv.Atoi(val)
	if err != nil {
		return def
	}
	return ival
}
