package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var passwordCmd = &cobra.Command{
	Use:   "password [username] [password]",
	Short: "Set a database user's password",
	Long:  `Set a database user's password`,
	Run: func(cmd *cobra.Command, args []string) {
		username, conn := givenUserModification(cmd, 2, args, true)
		password := args[1]

		password, err := conn.SetPassword(username, password)
		if err != nil {
			cmd.Println(fmt.Errorf("could not set password for user: %v", err))
			os.Exit(1)
		}

		user, err := conn.GetUser(username)
		if err != nil {
			cmd.Println(fmt.Errorf("could not get user after password set: %v", err))
			os.Exit(1)
		}

		output := getOutputType(cmd)
		if output == OutputTypeJson {
			outputSetPasswordJson(cmd, password)
		} else if output == OutputTypeTable {
			cmd.Println(fmt.Sprintf("Successfully set user password for %s in %s", username, conn.Config.Host))
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

func outputSetPasswordJson(cmd *cobra.Command, password string) {
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
	rootCmd.AddCommand(passwordCmd)
	addRequiredHostFlag(passwordCmd)
}
