package main

import (
	"fmt"
	"github.com/frederikhs/pg_user/cmd"
)

// these information will be collected when built, by goreleaser
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func buildVersion() string {
	return fmt.Sprintf("%s, %s, %s", version, commit, date)
}

func main() {
	cmd.Execute(buildVersion())
}
