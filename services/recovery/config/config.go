package config

import (
	"fmt"

	"github.com/scrapnode/kanthor/configuration"
	"github.com/scrapnode/kanthor/pkg/validator"
)

// @TODO:mapstructure with env
func New(provider configuration.Provider) (*Config, error) {
	var conf Wrapper
	return &conf.Recovery, provider.Unmarshal(&conf)
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
	Scanner RecoveryScanner `json:"scanner" yaml:"scanner" mapstructure:"scanner"`
}

func (conf *Config) Validate() error {
	if err := conf.Scanner.Validate(); err != nil {
		return err
	}

	return nil
}

type RecoveryScanner struct {
	Scheduler string                  `json:"scheduler" yaml:"scheduler" mapstructure:"scheduler"`
	Timeout   int64                   `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
	BatchSize int                     `json:"batch_size" yaml:"batch_size" mapstructure:"batch_size"`
	Buckets   []RecoveryScannerBucket `json:"buckets" yaml:"buckets" mapstructure:"buckets"`
}

func (conf *RecoveryScanner) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringRequired("CONFIG.RECOVERY.SCANNER.SCHEDULER", conf.Scheduler),
		validator.NumberGreaterThanOrEqual("CONFIG.RECOVERY.SCANNER.TIMEOUT", conf.Timeout, 1000),
		validator.NumberGreaterThan("CONFIG.RECOVERY.SCANNER.BATCH_SIZE", conf.BatchSize, 0),
		validator.SliceRequired("CONFIG.RECOVERY.SCANNER.BUCKETS", conf.Buckets),
		validator.Slice(conf.Buckets, func(i int, item *RecoveryScannerBucket) error {
			return item.Validate(fmt.Sprintf("CONFIG.RECOVERY.SCANNER.BUCKETS[%d]", i))
		}),
	)
}

type RecoveryScannerBucket struct {
	Offset   int64 `json:"offset" yaml:"offset" mapstructure:"offset"`
	Duration int64 `json:"duration" yaml:"duration" mapstructure:"duration"`
}

func (conf *RecoveryScannerBucket) Validate(prefix string) error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThanOrEqual(prefix+".OFFSET", conf.Offset, 0),
		validator.NumberGreaterThan(prefix+".DURATION", conf.Duration, 0),
	)
}
