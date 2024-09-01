package config

import (
	"fmt"


	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func InitializeConfig() error {
	v := viper.New()
	
	v.AddConfigPath("/etc/kittypass")
	v.AddConfigPath("$HOME/.config/kittypass")
	v.AddConfigPath("$XDG_CONFIG_HOME/.config/kittypass")
	v.SetEnvPrefix("KTPS")
	v.AutomaticEnv()
	return nil
}

// From https://github.com/carolynvs/stingoftheviper
func BindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		configName := f.Name

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && v.IsSet(configName) {
			val := v.Get(configName)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}
