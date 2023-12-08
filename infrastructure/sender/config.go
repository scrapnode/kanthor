package sender

import "github.com/scrapnode/kanthor/pkg/validator"

type Config struct {
	Trace   bool  `json:"trace" yaml:"trace" mapstructure:"trace"`
	Timeout int   `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
	Retry   Retry `json:"retry" yaml:"retry" mapstructure:"retry"`
}

func (conf *Config) Validate() error {
	err := validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThanOrEqual("CONFIG.INFRA.SENDER.TIMEOUT", conf.Timeout, 0),
	)
	if err != nil {
		return err
	}

	return conf.Retry.Validate()
}

type Retry struct {
	Count    int `json:"count" yaml:"count" mapstructure:"count"`
	WaitTime int `json:"wait_time" yaml:"wait_time" mapstructure:"wait_time"`
}

func (conf *Retry) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThanOrEqual("CONFIG.INFRA.SENDER.RETRY.COUNT", conf.Count, 0),
		validator.NumberGreaterThanOrEqual("CONFIG.INFRA.SENDER.RETRY.WAIT_TIME", conf.WaitTime, 100),
	)
}

var DefaultConfig = &Config{
	Trace:   false,
	Timeout: 5000,
	Retry: Retry{
		Count:    3,
		WaitTime: 100,
	},
}
