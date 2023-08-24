package sender

import "github.com/go-playground/validator/v10"

type Config struct {
	EnableTrace bool  `json:"enable_trace" yaml:"enable_trace" mapstructure:"enable_trace"`
	Timeout     int   `json:"timeout" yaml:"timeout" mapstructure:"timeout" validate:"required,number,gte=0"`
	Retry       Retry `json:"retry" yaml:"retry" mapstructure:"retry" validate:"required"`
}

func (conf *Config) Validate() error {
	if err := conf.Retry.Validate(); err != nil {
		return err
	}

	return validator.New().Struct(conf)
}

type Retry struct {
	Count       int `json:"count" yaml:"count" mapstructure:"count" validate:"required,number,gte=0"`
	WaitTime    int `json:"wait_time" yaml:"wait_time" mapstructure:"wait_time" validate:"required,number,gte=0"`
	WaitTimeMax int `json:"wait_time_max" yaml:"wait_time_max" mapstructure:"wait_time_max" validate:"required,number,gte=0,gtfield=Count"`
}

func (conf *Retry) Validate() error {
	return validator.New().Struct(conf)
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
