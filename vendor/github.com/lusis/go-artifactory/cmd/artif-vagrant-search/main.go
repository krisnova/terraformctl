package main

import (
	"fmt"
	"os"

	artifactory "github.com/lusis/go-artifactory/artifactory.v51"
	"github.com/olekukonko/tablewriter"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	criteria = kingpin.Arg("criteria", "what to search for").Required().String()
)

func main() {
	kingpin.Parse()
	client := artifactory.NewClientFromEnv()
	data, err := client.VagrantSearch(*criteria)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)
	table.SetBorder(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	theaders := []string{
		"NAME",
		"VERSION",
		"PROVIDER",
		"MODIFIED",
		"MODIFIED BY",
	}
	table.SetHeader(theaders)

	for _, d := range data {
		props := make(map[string]string)
		for _, prop := range d.Properties {
			props[prop.Key] = prop.Value
		}
		tdata := []string{
			fmt.Sprintf("%s/%s\n", d.Repo, props["box_name"]),
			props["box_version"],
			props["box_provider"],
			d.Modified,
			d.ModifiedBy,
		}
		table.Append(tdata)
	}
	table.Render()
	os.Exit(0)
}
