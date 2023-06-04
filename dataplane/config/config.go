package config

import (
	"github.com/scrapnode/kanthor/infrastructure/config"
	"github.com/scrapnode/kanthor/infrastructure/database"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/msgbroker"
)

type Config struct {
	Logger    *logging.Config   `json:"logger"`
	Database  *database.Config  `json:"database"`
	MsgBroker *msgbroker.Config `json:"msgbroker"`

	Server struct {
		Addr string `json:"addr"`
	} `json:"server"`
}

const Name = "dataplane"

func New(provider config.Provider) (*Config, error) {
	var cfg Config
	err := provider.UnmarshalKey(Name, &cfg)
	return &cfg, err
}
