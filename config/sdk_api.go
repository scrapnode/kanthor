package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/gateway"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metrics"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
)

type SdkApi struct {
	Gateway      gateway.Config            `json:"gateway" yaml:"gateway" mapstructure:"gateway" validate:"required"`
	Authorizator authorizator.Config       `json:"authorizator" yaml:"authorizator" mapstructure:"authorizator" validate:"required"`
	Publisher    streaming.PublisherConfig `json:"publisher" yaml:"publisher" mapstructure:"publisher" validate:"required"`
	Cache        cache.Config              `json:"cache" yaml:"cache" mapstructure:"cache" validate:"required"`
	Metrics      metrics.Config            `json:"metrics" yaml:"metrics" mapstructure:"metrics" validate:"required"`
}

func (conf *SdkApi) Validate() error {
	if err := validator.New().Struct(conf); err != nil {
		return fmt.Errorf("config.SdkApi: %v", err)
	}

	if err := conf.Gateway.Validate(); err != nil {
		return fmt.Errorf("config.SdkApi.GRPC: %v", err)
	}
	if err := conf.Authorizator.Validate(); err != nil {
		return fmt.Errorf("config.SdkApi.Enforcer: %v", err)
	}
	if err := conf.Cache.Validate(); err != nil {
		return fmt.Errorf("config.SdkApi.Cache: %v", err)
	}
	if err := conf.Metrics.Validate(); err != nil {
		return fmt.Errorf("config.SdkApi.Metrics: %v", err)
	}

	return nil
}
