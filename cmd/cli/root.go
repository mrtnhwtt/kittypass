package cli

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kittypass",
		Short: "A CLI password manager",
		Long:  "A CLI password manager to securely stock your password. Interact with your password in a declarative or interactive mod.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if reterr, _ := cmd.Flags().GetBool("error"); reterr {
				return errors.New("an error")
			}
			fmt.Println("kittypass")
			return nil
		},
	}
	cmd.Flags().Bool("error", false, "Return error in RunE")
	return cmd
}
