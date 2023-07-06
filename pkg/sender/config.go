package sender

import "github.com/go-playground/validator/v10"

type Config struct {
	EnableTrace bool  `json:"enable_trace" mapstructure:"enable_trace"`
	Timeout     int   `json:"timeout" mapstructure:"timeout" validate:"required,number,gte=0"`
	Retry       Retry `json:"retry" mapstructure:"retry" validate:"required"`
}

func (conf Config) Validate() error {
	if err := conf.Retry.Validate(); err != nil {
		return err
	}

	return validator.New().Struct(conf)
}

type Retry struct {
	Count       int `json:"count" mapstructure:"count" validate:"required,number,gte=0"`
	WaitTime    int `json:"wait_time" mapstructure:"wait_time" validate:"required,number,gte=0"`
	WaitTimeMax int `json:"wait_time_max" mapstructure:"wait_time_max" validate:"required,number,gte=0"`
}

func (conf Retry) Validate() error {
	return validator.New().Struct(conf)
}
