package cmd

import (
	"github.com/hiperdk/pg_user/database"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
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

		listUsersForConnection(conn)
	},
}

func listUsersForConnection(conn *database.DBConn) {
	users, err := conn.GetAllUsers()
	if err != nil {
		panic(err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"username", "valid until"})

	for _, u := range users {
		table.Append([]string{
			u.Username, *u.ValidUntil,
		})
	}

	table.Render()
}

func init() {
	rootCmd.AddCommand(listCmd)
}
