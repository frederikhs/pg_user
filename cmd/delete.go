package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [username]",
	Short: "Delete a database user",
	Long:  `Delete a database user from a specific database`,
	Run: func(cmd *cobra.Command, args []string) {
		username, conn := givenUserModification(cmd, args, true)

		err := conn.DeleteUser(username)
		if err != nil {
			cmd.Println(fmt.Errorf("could not delete user: %v", err))
			os.Exit(1)
		}

		cmd.Println(fmt.Sprintf("successfully deleted user %s from %s", username, conn.Config.Host))
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	addRequiredHostFlag(deleteCmd)
}
