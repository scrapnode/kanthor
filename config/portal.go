package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/gateway"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
)

type Portal struct {
	Gateway       gateway.Config       `json:"gateway" yaml:"gateway" mapstructure:"gateway" validate:"required"`
	Authenticator authenticator.Config `json:"authenticator" yaml:"authenticator" mapstructure:"authenticator" validate:"required"`
	Authorizator  authorizator.Config  `json:"authorizator" yaml:"authorizator" mapstructure:"authorizator" validate:"required"`
	Cache         cache.Config         `json:"cache" yaml:"cache" mapstructure:"cache" validate:"required"`

	Metrics metric.Config `json:"metrics" yaml:"metrics" mapstructure:"metrics" validate:"-"`
}

func (conf *Portal) Validate() error {
	if err := validator.New().Struct(conf); err != nil {
		return fmt.Errorf("config.Portal: %v", err)
	}

	if err := conf.Gateway.Validate(); err != nil {
		return fmt.Errorf("config.Portal.GRPC: %v", err)
	}
	if err := conf.Authenticator.Validate(); err != nil {
		return fmt.Errorf("config.Portal.Authenticator: %v", err)
	}
	if err := conf.Authorizator.Validate(); err != nil {
		return fmt.Errorf("config.Portal.Enforcer: %v", err)
	}
	if err := conf.Cache.Validate(); err != nil {
		return fmt.Errorf("config.Portal.Cache: %v", err)
	}
	if err := conf.Metrics.Validate(); err != nil {
		return fmt.Errorf("config.Portal.Metrics: %v", err)
	}

	return nil
}
