package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/enforcer"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
)

type Controlplane struct {
	GRPC          Server               `json:"grpc" yaml:"grpc" mapstructure:"grpc" validate:"required"`
	Cache         *cache.Config        `json:"cache" yaml:"cache" mapstructure:"cache" validate:"-"`
	Authenticator authenticator.Config `json:"authenticator" yaml:"authenticator" mapstructure:"authenticator" validate:"required"`
	Enforcer      enforcer.Config      `json:"enforcer" yaml:"enforcer" mapstructure:"enforcer" validate:"required"`

	Metrics metric.Config `json:"metrics" yaml:"metrics" mapstructure:"metrics" validate:"-"`
}

func (conf *Controlplane) Validate() error {
	if err := validator.New().Struct(conf); err != nil {
		return fmt.Errorf("config.Controlplane: %v", err)
	}

	if err := conf.GRPC.Validate(); err != nil {
		return fmt.Errorf("config.Controlplane.GRPC: %v", err)
	}
	if err := conf.Authenticator.Validate(); err != nil {
		return fmt.Errorf("config.Controlplane.Authenticator: %v", err)
	}
	if err := conf.Enforcer.Validate(); err != nil {
		return fmt.Errorf("config.Controlplane.Enforcer: %v", err)
	}
	if conf.Cache != nil {
		if err := conf.Cache.Validate(); err != nil {
			return fmt.Errorf("config.Controlplane.Cache: %v", err)
		}
	}
	if err := conf.Metrics.Validate(); err != nil {
		return fmt.Errorf("config.Controlplane.Metrics: %v", err)
	}

	return nil
}
