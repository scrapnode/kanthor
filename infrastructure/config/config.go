package config

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

type Provider interface {
	UnmarshalKey(name string, dest interface{}) error
}

var DefaultDirs = []string{"./secrets", ".", "$HOME/.kanthor"}

func New() Provider {
	provider := viper.New()
	provider.SetConfigName("configs") // name of config file (without extension)
	provider.SetConfigType("yaml")

	for _, dir := range DefaultDirs {
		provider.AddConfigPath(dir)
		if err := provider.MergeInConfig(); err != nil {
			// ignore not found files, otherwise return error
			if _, notfound := err.(viper.ConfigFileNotFoundError); !notfound {
				panic(fmt.Sprintf("config.provider.MergeInConfig(): %v", err))
			}
		}
	}

	provider.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	provider.SetEnvPrefix("KANTHOR")
	provider.AutomaticEnv()

	return &config{provider: provider}
}

type config struct {
	provider *viper.Viper
}

func (cfg *config) UnmarshalKey(name string, dest interface{}) error {
	return cfg.provider.UnmarshalKey(name, dest)
}
