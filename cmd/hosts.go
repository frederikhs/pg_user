package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/frederikhs/pg_user/database"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/tg/pgpass"
	"os"
	"strings"
)

var hostsCmd = &cobra.Command{
	Use:   "hosts",
	Short: "List hosts from .pgpass",
	Long:  `List all hosts from your .pgpass file`,
	Run: func(cmd *cobra.Command, args []string) {
		entries, err := database.ListHosts()
		if err != nil {
			cmd.Println(fmt.Errorf("could not read .pgpass file: %v", err))
			os.Exit(1)
		}

		output := getOutputType(cmd)
		if output == OutputTypeJson {
			outputHostsJson(entries, cmd)
		} else if output == OutputTypeTable {
			outputHostsTable(entries, cmd)
		}
	},
}

func outputHostsTable(es []pgpass.Entry, cmd *cobra.Command) {
	out := cmd.OutOrStdout()

	table := tablewriter.NewWriter(out)
	table.SetHeader([]string{"host", "port", "database", "username", "password"})

	for _, e := range es {
		table.Append([]string{
			e.Hostname, e.Port, e.Database, e.Username, "********",
		})
	}

	table.Render()
}

type Entry struct {
	Hostname string `json:"hostname"`
	Port     string `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func outputHostsJson(es []pgpass.Entry, cmd *cobra.Command) {
	var entries []Entry
	for _, e := range es {
		entries = append(entries, Entry{
			Hostname: e.Hostname,
			Port:     e.Port,
			Database: e.Database,
			Username: e.Username,
			Password: "********",
		})
	}

	b, err := json.MarshalIndent(entries, "", strings.Repeat(" ", 4))
	if err != nil {
		panic(err)
	}

	cmd.Println(string(b))
}

func init() {
	rootCmd.AddCommand(hostsCmd)
}
