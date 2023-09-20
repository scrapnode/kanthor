package config

import (
	"fmt"

	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/gateway"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
)

type PortalApi struct {
	Gateway       gateway.Config       `json:"gateway" yaml:"gateway" mapstructure:"gateway"`
	Authenticator authenticator.Config `json:"authenticator" yaml:"authenticator" mapstructure:"authenticator"`
	Authorizator  authorizator.Config  `json:"authorizator" yaml:"authorizator" mapstructure:"authorizator"`
	Cache         cache.Config         `json:"cache" yaml:"cache" mapstructure:"cache"`
	Metrics       metric.Config        `json:"metrics" yaml:"metrics" mapstructure:"metrics"`
}

func (conf *PortalApi) Validate() error {
	if err := conf.Gateway.Validate(); err != nil {
		return fmt.Errorf("config.portalapi.gateway: %v", err)
	}
	if err := conf.Authenticator.Validate(); err != nil {
		return fmt.Errorf("config.portalapi.authenticator: %v", err)
	}
	if err := conf.Authorizator.Validate(); err != nil {
		return fmt.Errorf("config.portalapi.enforcer: %v", err)
	}
	if err := conf.Cache.Validate(); err != nil {
		return fmt.Errorf("config.portalapi.cache: %v", err)
	}
	if err := conf.Metrics.Validate(); err != nil {
		return fmt.Errorf("config.portalapi.metrics: %v", err)
	}

	return nil
}
