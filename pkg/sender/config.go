package sender

import "github.com/scrapnode/kanthor/pkg/validator"

type Config struct {
	EnableTrace bool  `json:"enable_trace" yaml:"enable_trace" mapstructure:"enable_trace"`
	Timeout     int   `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
	Retry       Retry `json:"retry" yaml:"retry" mapstructure:"retry"`
}

func (conf *Config) Validate() error {
	err := validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThanOrEqual("sender.timeout", conf.Timeout, 0),
	)
	if err != nil {
		return err
	}

	return conf.Retry.Validate()
}

type Retry struct {
	Count       int `json:"count" yaml:"count" mapstructure:"count"`
	WaitTime    int `json:"wait_time" yaml:"wait_time" mapstructure:"wait_time"`
	WaitTimeMax int `json:"wait_time_max" yaml:"wait_time_max" mapstructure:"wait_time_max"`
}

func (conf *Retry) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThanOrEqual("sender.retry.count", conf.Count, 0),
		validator.NumberGreaterThanOrEqual("sender.retry.wait_time", conf.WaitTime, 0),
		validator.NumberGreaterThanOrEqual("sender.retry.wait_time_max", conf.WaitTimeMax, 0),
		validator.NumberGreaterThan("sender.retry.wait_time_max>sender.retry.wait_time", conf.WaitTimeMax, conf.WaitTime),
	)
}

var DefaultConfig = &Config{
	EnableTrace: false,
	Timeout:     3000,
	Retry: Retry{
		Count:       1,
		WaitTime:    200,
		WaitTimeMax: 5000,
	},
}
