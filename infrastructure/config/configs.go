package config

import (
	"log"

	"github.com/scrapnode/kanthor/configuration"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/circuitbreaker"
	"github.com/scrapnode/kanthor/infrastructure/dlm"
	"github.com/scrapnode/kanthor/infrastructure/idempotency"
	"github.com/scrapnode/kanthor/infrastructure/sender"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/utils"
)

func New(provider configuration.Provider) (*Config, error) {
	var conf Wrapper
	if err := provider.Unmarshal(&conf); err != nil {
		return nil, err
	}
	if err := conf.Validate(); err != nil {
		log.Println(string(utils.Stringify(conf)))
		return nil, err
	}

	return &conf.Infrastructure, nil
}

type Wrapper struct {
	Infrastructure Config `json:"infrastructure" yaml:"infrastructure" mapstructure:"infrastructure"`
}

func (conf *Wrapper) Validate() error {
	if err := conf.Infrastructure.Validate(); err != nil {
		return err
	}
	return nil
}

type Config struct {
	Sender                 sender.Config          `json:"sender" yaml:"sender" mapstructure:"sender"`
	CircuitBreaker         circuitbreaker.Config  `json:"circuit_breaker" yaml:"circuit_breaker" mapstructure:"circuit_breaker"`
	Idempotency            idempotency.Config     `json:"idempotency" yaml:"idempotency" mapstructure:"idempotency"`
	DistributedLockManager dlm.Config             `json:"distributed_lock_manager" yaml:"distributed_lock_manager" mapstructure:"distributed_lock_manager"`
	Cache                  cache.Config           `json:"cache" yaml:"cache" mapstructure:"cache"`
	Authenticators         []authenticator.Config `json:"authenticators" yaml:"authenticators" mapstructure:"authenticators"`
	Authorizator           authorizator.Config    `json:"authorizator" yaml:"authorizator" mapstructure:"authorizator"`
	Streaming              streaming.Config       `json:"streaming" yaml:"streaming" mapstructure:"streaming"`
}

func (conf *Config) Validate() error {
	if err := conf.Sender.Validate(); err != nil {
		return err
	}
	if err := conf.CircuitBreaker.Validate(); err != nil {
		return err
	}
	if err := conf.Idempotency.Validate(); err != nil {
		return err
	}
	if err := conf.DistributedLockManager.Validate(); err != nil {
		return err
	}
	if err := conf.Cache.Validate(); err != nil {
		return err
	}
	if len(conf.Authenticators) > 0 {
		for _, authenticator := range conf.Authenticators {
			if err := authenticator.Validate(); err != nil {
				return err
			}
		}
	}
	if err := conf.Authorizator.Validate(); err != nil {
		return err
	}
	if err := conf.Streaming.Validate(); err != nil {
		return err
	}

	return nil
}
