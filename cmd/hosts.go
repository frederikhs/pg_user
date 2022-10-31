package cmd

import (
	"fmt"
	"github.com/hiperdk/pg_user/database"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
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

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"host", "port", "database", "username"})

		for _, e := range entries {
			table.Append([]string{
				e.Hostname, e.Port, e.Database, e.Username,
			})
		}

		table.Render()
	},
}

func init() {
	rootCmd.AddCommand(hostsCmd)
}
