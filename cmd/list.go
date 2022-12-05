package cmd

import (
	"encoding/json"
	"github.com/frederikhs/pg_user/database"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var listCmd = &cobra.Command{
	Use:   "list [host]",
	Short: "List database users",
	Long:  `List database users for a specific database`,
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

		users, err := conn.GetAllUsers()
		if err != nil {
			cmd.Println(err)
			os.Exit(1)
		}

		output := getOutputType(cmd)
		if output == OutputTypeJson {
			outputListJson(users, cmd)
		} else if output == OutputTypeTable {
			outputListTable(users, cmd)
		}
	},
}

func outputListJson(users []database.User, cmd *cobra.Command) {
	b, err := json.MarshalIndent(users, "", strings.Repeat(" ", 4))
	if err != nil {
		panic(err)
	}

	cmd.Println(string(b))
}

func outputListTable(users []database.User, cmd *cobra.Command) {
	out := cmd.OutOrStdout()

	table := tablewriter.NewWriter(out)
	table.SetHeader([]string{"username", "valid until", "roles"})

	for _, u := range users {
		table.Append([]string{
			u.Username, *u.ValidUntil, strings.Join(u.Roles, ", "),
		})
	}

	table.Render()
}

func init() {
	rootCmd.AddCommand(listCmd)
}
