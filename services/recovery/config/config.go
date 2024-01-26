package config

import (
	"fmt"
	"log"

	"github.com/scrapnode/kanthor/configuration"
	"github.com/scrapnode/kanthor/pkg/validator"
)

func New(provider configuration.Provider) (*Config, error) {
	var conf Wrapper
	if err := provider.Unmarshal(&conf); err != nil {
		return nil, err
	}
	return &conf.Recovery, conf.Validate()
}

type Wrapper struct {
	Recovery Config `json:"recovery" yaml:"recovery" mapstructure:"recovery"`
}

func (conf *Wrapper) Validate() error {
	if err := conf.Recovery.Validate(); err != nil {
		return err
	}
	return nil
}

type Config struct {
	Cronjob  RecoveryCronjob  `json:"cronjob" yaml:"cronjob" mapstructure:"cronjob"`
	Consumer RecoveryConsumer `json:"consumer" yaml:"consumer" mapstructure:"consumer"`
}

func (conf *Config) Validate() error {
	if err := conf.Cronjob.Validate(); err != nil {
		return err
	}
	if err := conf.Consumer.Validate(); err != nil {
		return err
	}

	return nil
}

type RecoveryCronjob struct {
	Scheduler string                  `json:"scheduler" yaml:"scheduler" mapstructure:"scheduler"`
	Timeout   int64                   `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
	BatchSize int                     `json:"batch_size" yaml:"batch_size" mapstructure:"batch_size"`
	Buckets   []RecoveryCronjobBucket `json:"buckets" yaml:"buckets" mapstructure:"buckets"`
}

func (conf *RecoveryCronjob) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringRequired("RECOVERY.CONFIG.CRONJOB.SCHEDULER", conf.Scheduler),
		validator.NumberGreaterThanOrEqual("RECOVERY.CONFIG.CRONJOB.TIMEOUT", conf.Timeout, 1000),
		validator.NumberGreaterThan("RECOVERY.CONFIG.CRONJOB.BATCH_SIZE", conf.BatchSize, 0),
		validator.SliceRequired("RECOVERY.CONFIG.CRONJOB.BUCKETS", conf.Buckets),
		validator.Slice(conf.Buckets, func(i int, item *RecoveryCronjobBucket) error {
			return item.Validate(fmt.Sprintf("RECOVERY.CONFIG.CRONJOB.BUCKETS[%d]", i))
		}),
	)
}

type RecoveryCronjobBucket struct {
	Offset   int64 `json:"offset" yaml:"offset" mapstructure:"offset"`
	Duration int64 `json:"duration" yaml:"duration" mapstructure:"duration"`
}

func (conf *RecoveryCronjobBucket) Validate(prefix string) error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThanOrEqual(prefix+".OFFSET", conf.Offset, 0),
		validator.NumberGreaterThan(prefix+".DURATION", conf.Duration, 0),
	)
}

type RecoveryConsumer struct {
	Timeout   int64 `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
	BatchSize int   `json:"batch_size" yaml:"batch_size" mapstructure:"batch_size"`
}

func (conf *RecoveryConsumer) Validate() error {
	log.Println(validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThanOrEqual("RECOVERY.CONFIG.CRONJOB.TIMEOUT", conf.Timeout, 1000),
		validator.NumberGreaterThan("RECOVERY.CONFIG.CRONJOB.BATCH_SIZE", conf.BatchSize, 0),
	))
	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThanOrEqual("RECOVERY.CONFIG.CRONJOB.TIMEOUT", conf.Timeout, 1000),
		validator.NumberGreaterThan("RECOVERY.CONFIG.CRONJOB.BATCH_SIZE", conf.BatchSize, 0),
	)
}
