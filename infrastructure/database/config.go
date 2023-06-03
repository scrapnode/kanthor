package database

import "github.com/scrapnode/kanthor/infrastructure/config"

type Config struct {
	Uri           string `json:"uri"`
	RetryAttempts uint   `json:"retry_attempts"`
}

const ConfigName = "database"

func GetConfig(provider config.Provider) (*Config, error) {
	var cfg Config
	err := provider.UnmarshalKey(ConfigName, &cfg)
	return &cfg, err
}
