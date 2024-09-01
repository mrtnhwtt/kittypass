package cli

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/briandowns/spinner"
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
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
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
			s := spinner.New(spinner.CharSets[26], 150*time.Millisecond)
			s.Color("green")
			s.Prefix = "Creating Vault"
			s.Start()

			err := vault.CreateVault()
			if err != nil {
				s.FinalMSG = red("Vault creation failed.\n")
				s.Stop()
				return err
			}
			s.FinalMSG = green("✓ Vault created successfully.\n")
			s.Stop()
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
	login.Vault = &vault
	cmd := &cobra.Command{
		Use:     "login",
		Aliases: []string{"pass", "password"},
		Short:   "Create a new Login",
		Long:    "Create a new Login storing a username and password pair. If no password are provided, generates a new password for the login.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if login.Generator.Length < 5 || login.Generator.Length > 64 {
				return fmt.Errorf("invalid password length %d, please set a password length between 5 and 64", login.Generator.Length)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			err := login.Vault.Get()
			if err != nil {
				return err
			}

			s := spinner.New(spinner.CharSets[26], 150*time.Millisecond)
			s.Color("green")
			s.Prefix = "Checking Master Password"

			login.Vault.Masterpass = strings.TrimSpace(prompt.PasswordPrompt("Input master password:"))
			if login.Vault.Masterpass == "" {
				return errors.New("invalid empty master password")
			}

			s.Start()

			if err := login.Vault.MasterpassMatch(); err != nil {
				s.FinalMSG = red("Master Password check failed.\n")
				s.Stop()
				return err
			}
			err = login.Vault.RecreateDerivationKey()
			if err != nil {
				s.FinalMSG = red("Opening Vault failed.\n")
				s.Stop()
				return err
			}
			s.FinalMSG = green("✓ Successfully opened Vault.\n")
			s.Stop()

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
			s.Prefix = "Adding login to Vault"
			s.Start()
			err = login.Add()
			if err != nil {
				s.FinalMSG = red("Adding login to Vault failed.\n")
				s.Stop()
				return err
			}
			s.FinalMSG = fmt.Sprintf("%s %s %s %s.\n",green("✓ Successfully added login"), blue(login.Name), green("to Vault"), blue(login.Vault.Name))
			s.Stop()
			return nil
		},
	}
	cmd.Flags().StringVarP(&login.Name, "name", "n", "", "Name of the login")
	cmd.Flags().StringVarP(&login.Username, "username", "u", "", "Username or email for the login")
	cmd.Flags().StringVarP(&login.Vault.Name, "vault-name", "v", "", "Name of the Vault")
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
