// Package cmd contains the CLI commands for speakeasy-helpers.
package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "speakeasy-helpers",
	Short: "A collection of helper utilities for working with Speakeasy generated API clients",
	Long: `speakeasy-helpers is a CLI tool that provides a set of utilities
to help you work more effectively with Speakeasy generated API clients.

This tool includes various subcommands for common tasks.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
}
