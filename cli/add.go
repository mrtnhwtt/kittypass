package cli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mrtnhwtt/kittypass/internal/kittypass"
	"github.com/mrtnhwtt/kittypass/internal/prompt"
	"github.com/spf13/cobra"
)

func NewAddCmd() *cobra.Command {
	kp := kittypass.New()

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new password",
		Long:  "Add a new named password and email pair to your secret storage",
		RunE: func(cmd *cobra.Command, args []string) error {
			master := strings.TrimSpace(prompt.PasswordPrompt("Input master password:"))
			if master == "" {
				return errors.New("invalid empty master password")
			}
			kp.UseMasterPassword(master)

			kp.Password = strings.TrimSpace(prompt.PasswordPrompt("Input new password:"))
			if kp.Password == "" {
				return errors.New("invalid empty password")
			}
			fmt.Println(kp)
			err := kp.Add()
			if err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&kp.Name, "name", "n", "", "Secret's identifier")
	cmd.Flags().StringVarP(&kp.Username, "username", "u", "", "Username or email associated with the password")
	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("username")
	return cmd
}
