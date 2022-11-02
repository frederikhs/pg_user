package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
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

		output := getOutputType(cmd)
		if output == OutputTypeJson {
			outputDeleteJson(cmd)
		} else if output == OutputTypeTable {
			cmd.Println(fmt.Sprintf("successfully deleted user %s from %s", username, conn.Config.Host))
		}
	},
}

func outputDeleteJson(cmd *cobra.Command) {
	b, err := json.MarshalIndent(struct {
		Ok bool `json:"ok"`
	}{
		Ok: true,
	}, "", strings.Repeat(" ", 4))
	if err != nil {
		panic(err)
	}

	cmd.Println(string(b))
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	addRequiredHostFlag(deleteCmd)
}
