package cli

import (
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "kittypass",
		Short: "A CLI password manager",
		Long:  "A CLI password manager to securely stock your password. Interact with your password in a declarative or interactive mod.",
	}

	cmd.AddCommand(
		NewAddCmd(),
		NewGetCmd(),
		NewListCmd(),
		NewDeleteCmd(),
	)

	return cmd
}
