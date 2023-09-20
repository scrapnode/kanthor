package config

import (
	"fmt"

	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/gateway"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type SdkApi struct {
	Gateway      gateway.Config            `json:"gateway" yaml:"gateway" mapstructure:"gateway"`
	Authorizator authorizator.Config       `json:"authorizator" yaml:"authorizator" mapstructure:"authorizator"`
	Publisher    streaming.PublisherConfig `json:"publisher" yaml:"publisher" mapstructure:"publisher"`
	Cache        cache.Config              `json:"cache" yaml:"cache" mapstructure:"cache"`
	Metrics      metric.Config             `json:"metrics" yaml:"metrics" mapstructure:"metrics"`

	PortalConnection *SdkApiPortalConnection `json:"portal_connection" yaml:"portal_connection" mapstructure:"portal_connection"`
}

func (conf *SdkApi) Validate() error {
	if err := conf.Gateway.Validate(); err != nil {
		return fmt.Errorf("config.sdkapi.gateway: %v", err)
	}
	if err := conf.Authorizator.Validate(); err != nil {
		return fmt.Errorf("config.sdkapi.enforcer: %v", err)
	}
	if err := conf.Cache.Validate(); err != nil {
		return fmt.Errorf("config.sdkapi.cache: %v", err)
	}
	if err := conf.Metrics.Validate(); err != nil {
		return fmt.Errorf("config.sdkapi.metrics: %v", err)
	}

	if conf.PortalConnection != nil {
		if err := conf.PortalConnection.Validate(); err != nil {
			return fmt.Errorf("config.sdkapi.portal_connection: %v", err)
		}
	}

	return nil
}

type SdkApiPortalConnection struct {
	Account string `json:"account" yaml:"account" mapstructure:"account"`
}

func (conf *SdkApiPortalConnection) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringUri("config.sdkapi.portal_connection.account", conf.Account),
	)
}
