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
}

func (conf *Config) Validate() error {
	if err := conf.Cronjob.Validate(); err != nil {
		return err
	}

	return nil
}

type AttemptCronjob struct {
	Scheduler string                 `json:"scheduler" yaml:"scheduler" mapstructure:"scheduler"`
	Timeout   int64                  `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
	BatchSize int                    `json:"batch_size" yaml:"batch_size" mapstructure:"batch_size"`
	Buckets   []AttemptCronjobBucket `json:"buckets" yaml:"buckets" mapstructure:"buckets"`
}

func (conf *AttemptCronjob) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringRequired("CONFIG.ATTEMPT.CRONJOB.SCHEDULER", conf.Scheduler),
		validator.NumberGreaterThanOrEqual("CONFIG.ATTEMPT.CRONJOB.TIMEOUT", conf.Timeout, 1000),
		validator.NumberGreaterThan("CONFIG.ATTEMPT.CRONJOB.BATCH_SIZE", conf.BatchSize, 0),
		validator.SliceRequired("CONFIG.ATTEMPT.CRONJOB.BUCKETS", conf.Buckets),
		validator.Slice(conf.Buckets, func(i int, item *AttemptCronjobBucket) error {
			return item.Validate(fmt.Sprintf("CONFIG.ATTEMPT.CRONJOB.BUCKETS[%d]", i))
		}),
	)
}

type AttemptCronjobBucket struct {
	Offset   int64 `json:"offset" yaml:"offset" mapstructure:"offset"`
	Duration int64 `json:"duration" yaml:"duration" mapstructure:"duration"`
}

func (conf *AttemptCronjobBucket) Validate(prefix string) error {
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
