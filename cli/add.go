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
	cmd := &cobra.Command{
		Use:     "add",
		Aliases: []string{"new"},
		Short:   "Add a new vault or login.",
		Long:    "Add a new vault or login. When creating a new vault, choose a master password to be used to add new logins",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Help()
			return nil
		},
	}

	cmd.AddCommand(
		NewAddVaultCmd(),
		NewAddLoginCmd(),
	)
	return cmd
}

func NewAddVaultCmd() *cobra.Command {
	vault := kittypass.NewVault()
	cmd := &cobra.Command{
		Use:     "vault",
		Aliases: []string{"folder"},
		Short:   "Create a new Vault.",
		Long:    "Create a new Vault to store login infornmation. Requires a master password.",
		RunE: func(cmd *cobra.Command, args []string) error {
			vault.Masterpass = strings.TrimSpace(prompt.PasswordPrompt("Input master password:"))
			if vault.Masterpass == "" {
				return errors.New("invalid empty master password")
			}
			confirm := strings.TrimSpace(prompt.PasswordPrompt("Confirm master password:"))
			if vault.Masterpass != confirm {
				return errors.New("master password does not match")
			}
			vault.HashMasterpass()

			err := vault.CreateVault()
			if err != nil {
				return err
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&vault.Name, "name", "n", "", "Name of the Vault")
	cmd.Flags().StringVarP(&vault.Description, "description", "d", "", "Description of the Vault")
	cmd.MarkFlagRequired("name")

	return cmd
}

func NewAddLoginCmd() *cobra.Command {
	login := kittypass.NewLogin()
	vault := kittypass.NewVault()
	cmd := &cobra.Command{
		Use:     "login",
		Aliases: []string{"password", "pass"},
		Short:   "Create a new Login",
		Long: "Create a new Login storing a username and password pair. If no password are provided, generates a new password for the login.",
		RunE: func(cmd *cobra.Command, args []string) error {
			vault.Masterpass = strings.TrimSpace(prompt.PasswordPrompt("Input master password:"))
			if vault.Masterpass == "" {
				return errors.New("invalid empty master password")
			}
			err := vault.Get()
			if err != nil {
				return err
			}
			if match := vault.MasterpassMatch(); match != nil {
				return fmt.Errorf("failed master password check: %s", match.Error())
			}
			
			if login.ProvidePassword {
				login.Password = strings.TrimSpace(prompt.PasswordPrompt("Input password:"))
				if login.Password == "" {
					return errors.New("invalide empty password")
				}
				confirm := strings.TrimSpace(prompt.PasswordPrompt("Confirm password:"))
				if login.Password != confirm {
					return errors.New("password does not match")
				}
			} else {
				login.Password = login.Generator.GeneratePassword()
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&login.Name, "name", "n", "", "Name of the login")
	cmd.Flags().StringVarP(&login.Username, "username", "u", "", "Username or email for the login")
	cmd.Flags().StringVarP(&vault.Name, "vault-name", "v", "", "Name of the Vault")
	cmd.Flags().BoolVarP(&login.ProvidePassword, "password", "p", false, "Use to set the password instead of generating a new password")
	cmd.Flags().IntVarP(&login.Generator.Length, "lenght", "l", 16, "Length of the generated password")
	cmd.Flags().BoolVarP(&login.Generator.SpecialChar, "special-char", "s", false, "Use special character in the password")
	cmd.Flags().BoolVarP(&login.Generator.Numeral, "numeral", "N", false, "Add number in the password")
	cmd.Flags().BoolVarP(&login.Generator.Uppercase, "uppercase", "U", false, "Use uppercase and lowercase characters")
	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("username")
	cmd.MarkFlagRequired("vault-name")
	return cmd
}
