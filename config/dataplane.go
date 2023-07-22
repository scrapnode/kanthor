package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/gateway"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
)

type Dataplane struct {
	Gateway       gateway.Config            `json:"gateway" yaml:"gateway" mapstructure:"gateway" validate:"required"`
	Publisher     streaming.PublisherConfig `json:"publisher" yaml:"publisher" mapstructure:"publisher" validate:"required"`
	Cache         cache.Config              `json:"cache" yaml:"cache" mapstructure:"cache" validate:"required"`
	Authenticator authenticator.Config      `json:"authenticator" yaml:"authenticator" mapstructure:"authenticator" validate:"required"`
	Authorizator  authorizator.Config       `json:"authorizator" yaml:"authorizator" mapstructure:"authorizator" validate:"required"`

	Metrics metric.Config `json:"metrics" yaml:"metrics" mapstructure:"metrics" validate:"-"`
}

func (conf *Dataplane) Validate() error {
	if err := validator.New().Struct(conf); err != nil {
		return fmt.Errorf("config.Dataplane: %v", err)
	}

	if err := conf.Gateway.Validate(); err != nil {
		return fmt.Errorf("config.Dataplane.GRPC: %v", err)
	}
	if err := conf.Publisher.Validate(); err != nil {
		return fmt.Errorf("config.Dataplane.Publisher: %v", err)
	}
	if err := conf.Authenticator.Validate(); err != nil {
		return fmt.Errorf("config.Dataplane.Authenticator: %v", err)
	}
	if err := conf.Cache.Validate(); err != nil {
		return fmt.Errorf("config.Dataplane.Cache: %v", err)
	}
	if err := conf.Metrics.Validate(); err != nil {
		return fmt.Errorf("config.Dataplane.Metrics: %v", err)
	}

	return nil
}
