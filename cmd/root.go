// Package cmd contains the CLI commands for speakeasy-helpers.
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version information - will be set by SetVersionInfo
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "speakeasy-helpers",
	Short: "A collection of helper utilities for working with Speakeasy generated API clients",
	Long: `speakeasy-helpers is a CLI tool that provides a set of utilities
to help you work more effectively with Speakeasy generated API clients.

This tool includes various subcommands for common tasks.`,
	Version: version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

// SetVersionInfo sets the version information for the CLI
func SetVersionInfo(v, c, d string) {
	version = v
	commit = c
	date = d
	rootCmd.Version = version
	rootCmd.SetVersionTemplate(fmt.Sprintf("speakeasy-helpers %s (commit: %s, built: %s)\n", version, commit, date))
}

func init() {
	rootCmd.SetVersionTemplate(fmt.Sprintf("speakeasy-helpers %s (commit: %s, built: %s)\n", version, commit, date))
}
