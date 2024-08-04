package cli

import (
	"github.com/mrtnhwtt/kittypass/internal/kittypass"
	"github.com/spf13/cobra"
)

func NewDeleteCmd() *cobra.Command {
	kp := kittypass.New()
	cmd := &cobra.Command{
		Use:   "delete",
		Aliases: []string{"rm", "remove"},
		Short: "delete a saved login",
		Long:  "delete a saved logins from your secret storage using the login name",
		RunE: func(cmd *cobra.Command, args []string) error {
			return kp.Delete()
		},
	}
	cmd.Flags().StringVarP(&kp.Name, "name", "n", "", "Secret's identifier")
	cmd.MarkFlagRequired("name")
	return cmd
}
