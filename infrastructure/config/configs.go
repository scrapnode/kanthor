package config

import (
	"fmt"

	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/circuitbreaker"
	"github.com/scrapnode/kanthor/infrastructure/configuration"
	"github.com/scrapnode/kanthor/infrastructure/cryptography"
	"github.com/scrapnode/kanthor/infrastructure/dlm"
	"github.com/scrapnode/kanthor/infrastructure/idempotency"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/sender"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
)

// @TODO: mapstructure with env
func New(provider configuration.Provider) (*Config, error) {
	var conf Config
	return &conf, provider.Unmarshal(&conf)
}

type Config struct {
	Cryptography           cryptography.Config   `json:"cryptography" yaml:"cryptography" mapstructure:"cryptography"`
	Sender                 sender.Config         `json:"sender" yaml:"sender" mapstructure:"sender"`
	CircuitBreaker         circuitbreaker.Config `json:"circuit_breaker" yaml:"circuit_breaker" mapstructure:"circuit_breaker"`
	Idempotency            idempotency.Config    `json:"idempotency" yaml:"idempotency" mapstructure:"idempotency"`
	DistributedLockManager dlm.Config            `json:"distributed_lock_manager" yaml:"distributed_lock_manager" mapstructure:"distributed_lock_manager"`
	Cache                  cache.Config          `json:"cache" yaml:"cache" mapstructure:"cache"`
	Stream                 streaming.Config      `json:"stream" yaml:"stream" mapstructure:"stream"`
	Metric                 metric.Config         `json:"metric" yaml:"metric" mapstructure:"metric"`
	Authorizator           authorizator.Config   `json:"authorizator" yaml:"authorizator" mapstructure:"authorizator"`
}

func (conf *Config) Validate() error {
	if err := conf.Cryptography.Validate(); err != nil {
		return fmt.Errorf("infrastructure.cryptography: %v", err)
	}
	if err := conf.Sender.Validate(); err != nil {
		return fmt.Errorf("infrastructure.sender: %v", err)
	}
	if err := conf.CircuitBreaker.Validate(); err != nil {
		return fmt.Errorf("infrastructure.circuit_breaker: %v", err)
	}
	if err := conf.Idempotency.Validate(); err != nil {
		return fmt.Errorf("infrastructure.idempotency: %v", err)
	}
	if err := conf.DistributedLockManager.Validate(); err != nil {
		return fmt.Errorf("infrastructure.distributed_lock_manager: %v", err)
	}
	if err := conf.Cache.Validate(); err != nil {
		return fmt.Errorf("infrastructure.cache: %v", err)
	}
	if err := conf.Stream.Validate(); err != nil {
		return fmt.Errorf("infrastructure.stream: %v", err)
	}
	if err := conf.Metric.Validate(); err != nil {
		return fmt.Errorf("infrastructure.metric: %v", err)
	}
	if err := conf.Authorizator.Validate(); err != nil {
		return fmt.Errorf("infrastructure.authorizator: %v", err)
	}

	return nil
}
