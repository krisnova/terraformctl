package main

import (
	"fmt"
	"os"

	artifactory "github.com/lusis/go-artifactory/artifactory.v51"
)

func main() {
	client := artifactory.NewClientFromEnv()
	p, err := client.CreateUserAPIKey()
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	} else {
		fmt.Printf("%s\n", p)
		os.Exit(0)
	}
}
