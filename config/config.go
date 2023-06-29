package config

import (
	"github.com/scrapnode/kanthor/infrastructure/configuration"
	"github.com/scrapnode/kanthor/infrastructure/database"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
)

type Config struct {
	Logger    logging.Config             `json:"logger" mapstructure:"logger"`
	Bucket    Bucket                     `json:"bucket" mapstructure:"bucket"`
	Database  database.Config            `json:"database" mapstructure:"database"`
	Streaming streaming.ConnectionConfig `json:"streaming" mapstructure:"streaming"`

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

	conf.Scheduler.Consumer.ConnectionConfig = conf.Streaming
	conf.Dispatcher.Consumer.ConnectionConfig = conf.Streaming

	return &conf, err
}
