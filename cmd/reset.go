package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var resetCmd = &cobra.Command{
	Use:   "reset [username]",
	Short: "Reset a database users password",
	Long:  `Reset a database users password`,
	Run: func(cmd *cobra.Command, args []string) {
		username, conn := givenUserModification(cmd, args, true)

		password, err := conn.ResetPassword(username)
		if err != nil {
			cmd.Println(fmt.Errorf("could not reset password for user: %v", err))
			os.Exit(1)
		}

		user, err := conn.GetUser(username)
		if err != nil {
			cmd.Println(fmt.Errorf("could not get user after password reset: %v", err))
			os.Exit(1)
		}

		output := getOutputType(cmd)
		if output == OutputTypeJson {
			outputResetJson(cmd, password)
		} else if output == OutputTypeTable {
			cmd.Println(fmt.Sprintf("Successfully reset user password for %s in %s", username, conn.Config.Host))
			cmd.Println("\ncredentials")
			cmd.Println("Valid until:", *user.ValidUntil)
			cmd.Println("Username:", username)
			cmd.Println("Password:", password)
			cmd.Println("Hostname:", conn.Config.Host)
			cmd.Println("Port:", conn.Config.Port)
			cmd.Println("Database:", conn.Config.Database)
		}
	},
}

func outputResetJson(cmd *cobra.Command, password string) {
	b, err := json.MarshalIndent(struct {
		Password string `json:"password"`
	}{
		Password: password,
	}, "", strings.Repeat(" ", 4))
	if err != nil {
		panic(err)
	}

	cmd.Println(string(b))
}

func init() {
	rootCmd.AddCommand(resetCmd)
	addRequiredHostFlag(resetCmd)
}
