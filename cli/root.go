package cli

import (
	"github.com/fatih/color"
	"github.com/mrtnhwtt/kittypass/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	green   = color.New(color.FgGreen).SprintFunc()
	red     = color.New(color.FgRed).SprintFunc()
	blue    = color.New(color.FgBlue).SprintFunc()
	magenta = color.New(color.FgMagenta).SprintFunc()
)

func NewRootCmd(conf *viper.Viper) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "kittypass",
		Short:        "A CLI password manager",
		Long:         "A CLI password manager to securely stock your password.",
		SilenceUsage: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			config.BindFlags(cmd, conf)
		},
	}

	cmd.AddCommand(
		NewAddCmd(),
		NewGetCmd(),
		NewListCmd(),
		NewDeleteCmd(),
		NewUpdateCmd(),
	)
	// TODO: implement a migration command to migrate a vault between different storage.

	return cmd
}
