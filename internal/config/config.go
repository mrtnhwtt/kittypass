package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func InitializeConfig() *viper.Viper {
	v := viper.New()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	v.SetDefault("encryption", "aes") // TODO: this should do something once another encryption algo is implemented
	v.SetDefault("storage", "sqlite") // TODO: this should do something once another storage option is implemented
	v.SetDefault("log_path", filepath.Join(homeDir, ".local/share/kittypass/logs/kittypass.log"))

	v.AddConfigPath(filepath.Join(homeDir, ".config/kittypass"))
	v.AddConfigPath("/etc/kittypass")

	v.SetEnvPrefix("KTPS")

	v.AutomaticEnv()
	return v
}

// From https://github.com/carolynvs/stingoftheviper
func BindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		configName := f.Name
		if !f.Changed && v.IsSet(configName) {
			val := v.Get(configName)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}
