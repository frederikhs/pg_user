package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var alterCmd = &cobra.Command{
	Use:   "alter [username]",
	Short: "Alter a database users roles",
	Long:  "Alter a database users roles in a specific database",
	Run: func(cmd *cobra.Command, args []string) {
		username, conn := givenUserModification(cmd, 1, args, true)

		add_roles, err := cmd.Flags().GetStringSlice("add")
		if err != nil {
			cmd.Println(fmt.Errorf("could not alter user: %v", err))
			os.Exit(1)
		}

		remove_roles, err := cmd.Flags().GetStringSlice("remove")
		if err != nil {
			cmd.Println(fmt.Errorf("could not alter user: %v", err))
			os.Exit(1)
		}

		if len(add_roles) == 0 && len(remove_roles) == 0 {
			cmd.Println("at least one role must be added or removed")
			os.Exit(1)
		}

		all_roles := append(add_roles, remove_roles...)

		for _, role := range all_roles {
			exists, err := conn.RoleExist(role)

			if err != nil {
				cmd.Println(fmt.Errorf("could not alter user: %v", err))
				os.Exit(1)
			}

			if !exists {
				cmd.Println(fmt.Sprintf("role %s does not exist", role))
				os.Exit(1)
			}
		}

		tx := conn.BeginTransaction()

		err = conn.AddRole(tx, username, add_roles)
		if err != nil {
			cmd.Println(fmt.Errorf("could not add role: %v", err))
			os.Exit(1)
		}
		err = conn.RemoveRole(tx, username, remove_roles)
		if err != nil {
			cmd.Println(fmt.Errorf("could not remove role: %v", err))
			os.Exit(1)
		}

		err = tx.Commit()
		if err != nil {
			cmd.Println(fmt.Errorf("could not commit: %v", err))
			os.Exit(1)
		}

		output := getOutputType(cmd)
		if output == OutputTypeJson {
			outputAlterJson(cmd, add_roles, remove_roles)
		} else if output == OutputTypeTable {
			outputAlterTable(cmd, add_roles, remove_roles, username, conn.Config.Host)
		}
	},
}

func outputAlterTable(cmd *cobra.Command, add_roles, remove_roles []string, username, host string) {
	cmd.Println(fmt.Sprintf("successfully altered user %s in %s", username, host))
	cmd.Println("Added roles:", add_roles)
	cmd.Println("Removed roles:", remove_roles)
}

func outputAlterJson(cmd *cobra.Command, add_roles, remove_roles []string) {
	b, err := json.MarshalIndent(struct {
		Added   []string `json:"added"`
		Removed []string `json:"removed"`
	}{
		Added:   add_roles,
		Removed: remove_roles,
	}, "", strings.Repeat(" ", 4))
	if err != nil {
		panic(err)
	}

	cmd.Println(string(b))
}

func init() {
	rootCmd.AddCommand(alterCmd)
	addRequiredHostFlag(alterCmd)

	alterCmd.Flags().StringSlice("add", []string{}, "--add=roleA,roleB (optional)")
	alterCmd.Flags().StringSlice("remove", []string{}, "--remove=roleA,roleB (optional)")
}
