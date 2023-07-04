package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
)

type Dataplane struct {
	GRPC      Server                    `json:"grpc" mapstructure:"grpc" validate:"required"`
	Publisher streaming.PublisherConfig `json:"publisher" mapstructure:"publisher" validate:"required"`
	Cache     *cache.Config             `json:"cache" mapstructure:"cache" validate:"-"`

	Metrics metric.Config `json:"metrics" mapstructure:"metrics" validate:"-"`
}

func (conf Dataplane) Validate() error {
	if err := validator.New().Struct(conf); err != nil {
		return fmt.Errorf("config.Dataplane: %v", err)
	}

	if err := conf.GRPC.Validate(); err != nil {
		return fmt.Errorf("config.Dataplane.GRPC: %v", err)
	}
	if err := conf.Publisher.Validate(); err != nil {
		return fmt.Errorf("config.Dataplane.Publisher: %v", err)
	}
	if conf.Cache != nil {
		if err := conf.Cache.Validate(); err != nil {
			return fmt.Errorf("config.Dataplane.Cache: %v", err)
		}
	}
	if err := conf.Metrics.Validate(); err != nil {
		return fmt.Errorf("config.Dataplane.Metrics: %v", err)
	}

	return nil
}
