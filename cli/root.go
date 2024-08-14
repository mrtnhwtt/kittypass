package cli

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	
)

var (
	green = color.New(color.FgGreen).SprintFunc()
	red   = color.New(color.FgRed).SprintFunc()
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
		NewUpdateCmd(),
	)

	return cmd
}
