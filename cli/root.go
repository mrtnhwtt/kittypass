package cli

import (
	"errors"

	"github.com/mrtnhwtt/kittypass/internal/kittypass"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	kp := kittypass.New()

	cmd := &cobra.Command{
		Use:   "kittypass",
		Short: "A CLI password manager",
		Long:  "A CLI password manager to securely stock your password. Interact with your password in a declarative or interactive mod.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if kp.Password == "" {
				return errors.New("no master password provided")
			}
			kp.Run(args[0])

			return nil
		},
	}
	cmd.Flags().StringVarP(&kp.Password, "password", "p", "", "Master password")
	return cmd
}
