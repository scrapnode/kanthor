package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metrics"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
)

type Storage struct {
	Subscriber streaming.SubscriberConfig `json:"subscriber" yaml:"subscriber" mapstructure:"subscriber" validate:"required"`
	Metrics    metrics.Config             `json:"metrics" yaml:"metrics" mapstructure:"metrics" validate:"required"`
}

func (conf *Storage) Validate() error {
	if err := validator.New().Struct(conf); err != nil {
		return fmt.Errorf("config.Storage: %v", err)
	}

	if err := conf.Subscriber.Validate(); err != nil {
		return fmt.Errorf("config.Storage.Subscriber: %v", err)
	}
	if err := conf.Metrics.Validate(); err != nil {
		return fmt.Errorf("config.Storage.Subscriber: %v", err)
	}

	return nil
}
