package config

import (
	"github.com/scrapnode/kanthor/infrastructure/config"
	"github.com/scrapnode/kanthor/infrastructure/database"
	"github.com/scrapnode/kanthor/infrastructure/logging"
)

type Config struct {
	Database database.Config `json:"database"`

	Migration Migration `json:"migration"`
}

func New(provider config.Provider) (*Config, error) {
	var cfg Config
	err := provider.Unmarshal(&cfg)
	return &cfg, err
}

type Migration struct {
	Logger logging.Config `json:"logger"`
	Source string         `json:"source"`
}
