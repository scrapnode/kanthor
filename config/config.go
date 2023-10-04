package config

import (
	"fmt"

	"github.com/scrapnode/kanthor/infrastructure/configuration"
	"github.com/scrapnode/kanthor/infrastructure/coordinator"
	"github.com/scrapnode/kanthor/infrastructure/cryptography"
	"github.com/scrapnode/kanthor/infrastructure/database"
	"github.com/scrapnode/kanthor/infrastructure/datastore"
	"github.com/scrapnode/kanthor/infrastructure/idempotency"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/pkg/validator"
	"github.com/scrapnode/kanthor/services"
)

func New(provider configuration.Provider) (*Config, error) {
	var conf Config
	return &conf, provider.Unmarshal(&conf)
}

type Config struct {
	Version      string
	Cryptography cryptography.Config `json:"cryptography" yaml:"cryptography" mapstructure:"cryptography"`

	Logger      logging.Config     `json:"logger" yaml:"logger" mapstructure:"logger"`
	Database    database.Config    `json:"database" yaml:"database" mapstructure:"database"`
	Datastore   datastore.Config   `json:"datastore" yaml:"datastore" mapstructure:"datastore"`
	Idempotency idempotency.Config `json:"idempotency" yaml:"idempotency" mapstructure:"idempotency"`
	Coordinator coordinator.Config `json:"coordinator" yaml:"coordinator" mapstructure:"coordinator"`

	SdkApi     SdkApi     `json:"sdkapi" yaml:"sdkapi" mapstructure:"sdkapi"`
	PortalApi  PortalApi  `json:"portalapi" yaml:"portalapi" mapstructure:"portalapi"`
	Scheduler  Scheduler  `json:"scheduler" yaml:"scheduler" mapstructure:"scheduler"`
	Dispatcher Dispatcher `json:"dispatcher" yaml:"dispatcher" mapstructure:"dispatcher"`
	Storage    Storage    `json:"storage" yaml:"storage" mapstructure:"storage"`
}

func (conf *Config) Validate(service string) error {
	err := validator.Validate(validator.DefaultConfig, validator.StringRequired("config.version", conf.Version))
	if err != nil {
		return err
	}

	if err := conf.Cryptography.Validate(); err != nil {
		return fmt.Errorf("config.cryptography: %v", err)
	}
	if err := conf.Logger.Validate(); err != nil {
		return fmt.Errorf("config.logger: %v", err)
	}
	if err := conf.Database.Validate(); err != nil {
		return fmt.Errorf("config.database: %v", err)
	}
	if err := conf.Datastore.Validate(); err != nil {
		return fmt.Errorf("config.datastore: %v", err)
	}
	if err := conf.Idempotency.Validate(); err != nil {
		return fmt.Errorf("config.idempotency: %v", err)
	}
	if err := conf.Coordinator.Validate(); err != nil {
		return fmt.Errorf("config.coordinator: %v", err)
	}

	if !services.Valid(service) {
		return fmt.Errorf("config: unknown service [%s]", service)
	}

	if service == services.ALL || service == services.SDK_API {
		if err := conf.SdkApi.Validate(); err != nil {
			return err
		}
	}
	if service == services.ALL || service == services.PORTAL_API {
		if err := conf.PortalApi.Validate(); err != nil {
			return err
		}
	}
	if service == services.ALL || service == services.SCHEDULER {
		if err := conf.Scheduler.Validate(); err != nil {
			return err
		}
	}
	if service == services.ALL || service == services.DISPATCHER {
		if err := conf.Dispatcher.Validate(); err != nil {
			return err
		}
	}
	if service == services.ALL || service == services.STORAGE {
		if err := conf.Storage.Validate(); err != nil {
			return err
		}
	}

	return nil
}
