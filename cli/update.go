package cli

import "github.com/spf13/cobra"

func NewUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "update a login",
		Long:  "update a login name, password or username from your secret storage",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	return cmd
}
