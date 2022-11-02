package cmd

import (
	"encoding/json"
	"github.com/hiperdk/pg_user/database"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var rolesCmd = &cobra.Command{
	Use:   "roles [host]",
	Short: "List database roles",
	Long:  `List database roles for a specific database`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.Println(cmd.UsageString())
			os.Exit(1)
		}

		host := args[0]

		conn := database.GetDatabaseEntryConnection(host)
		if conn == nil {
			cmd.Println("no hosts found by that name")
			os.Exit(1)
		}

		roles, err := conn.GetAllRoles()
		if err != nil {
			panic(err)
		}

		output := getOutputType(cmd)
		if output == OutputTypeJson {
			outputRolesJson(roles, cmd)
		} else if output == OutputTypeTable {
			outputRolesTable(roles, cmd)
		}
	},
}

func outputRolesJson(roles []string, cmd *cobra.Command) {
	b, err := json.MarshalIndent(roles, "", strings.Repeat(" ", 4))
	if err != nil {
		panic(err)
	}

	cmd.Println(string(b))
}

func outputRolesTable(roles []string, cmd *cobra.Command) {
	out := cmd.OutOrStdout()

	table := tablewriter.NewWriter(out)
	table.SetHeader([]string{"rolename"})

	for _, r := range roles {
		table.Append([]string{
			r,
		})
	}

	table.Render()
}

func init() {
	rootCmd.AddCommand(rolesCmd)
}
