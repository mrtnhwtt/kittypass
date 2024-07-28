package cli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mrtnhwtt/kittypass/internal/kittypass"
	"github.com/mrtnhwtt/kittypass/internal/prompt"
	"github.com/spf13/cobra"
)

func NewGetCmd() *cobra.Command {
	kp := kittypass.New()

	cmd := &cobra.Command{
		Use:   "get",
		Short: "get a saved password",
		Long:  "get a saved named password and email pair from your secret storage",
		RunE: func(cmd *cobra.Command, args []string) error {
			master := strings.TrimSpace(prompt.PasswordPrompt("Input master password:"))
			if master == "" {
				return errors.New("invalid empty master password")
			}
			login, err := kp.Get(master)
			if err != nil {
				return err
			}
			fmt.Println(login)
			return nil
		},
	}
	cmd.Flags().StringVarP(&kp.Name, "name", "n", "", "Secret's identifier")
	cmd.MarkFlagRequired("name")
	return cmd
}
