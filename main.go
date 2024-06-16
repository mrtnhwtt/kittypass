package main

import (
	"os"

	root "github.com/mrtnhwtt/kittypass/cmd/cli"
)

func main() {
	cmd := root.NewRootCmd()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
