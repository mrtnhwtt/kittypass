package cli

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/mrtnhwtt/kittypass/internal/kittypass"
	"github.com/mrtnhwtt/kittypass/internal/prompt"
	"github.com/mrtnhwtt/kittypass/internal/utils"
	"github.com/spf13/cobra"
)

func NewGetCmd() *cobra.Command {
	login := kittypass.NewLogin()
	vault := kittypass.NewVault()
	login.Vault = &vault

	cmd := &cobra.Command{
		Use:     "get",
		Aliases: []string{"fetch", "copy"},
		Short:   "get a login",
		Long:    "get a login from a vault, adds the password to the clipboard",
		RunE: func(cmd *cobra.Command, args []string) error {
			s := spinner.New(spinner.CharSets[26], 150*time.Millisecond)
			s.Color("green")
			s.Prefix = "Checking Master Password"

			login.Vault.Masterpass = strings.TrimSpace(prompt.PasswordPrompt("Input master password:"))
			if login.Vault.Masterpass == "" {
				return errors.New("invalid empty master password")
			}

			s.Start()
			err := login.Vault.Get()
			if err != nil {
				s.FinalMSG = red("Master Password check failed.\n")
				s.Stop()
				return err
			}
			if match := login.Vault.MasterpassMatch(); match != nil {
				s.FinalMSG = red("Master Password check failed.\n")
				s.Stop()
				return errors.New("incorrect password")
			}
			err = login.Vault.RecreateDerivationKey()
			if err != nil {
				s.FinalMSG = red("Opening Vault failed.\n")
				s.Stop()
				return err
			}
			s.FinalMSG = green("✓ Successfully opened Vault.\n")
			s.Stop()
			login, err := login.Get()
			if err != nil {
				return err
			}
			fmt.Printf("\n%s%s\n", blue("Login Name: "), login["name"])
			fmt.Printf("%s%s\n", blue("Usename: "), login["username"])
			err = utils.AddToClipboard(login["password"])
			if err != nil {
				fmt.Printf("%s%s\n", blue("Password: "), login["password"])
				fmt.Println(red("Failed to add password to the clipboard, printed password to the console."))
			} else {
				fmt.Println(green("password added to clipboard"))
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&login.Name, "name", "n", "", "login's name")
	cmd.Flags().StringVarP(&login.Vault.Name, "vault", "v", "", "vault's name")
	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("vault")
	return cmd
}
