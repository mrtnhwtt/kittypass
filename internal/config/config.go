package config

import (
	"fmt"


	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func InitializeConfig() *viper.Viper {
	v := viper.New()

	v.SetDefault("encryption", "aes")
	v.SetDefault("storage", "sqlite")
	v.SetDefault("log_path", "$HOME/.local/share/kittypass/logs/")

	v.AddConfigPath("/etc/kittypass")
	v.AddConfigPath("$HOME/.config/kittypass")
	v.AddConfigPath("$XDG_CONFIG_HOME/.config/kittypass")

	v.SetEnvPrefix("KTPS")

	v.AutomaticEnv()
	return v
}

// From https://github.com/carolynvs/stingoftheviper
func BindFlags(cmd *cobra.Command, v *viper.Viper) { // TODO: add this in persistentprerun in the root command to bind all flags for all commands
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		configName := f.Name
		if !f.Changed && v.IsSet(configName) {
			val := v.Get(configName)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}
