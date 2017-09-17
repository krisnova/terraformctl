package main

import (
	"fmt"
	"os"

	artifactory "github.com/lusis/go-artifactory/artifactory.v51"
	"github.com/olekukonko/tablewriter"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	repo = kingpin.Arg("repo", "repo to list files").Required().String()
)

func main() {
	kingpin.Parse()
	client := artifactory.NewClientFromEnv()
	u, err := client.ListFiles(*repo)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	} else {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"URI", "Size", "SHA-1"})
		for _, v := range u.Files {
			table.Append([]string{v.URI, fmt.Sprintf("%d", v.Size), v.SHA1})
		}

		table.Render()
		os.Exit(0)
	}
}
