package metric

import "github.com/go-playground/validator/v10"

type Config struct {
	Enable   bool `json:"enable" mapstructure:"enable"`
	Exporter struct {
		Addr string `json:"addr" mapstructure:"addr" validate:"required"`
	} `json:"exporter" mapstructure:"exporter" validate:"required"`
}

func (conf Config) Validate() error {
	if !conf.Enable {
		return nil
	}
	return validator.New().Struct(conf)
}
