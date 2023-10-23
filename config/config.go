package config

import (
	"fmt"

	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/circuitbreaker"
	"github.com/scrapnode/kanthor/infrastructure/configuration"
	"github.com/scrapnode/kanthor/infrastructure/cryptography"
	"github.com/scrapnode/kanthor/infrastructure/database"
	"github.com/scrapnode/kanthor/infrastructure/datastore"
	"github.com/scrapnode/kanthor/infrastructure/dlm"
	"github.com/scrapnode/kanthor/infrastructure/idempotency"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/sender"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/validator"
)

// @TODO: mapstructure with env
func New(provider configuration.Provider) (*Config, error) {
	var conf Config
	return &conf, provider.Unmarshal(&conf)
}

type Config struct {
	Version     string
	Development bool `json:"development" yaml:"development" mapstructure:"development"`

	Logger                 logging.Config        `json:"logger" yaml:"logger" mapstructure:"logger"`
	Cryptography           cryptography.Config   `json:"cryptography" yaml:"cryptography" mapstructure:"cryptography"`
	Sender                 sender.Config         `json:"sender" yaml:"sender" mapstructure:"sender"`
	CircuitBreaker         circuitbreaker.Config `json:"circuit_breaker" yaml:"circuit_breaker" mapstructure:"circuit_breaker"`
	Idempotency            idempotency.Config    `json:"idempotency" yaml:"idempotency" mapstructure:"idempotency"`
	DistributedLockManager dlm.Config            `json:"distributed_lock_manager" yaml:"distributed_lock_manager" mapstructure:"distributed_lock_manager"`
	Cache                  cache.Config          `json:"cache" yaml:"cache" mapstructure:"cache"`
	Stream                 streaming.Config      `json:"stream" yaml:"stream" mapstructure:"stream"`
	Metric                 metric.Config         `json:"metric" yaml:"metric" mapstructure:"metric"`
	Authorizator           authorizator.Config   `json:"authorizator" yaml:"authorizator" mapstructure:"authorizator"`

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
	err := validator.Validate(
		validator.DefaultConfig,
		validator.StringRequired("config.version", conf.Version),
	)
	if err != nil {
		return err
	}
	if err := conf.Logger.Validate(); err != nil {
		return fmt.Errorf("config.logger: %v", err)
	}

	// infrastructure
	if err := conf.Cryptography.Validate(); err != nil {
		return fmt.Errorf("config.cryptography: %v", err)
	}
	if err := conf.Sender.Validate(); err != nil {
		return fmt.Errorf("config.sender: %v", err)
	}
	if err := conf.CircuitBreaker.Validate(); err != nil {
		return fmt.Errorf("config.circuit_breaker: %v", err)
	}
	if err := conf.Idempotency.Validate(); err != nil {
		return fmt.Errorf("config.idempotency: %v", err)
	}
	if err := conf.DistributedLockManager.Validate(); err != nil {
		return fmt.Errorf("config.distributed_lock_manager: %v", err)
	}
	if err := conf.Cache.Validate(); err != nil {
		return fmt.Errorf("config.cache: %v", err)
	}
	if err := conf.Stream.Validate(); err != nil {
		return fmt.Errorf("config.stream: %v", err)
	}
	if err := conf.Metric.Validate(); err != nil {
		return fmt.Errorf("config.metric: %v", err)
	}
	if err := conf.Authorizator.Validate(); err != nil {
		return fmt.Errorf("config.authorizator: %v", err)
	}

	// data
	if err := conf.Database.Validate(); err != nil {
		return fmt.Errorf("config.database: %v", err)
	}
	if err := conf.Datastore.Validate(); err != nil {
		return fmt.Errorf("config.datastore: %v", err)
	}

	// services
	if !IsValidServiceName(service) {
		return fmt.Errorf("config: unknown service [%s]", service)
	}

	if service == SERVICE_ALL || service == SERVICE_SDK_API {
		if err := conf.SdkApi.Validate(); err != nil {
			return err
		}
	}
	if service == SERVICE_ALL || service == SERVICE_PORTAL_API {
		if err := conf.PortalApi.Validate(); err != nil {
			return err
		}
	}
	if service == SERVICE_ALL || service == SERVICE_SCHEDULER {
		if err := conf.Scheduler.Validate(); err != nil {
			return err
		}
	}
	if service == SERVICE_ALL || service == SERVICE_DISPATCHER {
		if err := conf.Dispatcher.Validate(); err != nil {
			return err
		}
	}
	if service == SERVICE_ALL || service == SERVICE_STORAGE {
		if err := conf.Storage.Validate(); err != nil {
			return err
		}
	}

	attempt := service == SERVICE_ATTEMPT_TRIGGER_PLANNER ||
		service == SERVICE_ATTEMPT_TRIGGER_EXECUTOR
	if service == SERVICE_ALL || attempt {
		if err := conf.Attempt.Trigger.Validate(); err != nil {
			return err
		}
	}

	return nil
}
