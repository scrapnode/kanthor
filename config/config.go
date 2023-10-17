package config

import (
	"fmt"

	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/circuitbreaker"
	"github.com/scrapnode/kanthor/infrastructure/configuration"
	"github.com/scrapnode/kanthor/infrastructure/coordinator"
	"github.com/scrapnode/kanthor/infrastructure/cryptography"
	"github.com/scrapnode/kanthor/infrastructure/database"
	"github.com/scrapnode/kanthor/infrastructure/datastore"
	"github.com/scrapnode/kanthor/infrastructure/dlm"
	"github.com/scrapnode/kanthor/infrastructure/idempotency"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/pkg/validator"
	"github.com/scrapnode/kanthor/services"
)

// @TODO: mapstructure with env
func New(provider configuration.Provider) (*Config, error) {
	var conf Config
	return &conf, provider.Unmarshal(&conf)
}

type Config struct {
	Version string
	Logger  logging.Config `json:"logger" yaml:"logger" mapstructure:"logger"`

	Cryptography           cryptography.Config   `json:"cryptography" yaml:"cryptography" mapstructure:"cryptography"`
	Idempotency            idempotency.Config    `json:"idempotency" yaml:"idempotency" mapstructure:"idempotency"`
	Coordinator            coordinator.Config    `json:"coordinator" yaml:"coordinator" mapstructure:"coordinator"`
	CircuitBreaker         circuitbreaker.Config `json:"circuit_breaker" yaml:"circuit_breaker" mapstructure:"circuit_breaker"`
	DistributedLockManager dlm.Config            `json:"distributed_lock_manager" yaml:"distributed_lock_manager" mapstructure:"distributed_lock_manager"`
	Metric                 metric.Config         `json:"metric" yaml:"metric" mapstructure:"metric"`
	Cache                  cache.Config          `json:"cache" yaml:"cache" mapstructure:"cache"`

	Database  database.Config  `json:"database" yaml:"database" mapstructure:"database"`
	Datastore datastore.Config `json:"datastore" yaml:"datastore" mapstructure:"datastore"`

	SdkApi     SdkApi     `json:"sdkapi" yaml:"sdkapi" mapstructure:"sdkapi"`
	PortalApi  PortalApi  `json:"portalapi" yaml:"portalapi" mapstructure:"portalapi"`
	Scheduler  Scheduler  `json:"scheduler" yaml:"scheduler" mapstructure:"scheduler"`
	Dispatcher Dispatcher `json:"dispatcher" yaml:"dispatcher" mapstructure:"dispatcher"`
	Storage    Storage    `json:"storage" yaml:"storage" mapstructure:"storage"`
	Attempt    Attempt    `json:"attempt" yaml:"attempt" mapstructure:"attempt"`
}

func (conf *Config) Validate(service string) error {
	err := validator.Validate(validator.DefaultConfig, validator.StringRequired("config.version", conf.Version))
	if err != nil {
		return err
	}
	if err := conf.Logger.Validate(); err != nil {
		return fmt.Errorf("config.logger: %v", err)
	}

	if err := conf.Cryptography.Validate(); err != nil {
		return fmt.Errorf("config.cryptography: %v", err)
	}
	if err := conf.Idempotency.Validate(); err != nil {
		return fmt.Errorf("config.idempotency: %v", err)
	}
	if err := conf.Coordinator.Validate(); err != nil {
		return fmt.Errorf("config.coordinator: %v", err)
	}
	if err := conf.CircuitBreaker.Validate(); err != nil {
		return fmt.Errorf("config.circuit_breaker: %v", err)
	}
	if err := conf.DistributedLockManager.Validate(); err != nil {
		return fmt.Errorf("config.distributed_lock_manager: %v", err)
	}
	if err := conf.Metric.Validate(); err != nil {
		return fmt.Errorf("config.metric: %v", err)
	}
	if err := conf.Cache.Validate(); err != nil {
		return fmt.Errorf("config.cache: %v", err)
	}
	if err := conf.Database.Validate(); err != nil {
		return fmt.Errorf("config.database: %v", err)
	}
	if err := conf.Datastore.Validate(); err != nil {
		return fmt.Errorf("config.datastore: %v", err)
	}

	if !services.IsValidServiceName(service) {
		return fmt.Errorf("config: unknown service [%s]", service)
	}

	if service == services.SERVICE_ALL || service == services.SERVICE_SDK_API {
		if err := conf.SdkApi.Validate(); err != nil {
			return err
		}
	}
	if service == services.SERVICE_ALL || service == services.SERVICE_PORTAL_API {
		if err := conf.PortalApi.Validate(); err != nil {
			return err
		}
	}
	if service == services.SERVICE_ALL || service == services.SERVICE_SCHEDULER {
		if err := conf.Scheduler.Validate(); err != nil {
			return err
		}
	}
	if service == services.SERVICE_ALL || service == services.SERVICE_DISPATCHER {
		if err := conf.Dispatcher.Validate(); err != nil {
			return err
		}
	}
	if service == services.SERVICE_ALL || service == services.SERVICE_STORAGE {
		if err := conf.Storage.Validate(); err != nil {
			return err
		}
	}
	if service == services.SERVICE_ALL || service == services.SERVICE_ATTEMPT {
		if err := conf.Attempt.Validate(); err != nil {
			return err
		}
	}

	return nil
}
