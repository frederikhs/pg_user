package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var extendCmd = &cobra.Command{
	Use:   "extend [username]",
	Short: "Extend a database user's valid",
	Long:  `Extend a database user's valid until`,
	Run: func(cmd *cobra.Command, args []string) {
		username, conn := givenUserModification(cmd, args, true)

		validUntil, err := conn.ExtendUser(username, DefaultValidUntil)
		if err != nil {
			cmd.Println(fmt.Errorf("could not extend user: %v", err))
			os.Exit(1)
		}

		cmd.Println(fmt.Sprintf("Successfully extended user %s in %s", username, conn.Config.Host))
		cmd.Println("\ncredentials")
		cmd.Println("Valid until:", validUntil.Format("2006-01-02"))
		cmd.Println("Username:", username)
		cmd.Println("Hostname:", conn.Config.Host)
		cmd.Println("Port:", conn.Config.Port)
		cmd.Println("Database:", conn.Config.Database)
	},
}

func init() {
	rootCmd.AddCommand(extendCmd)
	addRequiredHostFlag(extendCmd)
}
