package circuitbreaker

import "github.com/go-playground/validator/v10"

type Config struct {
	Timeout               int `json:"timeout" mapstructure:"timeout" validate:"required,number,gte=0"`
	SleepWindow           int `json:"sleep_window" mapstructure:"sleep_window" validate:"required,number,gte=0"`
	ErrorPercentThreshold int `json:"error_percent_threshold" mapstructure:"error_percent_threshold" validate:"required,number,gte=0"`
}

func (conf Config) Validate() error {
	return validator.New().Struct(conf)
}
