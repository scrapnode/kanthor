package config

import (
	"fmt"

	"github.com/scrapnode/kanthor/configuration"
	"github.com/scrapnode/kanthor/pkg/validator"
)

func New(provider configuration.Provider) (*Config, error) {
	var conf Wrapper
	return &conf.Attempt, provider.Unmarshal(&conf)
}

type Wrapper struct {
	Attempt Config `json:"attempt" yaml:"attempt" mapstructure:"attempt"`
}

func (conf *Wrapper) Validate() error {
	if err := conf.Attempt.Validate(); err != nil {
		return err
	}
	return nil
}

type Config struct {
	Cronjob  AttemptCronjob  `json:"cronjob" yaml:"cronjob" mapstructure:"cronjob"`
	Consumer AttemptConsumer `json:"consumer" yaml:"consumer" mapstructure:"consumer"`
	Trigger  AttemptTrigger  `json:"trigger" yaml:"trigger" mapstructure:"trigger"`
	Selector AttemptSelector `json:"selector" yaml:"selector" mapstructure:"selector"`
	Endeavor AttemptEndeavor `json:"endeavor" yaml:"endeavor" mapstructure:"endeavor"`
}

func (conf *Config) Validate() error {
	if err := conf.Cronjob.Validate(); err != nil {
		return err
	}

	return nil
}

type AttemptCronjob struct {
	Scheduler string          `json:"scheduler" yaml:"scheduler" mapstructure:"scheduler"`
	Timeout   int64           `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
	BatchSize int             `json:"batch_size" yaml:"batch_size" mapstructure:"batch_size"`
	Buckets   []AttemptBucket `json:"buckets" yaml:"buckets" mapstructure:"buckets"`
}

func (conf *AttemptCronjob) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringRequired("CONFIG.ATTEMPT.CRONJOB.SCHEDULER", conf.Scheduler),
		validator.NumberGreaterThanOrEqual("CONFIG.ATTEMPT.CRONJOB.TIMEOUT", conf.Timeout, 1000),
		validator.NumberGreaterThan("CONFIG.ATTEMPT.CRONJOB.BATCH_SIZE", conf.BatchSize, 0),
		validator.SliceRequired("CONFIG.ATTEMPT.CRONJOB.BUCKETS", conf.Buckets),
		validator.Slice(conf.Buckets, func(i int, item *AttemptBucket) error {
			return item.Validate(fmt.Sprintf("CONFIG.ATTEMPT.CRONJOB.BUCKETS[%d]", i))
		}),
	)
}

type AttemptBucket struct {
	Offset   int64 `json:"offset" yaml:"offset" mapstructure:"offset"`
	Duration int64 `json:"duration" yaml:"duration" mapstructure:"duration"`
}

func (conf *AttemptBucket) Validate(prefix string) error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThanOrEqual(prefix+".OFFSET", conf.Offset, 0),
		validator.NumberGreaterThan(prefix+".DURATION", conf.Duration, 0),
	)
}

type AttemptConsumer struct {
	Timeout       int64 `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
	BatchSize     int   `json:"batch_size" yaml:"batch_size" mapstructure:"batch_size"`
	ScheduleDelay int   `json:"schedule_delay" yaml:"schedule_delay" mapstructure:"schedule_delay"`
}

func (conf *AttemptConsumer) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThanOrEqual("CONFIG.ATTEMPT.CRONJOB.TIMEOUT", conf.Timeout, 1000),
		validator.NumberGreaterThan("CONFIG.ATTEMPT.CRONJOB.BATCH_SIZE", conf.BatchSize, 0),
		validator.NumberGreaterThan("CONFIG.ATTEMPT.CRONJOB.SCHEDULE_DELAY", conf.ScheduleDelay, 0),
	)
}

type AttemptTrigger struct {
	Scheduler string          `json:"scheduler" yaml:"scheduler" mapstructure:"scheduler"`
	Timeout   int64           `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
	Buckets   []AttemptBucket `json:"buckets" yaml:"buckets" mapstructure:"buckets"`
}

func (conf *AttemptTrigger) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringRequired("CONFIG.ATTEMPT.TRIGGER.SCHEDULER", conf.Scheduler),
		validator.NumberGreaterThanOrEqual("CONFIG.ATTEMPT.TRIGGER.TIMEOUT", conf.Timeout, 1000),
		validator.SliceRequired("CONFIG.ATTEMPT.TRIGGER.BUCKETS", conf.Buckets),
		validator.Slice(conf.Buckets, func(i int, item *AttemptBucket) error {
			return item.Validate(fmt.Sprintf("CONFIG.ATTEMPT.TRIGGER.BUCKETS[%d]", i))
		}),
	)
}

type AttemptSelector struct {
	Timeout   int64 `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
	BatchSize int   `json:"batch_size" yaml:"batch_size" mapstructure:"batch_size"`
	Counter   int   `json:"counter" yaml:"counter" mapstructure:"counter"`
}

func (conf *AttemptSelector) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThanOrEqual("CONFIG.ATTEMPT.SELECTOR.TIMEOUT", conf.Timeout, 1000),
		validator.NumberGreaterThan("CONFIG.ATTEMPT.SELECTOR.BATCH_SIZE", conf.BatchSize, 0),
		validator.NumberGreaterThan("CONFIG.ATTEMPT.SELECTOR.COUNTER", conf.Counter, 0),
	)
}

type AttemptEndeavor struct {
	Concurrency int `json:"concurrency" yaml:"concurrency" mapstructure:"concurrency"`
	RetryDelay  int `json:"retry_delay" yaml:"retry_delay" mapstructure:"retry_delay"`
}

func (conf *AttemptEndeavor) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThanOrEqual("CONFIG.ATTEMPT.ENDEAVOR.CONCURRENCY", conf.Concurrency, 0),
		validator.NumberGreaterThanOrEqual("CONFIG.ATTEMPT.ENDEAVOR.RETRY_DELAY", conf.RetryDelay, 1000),
	)
}
