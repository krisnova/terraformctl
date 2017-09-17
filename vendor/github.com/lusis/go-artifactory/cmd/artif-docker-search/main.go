package main

import (
	"fmt"
	"os"
	"strings"

	artifactory "github.com/lusis/go-artifactory/artifactory.v51"
	"github.com/olekukonko/tablewriter"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	criteria   = kingpin.Arg("criteria", "what to search for").Required().String()
	showLabels = kingpin.Flag("labels", "show labels").Bool()
)

func main() {
	kingpin.Parse()
	client := artifactory.NewClientFromEnv()
	data, err := client.DockerSearch(*criteria)
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
		"DESCRIPTION",
		"LAST MODIFIED",
		"MODIFIED BY",
	}
	if *showLabels {
		theaders = append(theaders, "LABELS")
	}
	table.SetHeader(theaders)

	for _, d := range data {
		var description string
		if d.Properties["docker.label.description"] != nil {
			description = d.Properties["docker.label.description"][0]
		}
		tdata := []string{
			fmt.Sprintf("%s:%s\n", d.Properties["docker.repoName"][0], d.Properties["docker.manifest"][0]),
			description,
			d.LastModified,
			d.ModifiedBy,
		}
		if *showLabels {
			var allLabels string
			var desc []string
			for p := range d.Properties {
				var labels []string
				if strings.HasPrefix(p, "docker.label") {
					labels = append(labels, p)
				}
				for _, label := range labels {
					desc = append(desc, fmt.Sprintf("%s = %s", label, d.Properties[label][0]))
				}
				allLabels = strings.Join(desc, "\n")
			}
			tdata = append(tdata, allLabels)
		}
		table.Append(tdata)
	}
	table.Render()
	os.Exit(0)
}
