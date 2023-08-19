package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/gateway"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metrics"
)

type PortalApi struct {
	Gateway       gateway.Config       `json:"gateway" yaml:"gateway" mapstructure:"gateway" validate:"required"`
	Authenticator authenticator.Config `json:"authenticator" yaml:"authenticator" mapstructure:"authenticator" validate:"required"`
	Authorizator  authorizator.Config  `json:"authorizator" yaml:"authorizator" mapstructure:"authorizator" validate:"required"`
	Cache         cache.Config         `json:"cache" yaml:"cache" mapstructure:"cache" validate:"required"`
	Metrics       metrics.Config       `json:"metrics" yaml:"metrics" mapstructure:"metrics" validate:"required"`
}

func (conf *PortalApi) Validate() error {
	if err := validator.New().Struct(conf); err != nil {
		return fmt.Errorf("config.PortalApi: %v", err)
	}

	if err := conf.Gateway.Validate(); err != nil {
		return fmt.Errorf("config.PortalApi.GRPC: %v", err)
	}
	if err := conf.Authenticator.Validate(); err != nil {
		return fmt.Errorf("config.PortalApi.Authenticator: %v", err)
	}
	if err := conf.Authorizator.Validate(); err != nil {
		return fmt.Errorf("config.PortalApi.Enforcer: %v", err)
	}
	if err := conf.Cache.Validate(); err != nil {
		return fmt.Errorf("config.PortalApi.Cache: %v", err)
	}
	if err := conf.Metrics.Validate(); err != nil {
		return fmt.Errorf("config.PortalApi.Metrics: %v", err)
	}

	return nil
}
