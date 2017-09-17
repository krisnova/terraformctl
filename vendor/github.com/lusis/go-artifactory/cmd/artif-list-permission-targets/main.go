package main

import (
	"fmt"
	"os"

	artifactory "github.com/lusis/go-artifactory/artifactory.v51"
	"github.com/olekukonko/tablewriter"
)

func main() {
	client := artifactory.NewClientFromEnv()
	data, err := client.GetPermissionTargets()
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	} else {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Uri"})
		table.SetAutoWrapText(false)
		for _, u := range data {
			table.Append([]string{u.Name, u.URI})
		}
		table.Render()
		os.Exit(0)
	}
}
