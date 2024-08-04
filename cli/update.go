package cli

import (
	"errors"
	"strings"

	"github.com/mrtnhwtt/kittypass/internal/kittypass"
	"github.com/mrtnhwtt/kittypass/internal/prompt"
	"github.com/spf13/cobra"
)

func NewUpdateCmd() *cobra.Command {
	kp := kittypass.NewLogin()
	var targetName string
	cmd := &cobra.Command{
		Use:     "update",
		Aliases: []string{"put", "modify", "mod"},
		Short:   "update a login",
		Long:    "update a login password or username from your secret storage",
		RunE: func(cmd *cobra.Command, args []string) error {
			setPassword, err := cmd.Flags().GetBool("password")
			if err != nil {
				return err
			}
			if setPassword {
				master := strings.TrimSpace(prompt.PasswordPrompt("Input master password:"))
				if master == "" {
					return errors.New("invalid empty master password")
				}
				kp.Vault.UseMasterPassword()

				kp.Password = strings.TrimSpace(prompt.PasswordPrompt("Input new password:"))
				if kp.Password == "" {
					return errors.New("invalid empty password")
				}
			}
			kp.Update(targetName)
			return nil
		},
	}
	cmd.Flags().StringVarP(&targetName, "target", "t", "", "Name of the login to update")
	cmd.Flags().StringVarP(&kp.Name, "new-name", "n", "", "New name for the updated login")
	cmd.Flags().StringVarP(&kp.Username, "new-username", "u", "", "New Username or email for the login")
	cmd.Flags().BoolP("password", "p", false, "prompt to set a new password for the login")
	cmd.MarkFlagRequired("target")
	cmd.MarkFlagsOneRequired("password", "new-name", "new-username")
	return cmd
}
