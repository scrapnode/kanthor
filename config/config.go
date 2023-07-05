package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/configuration"
	"github.com/scrapnode/kanthor/infrastructure/database"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/services"
)

type Config struct {
	Bucket Bucket `json:"bucket" mapstructure:"bucket" validate:"required"`

	Logger    logging.Config             `json:"logger" mapstructure:"logger" validate:"required"`
	Streaming streaming.ConnectionConfig `json:"streaming" mapstructure:"streaming" validate:"required"`
	Database  database.Config            `json:"database" mapstructure:"database" validate:"required"`
	Cache     cache.Config               `json:"cache" mapstructure:"cache" validate:"required"`

	Migration  Migration  `json:"migration" mapstructure:"migration"`
	Dataplane  Dataplane  `json:"dataplane" mapstructure:"dataplane"`
	Scheduler  Scheduler  `json:"scheduler" mapstructure:"scheduler"`
	Dispatcher Dispatcher `json:"dispatcher" mapstructure:"dispatcher"`
}

func (conf Config) Validate(service string) error {
	if err := validator.New().Struct(conf); err != nil {
		return err
	}

	if err := conf.Bucket.Validate(); err != nil {
		return fmt.Errorf("config.Bucket: %v", err)
	}
	if err := conf.Logger.Validate(); err != nil {
		return fmt.Errorf("config.Logger: %v", err)
	}
	if err := conf.Database.Validate(); err != nil {
		return fmt.Errorf("config.Database: %v", err)
	}
	if err := conf.Streaming.Validate(); err != nil {
		return fmt.Errorf("config.Streaming: %v", err)
	}
	if err := conf.Cache.Validate(); err != nil {
		return fmt.Errorf("config.Cache: %v", err)
	}

	if service == services.MIGRATION {
		return conf.Migration.Validate()
	}
	if service == services.DATAPLANE {
		return conf.Dataplane.Validate()
	}
	if service == services.SCHEDULER {
		return conf.Scheduler.Validate()
	}
	if service == services.DISPATCHER {
		return conf.Dispatcher.Validate()
	}

	return fmt.Errorf("config: unknow service [%s]", service)
}

type Bucket struct {
	Layout string `json:"layout" mapstructure:"layout" validate:"required"`
}

func (conf Bucket) Validate() error {
	return validator.New().Struct(conf)
}

type Server struct {
	Addr string `json:"addr" mapstructure:"addr" validate:"required"`
}

func (conf Server) Validate() error {
	return validator.New().Struct(conf)
}

func New(provider configuration.Provider) (*Config, error) {
	var conf Config
	return &conf, provider.Unmarshal(&conf)
}
