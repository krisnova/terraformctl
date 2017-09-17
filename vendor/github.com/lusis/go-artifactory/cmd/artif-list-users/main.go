package main

import (
	"fmt"
	"os"

	artifactory "github.com/lusis/go-artifactory/artifactory.v51"
	"github.com/olekukonko/tablewriter"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	formatUsage = fmt.Sprintf("Format to show results [table, csv, list (usernames only - useful for piping)]")
	format      = kingpin.Flag("format", formatUsage).Short('F').Default("table").Enum("table", "list", "csv")
	sep         = kingpin.Flag("separator", "separator for csv output").Default(",").String()
)

func main() {
	kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version("1.0").Author("John E. Vincent")
	kingpin.CommandLine.Help = "List all users in Artifactory"
	kingpin.Parse()

	client := artifactory.NewClientFromEnv()
	data, err := client.GetUsers()
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	} else {
		if *format == "table" {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Name", "Uri"})
			table.SetAutoWrapText(false)
			for _, u := range data {
				table.Append([]string{u.Name, u.URI})
			}
			table.Render()
		} else if *format == "list" {
			for _, u := range data {
				fmt.Printf("%s\n", u.Name)
			}
		} else if *format == "csv" {
			for _, u := range data {
				fmt.Printf("%s%s%s\n", u.Name, *sep, u.URI)
			}
		}
		os.Exit(0)
	}
}
