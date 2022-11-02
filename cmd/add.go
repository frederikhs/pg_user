package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"time"
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

		output := getOutputType(cmd)
		if output == OutputTypeJson {
			outputAddJson(cmd, validUntil, roles, username, password, conn.Config.Host, conn.Config.Port, conn.Config.Database)
		} else if output == OutputTypeTable {
			outputAddTable(cmd, validUntil, roles, username, password, conn.Config.Host, conn.Config.Port, conn.Config.Database)
		}
	},
}

func outputAddTable(cmd *cobra.Command, validUntil time.Time, roles []string, username, password, host, port, database string) {
	cmd.Println(fmt.Sprintf("Successfully added user %s to %s", username, host))
	cmd.Println("\ncredentials")
	cmd.Println("Valid until:", validUntil.Format("2006-01-02"))
	cmd.Println("Roles:", roles)
	cmd.Println("Username:", username)
	cmd.Println("Password:", password)
	cmd.Println("Hostname:", host)
	cmd.Println("Port:", port)
	cmd.Println("Database:", database)
}

func outputAddJson(cmd *cobra.Command, validUntil time.Time, roles []string, username, password, host, port, database string) {
	b, err := json.MarshalIndent(struct {
		ValidUntil string   `json:"valid_until"`
		Roles      []string `json:"roles"`
		Username   string   `json:"username"`
		Password   string   `json:"password"`
		Hostname   string   `json:"hostname"`
		Port       string   `json:"port"`
		Database   string   `json:"database"`
	}{
		ValidUntil: validUntil.Format("2006-01-02"),
		Roles:      roles,
		Username:   username,
		Password:   password,
		Hostname:   host,
		Port:       port,
		Database:   database,
	}, "", strings.Repeat(" ", 4))
	if err != nil {
		panic(err)
	}

	cmd.Println(string(b))
}

func init() {
	rootCmd.AddCommand(addCmd)
	addRequiredHostFlag(addCmd)
	addCmd.Flags().StringSliceP("roles", "r", []string{}, "--roles=roleA,roleB (optional)")
}
