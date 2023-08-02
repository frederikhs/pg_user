package cmd

import (
	"fmt"
	"github.com/frederikhs/pg_user/database"
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

func Execute(buildVersion string) {
	rootCmd.Version = buildVersion

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("output", "o", "", "--output json|table (optional default table)")
	rootCmd.PersistentFlags().Bool("no-ssl", false, "do not use ssl when connecting")
}

func addRequiredHostFlag(cmd *cobra.Command) {
	cmd.Flags().String("host", "", "--host database.example.com")
	err := cmd.MarkFlagRequired("host")
	if err != nil {
		panic(err)
	}
}

func givenUserModification(cmd *cobra.Command, argsLen int, args []string, userMustExist bool) (string, *database.DBConn) {
	if len(args) != argsLen {
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

	withSSL := getUseSSL(cmd)

	conn := database.GetDatabaseEntryConnection(host, withSSL)
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

type OutputType string

const (
	OutputTypeTable = "table"
	OutputTypeJson  = "json"
)

func getOutputType(cmd *cobra.Command) OutputType {
	output, err := cmd.Flags().GetString("output")
	if err != nil {
		cmd.Println(err)
		os.Exit(1)
	}

	if output == "json" {
		return OutputTypeJson
	} else if output == "table" || output == "" {
		return OutputTypeTable
	} else {
		cmd.Println(cmd.UsageString())
		os.Exit(1)
		// for the compiler to be happy
		return OutputTypeTable
	}
}

func getUseSSL(cmd *cobra.Command) bool {
	ssl, err := cmd.Flags().GetBool("no-ssl")
	if err != nil {
		cmd.Println(err)
		os.Exit(1)
	}

	return !ssl
}
