package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	artifactory "github.com/lusis/go-artifactory/artifactory.v51"
	"github.com/olekukonko/tablewriter"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	formatUsage = fmt.Sprintf("Format to show results [table]")
	user        = kingpin.Arg("user", "User name to show").Required().String()
	format      = kingpin.Flag("format", formatUsage).Short('F').Default("table").Enum("table")
	attrs       = kingpin.Flag("attrs", "Columns to display. Cumulative last argument (i.e. -A name email)").Short('A').Strings()
)

func main() {
	kingpin.Parse()
	fmt.Printf("%s\n", *attrs)
	client := artifactory.NewClientFromEnv()
	u, err := client.GetUserDetails(*user)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	} else {
		if *format == "table" {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{
				"Name",
				"Email",
				"Password",
				"Admin?",
				"Updatable?",
				"Last Logged In",
				"Internal Password Disabled?",
				"Realm",
				"Groups",
			})
			table.SetAutoWrapText(false)
			table.Append([]string{
				u.Name,
				u.Email,
				"<hidden>",
				strconv.FormatBool(u.Admin),
				strconv.FormatBool(u.ProfileUpdatable),
				u.LastLoggedIn,
				strconv.FormatBool(u.InternalPasswordDisabled),
				u.Realm,
				strings.Join(u.Groups, "\n"),
			})
			table.Render()
		} else {
			fmt.Printf("Unknown format: %s", *format)
		}
		os.Exit(0)
	}
}
