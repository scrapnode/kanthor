package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
)

type Storage struct {
	Subscriber streaming.SubscriberConfig `json:"subscriber" yaml:"subscriber" mapstructure:"subscriber" validate:"required"`
}

func (conf *Storage) Validate() error {
	if err := validator.New().Struct(conf); err != nil {
		return fmt.Errorf("config.Dispatcher: %v", err)
	}

	if err := conf.Subscriber.Validate(); err != nil {
		return fmt.Errorf("config.Scheduler.Subscriber: %v", err)
	}

	return nil
}
