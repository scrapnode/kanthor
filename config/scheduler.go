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
	Metrics    metric.Config              `json:"metrics" yaml:"metrics" mapstructure:"metrics"`

	Request SchedulerRequest `json:"request" yaml:"request" mapstructure:"request"`
}

func (conf *Scheduler) Validate() error {
	if err := conf.Publisher.Validate(); err != nil {
		return fmt.Errorf("config.scheduler.publisher: %v", err)
	}
	if err := conf.Subscriber.Validate(); err != nil {
		return fmt.Errorf("config.scheduler.subscriber: %v", err)
	}
	if err := conf.Cache.Validate(); err != nil {
		return fmt.Errorf("config.scheduler.cache: %v", err)
	}
	if err := conf.Metrics.Validate(); err != nil {
		return fmt.Errorf("config.scheduler.metrics: %v", err)
	}
	if err := conf.Request.Validate(); err != nil {
		return fmt.Errorf("config.scheduler.request: %v", err)
	}

	return nil
}

type SchedulerRequest struct {
	Schedule SchedulerRequestSchedule `json:"schedule" yaml:"schedule" mapstructure:"schedule"`
}

func (conf *SchedulerRequest) Validate() error {
	if err := conf.Schedule.Validate(); err != nil {
		return fmt.Errorf("config.scheduler.request.schedule: %v", err)
	}
	return nil
}

type SchedulerRequestSchedule struct {
	ChunkSize    int   `json:"chunk_size" yaml:"chunk_size" mapstructure:"chunk_size"`
	ChunkTimeout int64 `json:"chunk_timeout" yaml:"chunk_timeout" mapstructure:"chunk_timeout"`
}

func (conf *SchedulerRequestSchedule) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThan("config.scheduler.request.schedule.chunk_size", conf.ChunkSize, 0),
		validator.NumberGreaterThanOrEqual("config.scheduler.request.schedule.chunk_timeout", conf.ChunkTimeout, 1000),
	)
}
