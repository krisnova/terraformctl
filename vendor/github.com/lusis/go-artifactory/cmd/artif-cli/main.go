package main

import (
	"fmt"

	artifactory "github.com/lusis/go-artifactory/artifactory.v51"
)

func main() {
	client := artifactory.NewClientFromEnv()

	fmt.Printf("%#v\n", client)
}
