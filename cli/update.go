package cli

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/mrtnhwtt/kittypass/internal/kittypass"
	"github.com/mrtnhwtt/kittypass/internal/prompt"
	"github.com/spf13/cobra"
)

func NewUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update",
		Aliases: []string{"put", "modify", "mod"},
		Short:   "update a login",
		Long:    "update a login password or username from your secret storage",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(
		NewUpdateLoginCmd(),
		NewUpdateVaultCmd(),
	)

	return cmd
}

func NewUpdateLoginCmd() *cobra.Command {
	login := kittypass.NewLogin()
	vault := kittypass.NewVault()
	login.Vault = &vault
	var targetName string

	cmd := &cobra.Command{
		Use:     "login",
		Aliases: []string{"put", "modify", "mod"},
		Short:   "update a login",
		Long:    "update a login password or username from your secret storage",
		RunE: func(cmd *cobra.Command, args []string) error {
			setPassword, err := cmd.Flags().GetBool("password")
			if err != nil {
				return err
			}
			generatePassword, err := cmd.Flags().GetBool("generate")
			if err != nil {
				return err
			}
			if setPassword || generatePassword {
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
					return fmt.Errorf("failed master password check: %s", match.Error())
				}
				err = login.Vault.RecreateDerivationKey()
				if err != nil {
					s.FinalMSG = red("Opening Vault failed.\n")
					s.Stop()
					return err
				}
				s.FinalMSG = green("Successfully opened Vault.\n")
				s.Stop()
				if setPassword {
					login.Password = strings.TrimSpace(prompt.PasswordPrompt("Input new password:"))
					if login.Password == "" {
						return errors.New("invalid empty password")
					}
					confirm := strings.TrimSpace(prompt.PasswordPrompt("Confirm new password:"))
					if login.Password != confirm {
						return errors.New("password does not match")
					}
				}
				if generatePassword {
					login.Password = login.Generator.GeneratePassword()
				}

			}

			if login.Vault.Name == "" {
				return fmt.Errorf("invalid vault name")
			}
			err = login.Vault.Get()
			if err != nil {
				return err
			}
			login.Update(targetName)
			return nil
		},
	}

	cmd.Flags().StringVarP(&targetName, "target", "t", "", "Name of the login to update")
	cmd.Flags().StringVarP(&login.Name, "new-name", "n", "", "New name for the updated login")
	cmd.Flags().StringVarP(&login.Username, "new-username", "u", "", "New Username or email for the login")
	cmd.Flags().StringVarP(&login.Vault.Name, "vault", "v", "", "vault of the target login")
	cmd.Flags().BoolP("password", "p", false, "prompt to set a user provided new password for the login")
	cmd.Flags().BoolP("generate", "g", false, "Generate a new password")
	cmd.Flags().IntVarP(&login.Generator.Length, "lenght", "l", 16, "Length of the generated password")
	cmd.Flags().BoolVarP(&login.Generator.SpecialChar, "special-char", "s", false, "Use special character in the password")
	cmd.Flags().BoolVarP(&login.Generator.Numeral, "numeral", "N", false, "Add number in the password")
	cmd.Flags().BoolVarP(&login.Generator.Uppercase, "uppercase", "U", false, "Use uppercase and lowercase characters")

	cmd.MarkFlagsMutuallyExclusive("password", "generate")
	cmd.MarkFlagRequired("target")
	cmd.MarkFlagRequired("vault")
	cmd.MarkFlagsOneRequired("password", "new-name", "new-username", "generate")

	return cmd
}

func NewUpdateVaultCmd() *cobra.Command {
	vault := kittypass.NewVault()
	var newName, newDescription, newPassword string
	var setNewPass bool
	cmd := &cobra.Command{
		Use:     "vault",
		Aliases: []string{"folder"},
		Short:   "update a vault",
		Long:    "update a vault name, description or master password",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := vault.Get()
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return errors.New("vault not found. Please check the vault name and try again")
				}
				return fmt.Errorf("an error occurred while retrieving the vault: %s", err)
			}
			s := spinner.New(spinner.CharSets[26], 150*time.Millisecond)
			s.Color("green")
			s.Prefix = "Checking Master Password"

			vault.Masterpass = strings.TrimSpace(prompt.PasswordPrompt("Input master password:"))
			if vault.Masterpass == "" {
				return errors.New("invalid empty master password")
			}

			s.Start()
			if match := vault.MasterpassMatch(); match != nil {
				s.FinalMSG = red("Master Password check failed.\n")
				s.Stop()
				return fmt.Errorf("failed master password check: %s", match.Error())
			}
			s.FinalMSG = green("Successfully opened Vault.\n")
			s.Stop()

			if setNewPass {
				newPassword = strings.TrimSpace(prompt.PasswordPrompt("Input new master password:"))
				if newPassword == "" {
					return errors.New("invalid empty master password")
				}
				confirm := strings.TrimSpace(prompt.PasswordPrompt("Confirm new master password:"))
				if newPassword != confirm {
					return errors.New("master password does not match")
				}
				err = vault.RecreateDerivationKey()
				if err != nil {
					return err
				}
			}
			affected, err := vault.Update(newPassword, newName, newDescription)
			fmt.Printf("%v", affected)
			return err
			// return nil
		},
	}
	cmd.Flags().StringVarP(&vault.Name, "target", "t", "", "Name of the vault to update")
	cmd.Flags().StringVarP(&newName, "new-name", "n", "", "New name for the vault")
	cmd.Flags().StringVarP(&newDescription, "new-description", "d", "", "New description for the vault")
	cmd.Flags().BoolVarP(&setNewPass, "new-password", "p", false, "prompt to set a new master password")
	cmd.MarkFlagRequired("target")
	cmd.MarkFlagsOneRequired("new-name", "new-description", "new-password")
	return cmd
}
