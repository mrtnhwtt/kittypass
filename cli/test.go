package cli

import (
	"fmt"

	"github.com/mrtnhwtt/kittypass/internal/kittypass"
	"github.com/spf13/cobra"
)

func NewTestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "test",
		RunE: func(cmd *cobra.Command, args []string) error {
			generator := kittypass.PasswordGenerator{
				Length:      5,
				SpecialChar: true,
				Numeral:     true,
				Uppercase:   true,
			}
			fmt.Printf("password: %s\n", generator.GeneratePassword())

			return nil
		},
	}

	cmd.AddCommand(
		NewAddVaultCmd(),
		NewAddLoginCmd(),
	)
	return cmd
}
