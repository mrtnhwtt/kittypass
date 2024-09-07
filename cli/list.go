package cli

import (
	"fmt"

	"github.com/mrtnhwtt/kittypass/internal/kittypass"
	"github.com/mrtnhwtt/kittypass/internal/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewListCmd(conf *viper.Viper) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "list vaults or logins",
		Long:    "list vaults or logins",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	cmd.AddCommand(
		NewListLoginCmd(conf),
		NewListVaultCmd(conf),
	)
	return cmd
}

func NewListLoginCmd(conf *viper.Viper) *cobra.Command {
	login := kittypass.NewLogin()
	vault := kittypass.NewVault()
	login.Vault = &vault

	cmd := &cobra.Command{
		Use:     "login",
		Aliases: []string{"passwords", "pass", "logins", "password"},
		Short:   "lists logins",
		Long:    "lists logins. Search for login from login name, username or email. Limit search to a specific vault",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return login.Vault.OpenDbConnection(conf)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if login.Vault.Name != "" {
				err := login.Vault.Get()
				if err != nil {
					return err
				}
			}
			loginList, err := login.List()
			if err != nil {
				return err
			}
			if len(loginList) < 1 {
				fmt.Println(red("No matching logins"))
				return nil
			}
			for _, login := range loginList {
				formattedTime, err := utils.ParseTimestamp(login["timestamp"])
				if err != nil {
					formattedTime = "unknown"
				}
				fmt.Println("------------------------------------------------------------------------------")
				fmt.Printf("Vault: %s\nLogin Name: %s\nUsername: %s\nCreated: %s\n", login["vault_name"], login["name"], login["username"], formattedTime)
			}
			fmt.Println("------------------------------------------------------------------------------")
			return nil
		},
	}
	cmd.Flags().StringVarP(&login.Name, "name", "n", "", "search for a login name")
	cmd.Flags().StringVarP(&login.Username, "username", "u", "", "search for a username or email associated with the password")
	cmd.Flags().StringVarP(&login.Vault.Name, "vault", "v", "", "limit search to a vault")
	return cmd
}

func NewListVaultCmd(conf *viper.Viper) *cobra.Command {
	vault := kittypass.NewVault()
	cmd := &cobra.Command{
		Use:     "vault",
		Aliases: []string{"folder", "vaults", "folders"},
		Short:   "lists vaults",
		Long:    "lists vaults. Search for a vault by providing a name.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return vault.OpenDbConnection(conf)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			vaultList, err := vault.List()
			if err != nil {
				return nil
			}
			if len(vaultList) < 1 {
				fmt.Println(red("No matching vaults"))
			}
			for _, vault := range vaultList {
				formattedTime, err := utils.ParseTimestamp(vault["date_created"])
				if err != nil {
					formattedTime = "unknown"
				}
				fmt.Println("---------------------------------------")
				fmt.Printf("Vault Name: %s\nDescription: %s\nCreation Date: %s\n", vault["name"], vault["description"], formattedTime)
			}
			fmt.Println("---------------------------------------")
			return nil
		},
	}
	cmd.Flags().StringVarP(&vault.Name, "name", "n", "", "search for a vault name")
	return cmd
}
