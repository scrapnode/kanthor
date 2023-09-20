package config

import (
	"fmt"

	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type Scheduler struct {
	Publisher  streaming.PublisherConfig  `json:"publisher" yaml:"publisher" mapstructure:"publisher"`
	Subscriber streaming.SubscriberConfig `json:"subscriber" yaml:"subscriber" mapstructure:"subscriber"`
	Cache      cache.Config               `json:"cache" yaml:"cache"`
	Request    Request                    `json:"request" yaml:"request" mapstructure:"request"`
	Metrics    metric.Config              `json:"metrics" yaml:"metrics" mapstructure:"metrics"`
}

func (conf *Scheduler) Validate() error {
	if err := conf.Publisher.Validate(); err != nil {
		return fmt.Errorf("config.scheduler.publisher: %v", err)
	}
	if err := conf.Subscriber.Validate(); err != nil {
		return fmt.Errorf("config.scheduler.subscriber: %v", err)
	}
	if err := conf.Request.Validate(); err != nil {
		return fmt.Errorf("config.scheduler.request: %v", err)
	}
	if err := conf.Cache.Validate(); err != nil {
		return fmt.Errorf("config.scheduler.cache: %v", err)
	}
	if err := conf.Metrics.Validate(); err != nil {
		return fmt.Errorf("config.scheduler.metrics: %v", err)
	}

	return nil
}

type Request struct {
	Arrange RequestArrange `json:"arrange" yaml:"arrange" mapstructure:"arrange"`
}

func (conf *Request) Validate() error {
	if err := conf.Arrange.Validate(); err != nil {
		return err
	}

	return nil
}

type RequestArrange struct {
	Concurrency int `json:"concurrency" yaml:"concurrency" mapstructure:"concurrency"`
}

func (conf *RequestArrange) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThan("config.request.arrange.concurrency", conf.Concurrency, 0),
	)
}
