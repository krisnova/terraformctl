package main

import (
	"fmt"
	"os"

	artifactory "github.com/lusis/go-artifactory/artifactory.v51"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	username = kingpin.Arg("username", "Username to delete").Required().String()
)

func main() {
	kingpin.Parse()
	client := artifactory.NewClientFromEnv()
	err := client.DeleteUser(*username)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	} else {
		fmt.Printf("User %s deleted\n", *username)
		os.Exit(0)
	}
}
