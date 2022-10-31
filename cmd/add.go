package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var addCmd = &cobra.Command{
	Use:   "add [username]",
	Short: "Add a database user",
	Long:  `Add a database user to a specific database`,
	Run: func(cmd *cobra.Command, args []string) {
		username, conn := givenUserModification(cmd, args, false)

		roles, err := cmd.Flags().GetStringSlice("roles")
		if err != nil {
			cmd.Println("roles flag nil")
			os.Exit(1)
		}

		password, validUntil, err := conn.CreateUser(username, DefaultValidUntil, roles)
		if err != nil {
			cmd.Println(fmt.Errorf("could not add user: %v", err))
			os.Exit(1)
		}

		cmd.Println(fmt.Sprintf("Successfully added user %s to %s", username, conn.Config.Host))
		cmd.Println("\ncredentials")
		cmd.Println("Valid until:", validUntil.Format("2006-01-02"))
		cmd.Println("Roles:", roles)
		cmd.Println("Username:", username)
		cmd.Println("Password:", password)
		cmd.Println("Hostname:", conn.Config.Host)
		cmd.Println("Port:", conn.Config.Port)
		cmd.Println("Database:", conn.Config.Database)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addRequiredHostFlag(addCmd)
	addCmd.Flags().StringSliceP("roles", "r", []string{}, "--roles=roleA,roleB (optional)")
}
