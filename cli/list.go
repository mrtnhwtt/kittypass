package cli

import (
	"fmt"

	"github.com/mrtnhwtt/kittypass/internal/kittypass"
	"github.com/spf13/cobra"
)

func NewListCmd() *cobra.Command {
	kp := kittypass.New()
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list saved logins",
		Long:  "list saved logins from your secret storage using the login name or username",
		RunE: func(cmd *cobra.Command, args []string) error {
			loginList, err := kp.List()
			if err != nil {
				return err
			}
			for _, login := range loginList {
				fmt.Println("---------------------------------------")
				fmt.Printf("Login Name: %s\nUsername: %s\n", login["name"], login["username"])
			}
			fmt.Println("---------------------------------------")
			return nil
		},
	}
	cmd.Flags().StringVarP(&kp.Name, "name", "n", "", "Secret's identifier")
	cmd.Flags().StringVarP(&kp.Username, "username", "u", "", "Username or email associated with the password")
	cmd.MarkFlagsOneRequired("name", "username")
	return cmd
}
