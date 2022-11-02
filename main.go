package main

import (
	"fmt"
	"github.com/hiperdk/pg_user/cmd"
)

const notSet string = "not set"

// these information will be collected when built, by `-ldflags "-X main.appVersion=v1.0.3"`
var (
	appVersion = notSet
	buildTime  = notSet
	gitCommit  = notSet
	gitRef     = notSet
)

func buildVersion() string {
	return fmt.Sprintf("%s, %s, %s, %s", appVersion, buildTime, gitCommit, gitRef)
}

func main() {
	cmd.Execute(buildVersion())
}
