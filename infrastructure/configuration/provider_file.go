package configuration

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

var FileLookingDirs = []string{"$KANTHOR_CONFIG_DIR", "$HOME/.kanthor", "."}

func NewFile(dirs []string) (Provider, error) {
	v := viper.New()
	v.SetConfigName("configs") // name of config file (without extension)
	v.SetConfigType("yaml")

	for _, dir := range dirs {
		v.AddConfigPath(dir)
		if err := v.MergeInConfig(); err != nil {
			// ignore not found files, otherwise return error
			if _, notfound := err.(viper.ConfigFileNotFoundError); !notfound {
				return nil, fmt.Errorf("config.viper.MergeInConfig(): %v", err)
			}
		}
	}

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetEnvPrefix("KANTHOR")
	v.AutomaticEnv()

	return &file{viper: v}, nil
}

type file struct {
	viper *viper.Viper
}

func (provider *file) Unmarshal(dest interface{}) error {
	return provider.viper.Unmarshal(dest)
}
