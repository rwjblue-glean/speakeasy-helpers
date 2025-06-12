package main

import (
	"os"

	"github.com/rwjblue-glean/speakeasy-helpers/cmd"
)

// Version information - these will be set by GoReleaser at build time
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmd.SetVersionInfo(version, commit, date)
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
