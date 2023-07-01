package config

import (
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/configuration"
	"github.com/scrapnode/kanthor/infrastructure/database"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
)

type Config struct {
	Bucket Bucket `json:"bucket" mapstructure:"bucket"`

	Logger    logging.Config             `json:"logger" mapstructure:"logger"`
	Database  database.Config            `json:"database" mapstructure:"database"`
	Streaming streaming.ConnectionConfig `json:"streaming" mapstructure:"streaming"`
	Cache     cache.Config               `json:"cache" mapstructure:"cache"`

	Migration  Migration  `json:"migration" mapstructure:"migration"`
	Dataplane  Dataplane  `json:"dataplane" mapstructure:"dataplane"`
	Scheduler  Scheduler  `json:"scheduler" mapstructure:"scheduler"`
	Dispatcher Dispatcher `json:"dispatcher" mapstructure:"dispatcher"`
}

type Bucket struct {
	Layout string `json:"layout" mapstructure:"layout"`
}

type Server struct {
	Addr string `json:"addr" mapstructure:"addr"`
}

func New(provider configuration.Provider) (*Config, error) {
	var conf Config
	err := provider.Unmarshal(&conf)

	return &conf, err
}
