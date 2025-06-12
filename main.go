package main

import (
	"os"

	"github.com/rwjblue-glean/speakeasy-helpers/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
