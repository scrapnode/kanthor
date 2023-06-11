package config

import (
	"github.com/scrapnode/kanthor/infrastructure/config"
	"github.com/scrapnode/kanthor/infrastructure/database"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
)

type Config struct {
	Bucket    config.Bucket
	Database  database.Config            `json:"database"`
	Streaming streaming.ConnectionConfig `json:"streaming_publisher"`

	Dataplane *Dataplane `json:"dataplane"`
}

type Dataplane struct {
	Logger logging.Config `json:"logger"`
	Server struct {
		Addr string `json:"addr"`
	} `json:"server"`
}

func New(provider config.Provider) (*Config, error) {
	var cfg Config
	err := provider.Unmarshal(&cfg)
	return &cfg, err
}
