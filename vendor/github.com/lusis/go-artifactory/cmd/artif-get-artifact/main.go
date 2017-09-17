package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	artifactory "github.com/lusis/go-artifactory/artifactory.v51"
)

var (
	repo   = kingpin.Arg("repo", "repository key for download").Required().String()
	file   = kingpin.Arg("filename", "full path and file to download").Required().String()
	output = kingpin.Flag("output", "output file").String()
	silent = kingpin.Flag("silent", "supress output").Bool()
)

func main() {
	kingpin.Parse()
	client := artifactory.NewClientFromEnv()
	_, destination := filepath.Split(*file)
	if *output != "" {
		destination = *output
	} else {
		curdir, err := os.Getwd()
		if err != nil {
			fmt.Printf("Unable to get current directory: %s", err.Error())
			os.Exit(1)
		}
		destination = curdir + "/" + destination
	}

	if !*silent {
		fmt.Printf("Writing file to: %s", destination)
	}
	i, err := client.RetrieveArtifact(*repo, *file)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	} else {
		err := ioutil.WriteFile(destination, i, 0600)
		if err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(1)
		} else {
			os.Exit(0)
		}
	}
}
