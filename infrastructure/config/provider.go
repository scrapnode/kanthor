package config

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

type Provider interface {
	Unmarshal(dest interface{}) error
}

var Dirs = []string{"$KANTHOR_CONFIG_DIR", "$HOME/.kanthor", ".", "./secrets"}

func New() (Provider, error) {
	v := viper.New()
	v.SetConfigName("configs") // name of config file (without extension)
	v.SetConfigType("yaml")

	for _, dir := range Dirs {
		v.AddConfigPath(dir)
		if err := v.MergeInConfig(); err != nil {
			// ignore not found files, otherwise return error
			if _, notfound := err.(viper.ConfigFileNotFoundError); !notfound {
				panic(fmt.Sprintf("config.viper.MergeInConfig(): %v", err))
			}
		}
	}

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetEnvPrefix("KANTHOR")
	v.AutomaticEnv()

	return &provider{viper: v}, nil
}

type provider struct {
	viper *viper.Viper
}

func (cfg *provider) Unmarshal(dest interface{}) error {
	return cfg.viper.Unmarshal(dest)
}
