package cmd

import (
	"fmt"
	"github.com/hiperdk/pg_user/database"
	"github.com/spf13/cobra"
	"os"
	"regexp"
	"time"
)

var DefaultValidUntil = time.Hour * 24 * 180

var rootCmd = &cobra.Command{
	Use:   "pg_user",
	Short: "Manage database users",
	Long:  `Cli tool for managing database users across databases`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func addRequiredHostFlag(cmd *cobra.Command) {
	cmd.Flags().String("host", "", "--host database.example.com")
	cmd.MarkFlagRequired("host")
}

func givenUserModification(cmd *cobra.Command, args []string, userMustExist bool) (string, *database.DBConn) {
	if len(args) != 1 {
		cmd.Println(cmd.UsageString())
		os.Exit(1)
	}

	hostnameFlag := cmd.Flag("host")
	if hostnameFlag == nil {
		cmd.Println("hostname flag nil")
		os.Exit(1)
	}

	host := hostnameFlag.Value.String()

	username := args[0]

	conn := database.GetDatabaseEntryConnection(host)
	if conn == nil {
		cmd.Println("no hosts found by that name")
		os.Exit(1)
	}

	alphanumeric := regexp.MustCompile(`^[a-zA-Z0-9@._-]*$`).MatchString(username)
	if !alphanumeric {
		cmd.Println("username must be in character set of a-zA-Z0-9@._-")
		os.Exit(1)
	}

	exists, err := conn.UserExist(username)
	if err != nil {
		cmd.Println(fmt.Errorf("could not check if user exists: %v", err))
		os.Exit(1)
	}

	if userMustExist && !exists {
		cmd.Println("username does not exist")
		os.Exit(1)
	}

	if !userMustExist && exists {
		cmd.Println("username already exists")
		os.Exit(1)
	}

	return username, conn

}
