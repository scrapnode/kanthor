package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
)

type Scheduler struct {
	Publisher       streaming.PublisherConfig  `json:"publisher" yaml:"publisher" mapstructure:"publisher" validate:"required"`
	Subscriber      streaming.SubscriberConfig `json:"subscriber" yaml:"subscriber" mapstructure:"subscriber" validate:"required"`
	Cache           *cache.Config              `json:"cache" yaml:"cache" validate:"-"`
	ArrangeRequests SchedulerArrangeRequests   `json:"arrange_requests" yaml:"arrange_requests" mapstructure:"arrange_requests" validate:"required"`

	Metrics metric.Config `json:"metrics" yaml:"metrics" mapstructure:"metrics" validate:"-"`
}

func (conf *Scheduler) Validate() error {
	if err := validator.New().Struct(conf); err != nil {
		return fmt.Errorf("config.Dispatcher: %v", err)
	}

	if err := conf.Publisher.Validate(); err != nil {
		return fmt.Errorf("config.Scheduler.Publisher: %v", err)
	}
	if err := conf.Subscriber.Validate(); err != nil {
		return fmt.Errorf("config.Scheduler.Subscriber: %v", err)
	}
	if err := conf.ArrangeRequests.Validate(); err != nil {
		return fmt.Errorf("config.Scheduler.Subscriber: %v", err)
	}
	if conf.Cache != nil {
		if err := conf.Cache.Validate(); err != nil {
			return fmt.Errorf("config.Scheduler.Cache: %v", err)
		}
	}
	if err := conf.Metrics.Validate(); err != nil {
		return fmt.Errorf("config.Scheduler.Metrics: %v", err)
	}

	return nil
}

type SchedulerArrangeRequests struct {
	Concurrency int `json:"concurrency" yaml:"concurrency" mapstructure:"concurrency" validate:"required,gt=0"`
}

func (conf *SchedulerArrangeRequests) Validate() error {
	return validator.New().Struct(conf)
}
