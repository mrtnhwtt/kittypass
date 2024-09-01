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

func NewDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete",
		Aliases: []string{"rm", "remove"},
		Short:   "delete a login or a vault",
		Long:    "delete a login or a vault",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	cmd.AddCommand(
		NewDeleteLoginCmd(),
		NewDeleteVaultCmd(),
	)
	return cmd
}

func NewDeleteLoginCmd() *cobra.Command {
	login := kittypass.NewLogin()
	vault := kittypass.NewVault()
	login.Vault = &vault

	cmd := &cobra.Command{
		Use:     "login",
		Aliases: []string{"pass", "password"},
		Short:   "delete a saved login",
		Long:    "delete a saved logins from a vault using the login and vault name",
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
			if err := login.Vault.MasterpassMatch(); err != nil {
				s.FinalMSG = red("Master Password check failed.\n")
				s.Stop()
				return err
			}
			s.FinalMSG = green("✓ Successfully opened Vault.\n")
			s.Stop()
			err = login.Delete()
			if err != nil {
				fmt.Println(red("Failed to delete login"))
				return err
			}
			fmt.Println(green("✓ Successfully deleted login"))
			return nil
		},
	}
	cmd.Flags().StringVarP(&login.Name, "name", "n", "", "login's name")
	cmd.Flags().StringVarP(&login.Vault.Name, "vault", "v", "", "vault's name")
	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("vault")
	return cmd
}

func NewDeleteVaultCmd() *cobra.Command {
	vault := kittypass.NewVault()

	cmd := &cobra.Command{
		Use:     "vault",
		Aliases: []string{"folder"},
		Short:   "Delete a Vault.",
		Long:    "Delete a Vault and all associated logins. Requires the master password.",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := vault.Get()
			if err != nil {
				return err
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
				s.FinalMSG = red("Failed master password check.\n")
				s.Stop()
				return err
			}
			s.FinalMSG = green("✓ Successfully opened Vault for deletion.\n")
			s.Stop()
			deleted, err := vault.Delete()
			if err != nil {
				return err
			}
			fmt.Printf("✓ Deleted logins: %d\n✓ Deleted vaults: %d\n", deleted["delete_login"], deleted["delete_vault"])
			return nil
		},
	}
	cmd.Flags().StringVarP(&vault.Name, "vault", "v", "", "vault's name")
	cmd.MarkFlagRequired("vault")
	return cmd
}
