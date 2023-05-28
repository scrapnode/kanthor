package logging

import "github.com/scrapnode/kanthor/infrastructure/config"

type Config struct {
	Debug bool
	Level string
}

const ConfigName = "logger"

func GetConfig(provider config.Provider) (*Config, error) {
	var cfg Config
	err := provider.UnmarshalKey(ConfigName, &cfg)
	return &cfg, err
}
