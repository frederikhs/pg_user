package cmd

import (
	"github.com/hiperdk/pg_user/database"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
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

		listRolesForConnection(conn)
	},
}

func listRolesForConnection(conn *database.DBConn) {
	roles, err := conn.GetAllRoles()
	if err != nil {
		panic(err)
	}

	table := tablewriter.NewWriter(os.Stdout)
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
