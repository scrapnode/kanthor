package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/scrapnode/kanthor/infrastructure/configuration"
	"github.com/scrapnode/kanthor/infrastructure/cryptography"
	"github.com/scrapnode/kanthor/infrastructure/database"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/services"
)

func New(provider configuration.Provider) (*Config, error) {
	var conf Config
	return &conf, provider.Unmarshal(&conf)
}

type Config struct {
	Version      string
	Bucket       Bucket              `json:"bucket" yaml:"bucket" mapstructure:"bucket" validate:"required"`
	Cryptography cryptography.Config `json:"cryptography" yaml:"cryptography" mapstructure:"cryptography" validate:"required"`

	Logger   logging.Config  `json:"logger" yaml:"logger" mapstructure:"logger" validate:"required"`
	Database database.Config `json:"database" yaml:"database" mapstructure:"database" validate:"required"`

	Migration  Migration  `json:"migration" yaml:"migration" mapstructure:"migration"`
	PortalApi  PortalApi  `json:"portalapi" yaml:"portalapi" mapstructure:"portalapi"`
	SdkApi     SdkApi     `json:"sdkapi" yaml:"sdkapi" mapstructure:"sdkapi"`
	Scheduler  Scheduler  `json:"scheduler" yaml:"scheduler" mapstructure:"scheduler"`
	Dispatcher Dispatcher `json:"dispatcher" yaml:"dispatcher" mapstructure:"dispatcher"`
}

func (conf *Config) Validate(service string) error {
	if err := validator.New().Struct(conf); err != nil {
		return err
	}

	if err := conf.Cryptography.Validate(); err != nil {
		return fmt.Errorf("config.Cryptography: %v", err)
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

	if service == services.ALL || service == services.MIGRATION {
		return conf.Migration.Validate()
	}
	if service == services.ALL || service == services.PORTAL_API {
		return conf.PortalApi.Validate()
	}
	if service == services.ALL || service == services.SDK_API {
		return conf.SdkApi.Validate()
	}
	if service == services.ALL || service == services.SCHEDULER {
		return conf.Scheduler.Validate()
	}
	if service == services.ALL || service == services.DISPATCHER {
		return conf.Dispatcher.Validate()
	}

	return fmt.Errorf("config: unknow service [%s]", service)
}

type Bucket struct {
	Layout string `json:"layout" yaml:"layout" mapstructure:"layout" validate:"required"`
}

func (conf *Bucket) Validate() error {
	if err := validator.New().Struct(conf); err != nil {
		return err
	}
	return nil
}
