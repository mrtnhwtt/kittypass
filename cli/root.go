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
		// Using PersistentPreRun will run this code on each command, binding the flags to the config.
		// If a child command implements PersistentPreRun, it will override this one, so it needs to run BindFlags too
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			config.BindFlags(cmd, conf)
		},
		// TODO: TUI implementation here. Calling a subcommand won't launch the TUI, they will always be imperative
		// RunE: func(cmd *cobra.Command, args []string) error { return nil },
	}

	cmd.AddCommand(
		NewAddCmd(conf),
		NewGetCmd(conf),
		NewListCmd(conf),
		NewDeleteCmd(conf),
		NewUpdateCmd(conf),
	)
	// TODO: implement a migration command to migrate a vault between different storage.

	return cmd
}
