package config

import (
	"fmt"

	"github.com/scrapnode/kanthor/infrastructure/dlocker"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type Attempt struct {
	Publisher streaming.PublisherConfig `json:"publisher" yaml:"publisher" mapstructure:"publisher"`
	DLocker   dlocker.Config            `json:"dlocker" yaml:"dlocker" mapstructure:"dlocker"`
	Trigger   AttemptTrigger            `json:"trigger" yaml:"trigger" mapstructure:"trigger"`
	Metrics   metric.Config             `json:"metrics" yaml:"metrics" mapstructure:"metrics"`
}

func (conf *Attempt) Validate() error {
	if err := conf.Publisher.Validate(); err != nil {
		return fmt.Errorf("config.attempt.publisher: %v", err)
	}
	if err := conf.DLocker.Validate(); err != nil {
		return fmt.Errorf("config.attempt.dlocker: %v", err)
	}
	if err := conf.Trigger.Validate(); err != nil {
		return fmt.Errorf("config.attempt.trigger: %v", err)
	}
	if err := conf.Metrics.Validate(); err != nil {
		return fmt.Errorf("config.attempt.metrics: %v", err)
	}

	return nil
}

type AttemptTrigger struct {
	Cron     AttemptTriggerCron     `json:"cron" yaml:"cron" mapstructure:"cron"`
	Consumer AttemptTriggerConsumer `json:"consumer" yaml:"consumer" mapstructure:"consumer"`
}

func (conf *AttemptTrigger) Validate() error {
	if err := conf.Cron.Validate(); err != nil {
		return fmt.Errorf("config.attempt.trigger.cron: %v", err)
	}
	if err := conf.Consumer.Validate(); err != nil {
		return fmt.Errorf("config.attempt.trigger.consumer: %v", err)
	}
	return nil
}

type AttemptTriggerCron struct {
	LockDuration int64 `json:"lock_duration" yaml:"lock_duration" mapstructure:"lock_duration"`
	ChunkTimeout int64 `json:"chunk_timeout" yaml:"chunk_timeout" mapstructure:"chunk_timeout"`
	ChunkSize    int   `json:"chunk_size" yaml:"chunk_size" mapstructure:"chunk_size"`

	ScanFrom int64 `json:"scan_from" yaml:"scan_from" mapstructure:"scan_from"`
	ScanTo   int64 `json:"scan_to" yaml:"scan_to" mapstructure:"scan_to"`
}

func (conf *AttemptTriggerCron) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThanOrEqual("config.attempt.trigger.cron.lock_duration", conf.LockDuration, 15000),
		validator.NumberGreaterThanOrEqual("config.attempt.trigger.cron.chunk_timeout", conf.ChunkTimeout, 3000),
		validator.NumberGreaterThan("config.attempt.trigger.cron.chunk_size", conf.ChunkSize, 0),
		validator.NumberGreaterThan("scan_from", conf.ScanFrom, conf.ScanTo),
		validator.NumberLessThan("scan_to", conf.ScanTo, conf.ScanFrom),
	)
}

type AttemptTriggerConsumer struct {
	Delay        int64 `json:"delay" yaml:"delay" mapstructure:"delay"`
	ChunkTimeout int64 `json:"chunk_timeout" yaml:"chunk_timeout" mapstructure:"chunk_timeout"`
	ChunkSize    int   `json:"chunk_size" yaml:"chunk_size" mapstructure:"chunk_size"`

	Subscriber streaming.SubscriberConfig `json:"subscriber" yaml:"subscriber" mapstructure:"subscriber"`
}

func (conf *AttemptTriggerConsumer) Validate() error {
	err := validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThanOrEqual("config.attempt.trigger.cron.chunk_timeout", conf.ChunkTimeout, 3000),
		validator.NumberGreaterThanOrEqual("config.attempt.trigger.cron.delay", conf.Delay, 3000),
		validator.NumberGreaterThan("config.attempt.trigger.cron.chunk_size", conf.ChunkSize, 0),
	)
	if err != nil {
		return err
	}

	if err := conf.Subscriber.Validate(); err != nil {
		return fmt.Errorf("config.attempt.trigger.consumer.subscriber: %v", err)
	}

	return nil
}
