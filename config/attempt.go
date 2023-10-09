package config

import (
	"fmt"

	"github.com/scrapnode/kanthor/infrastructure/dlocker"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/utils"
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
	Schedule AttemptTriggerSchedule `json:"schedule" yaml:"schedule" mapstructure:"schedule"`
	Create   AttemptTriggerCreate   `json:"create" yaml:"create" mapstructure:"create"`
}

func (conf *AttemptTrigger) Validate() error {
	if err := conf.Cron.Validate(); err != nil {
		return fmt.Errorf("config.attempt.trigger.cron: %v", err)
	}
	if err := conf.Consumer.Validate(); err != nil {
		return fmt.Errorf("config.attempt.trigger.consumer: %v", err)
	}
	if err := conf.Schedule.Validate(); err != nil {
		return fmt.Errorf("config.attempt.trigger.schedule: %v", err)
	}
	if err := conf.Create.Validate(); err != nil {
		return fmt.Errorf("config.attempt.trigger.create: %v", err)
	}
	return nil
}

type AttemptTriggerCron struct {
	LockDuration int `json:"lock_duration" yaml:"lock_duration" mapstructure:"lock_duration"`
	ScanSize     int `json:"scan_size" yaml:"scan_size" mapstructure:"scan_size"`
	PublishSize  int `json:"publish_size" yaml:"publish_size" mapstructure:"publish_size"`
}

func (conf *AttemptTriggerCron) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThanOrEqual("config.attempt.trigger.cron.lock_duration", conf.LockDuration, 300),
		validator.NumberGreaterThan("config.attempt.trigger.cron.scan_size", conf.ScanSize, 0),
		validator.NumberLessThan("config.attempt.trigger.cron.scan_size", conf.ScanSize, 100000),
		validator.NumberGreaterThan("config.attempt.trigger.cron.publish_size", conf.PublishSize, 0),
		validator.NumberLessThan("config.attempt.trigger.cron.publish_size", conf.PublishSize, 1000),
	)
}

type AttemptTriggerConsumer struct {
	Subscriber   streaming.SubscriberConfig `json:"subscriber" yaml:"subscriber" mapstructure:"subscriber"`
	LockDuration int                        `json:"lock_duration" yaml:"lock_duration" mapstructure:"lock_duration"`
	ScanFrom     int                        `json:"scan_from" yaml:"scan_from" mapstructure:"scan_from"`
	ScanTo       int                        `json:"scan_to" yaml:"scan_to" mapstructure:"scan_to"`
}

func (conf *AttemptTriggerConsumer) Validate() error {
	if err := conf.Subscriber.Validate(); err != nil {
		return fmt.Errorf("config.attempt.trigger.consumer.subscriber: %v", err)
	}

	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThanOrEqual("config.attempt.trigger.consumer.lock_duration", conf.LockDuration, 300),
		validator.NumberGreaterThan("config.attempt.trigger.consumer.scan_to", utils.AbsInt(conf.ScanTo), conf.LockDuration),
		validator.NumberLessThan("config.attempt.trigger.consumer.scan_to", conf.ScanTo, 0),
		validator.NumberLessThan("config.attempt.trigger.consumer.scan_from", conf.ScanFrom, conf.ScanTo),
	)
}

type AttemptTriggerSchedule struct {
	Concurrency int `json:"concurrency" yaml:"concurrency" mapstructure:"concurrency"`
}

func (conf *AttemptTriggerSchedule) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThan("config.attempt.trigger.schedule.concurrency", conf.Concurrency, 0),
	)
}

type AttemptTriggerCreate struct {
	Concurrency   int `json:"concurrency" yaml:"concurrency" mapstructure:"concurrency"`
	ScheduleDelay int `json:"schedule_delay" yaml:"schedule_delay" mapstructure:"schedule_delay"`
}

func (conf *AttemptTriggerCreate) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThan("config.attempt.trigger.create.concurrency", conf.Concurrency, 0),
		validator.NumberGreaterThanOrEqual("config.attempt.trigger.create.schedule_delay", conf.ScheduleDelay, 900),
	)
}
