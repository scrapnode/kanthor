package metric

import "github.com/go-playground/validator/v10"

type Config struct {
	Enable    bool   `json:"enable" yaml:"enable" mapstructure:"enable"`
	Namespace string `json:"namespace" yaml:"namespace" mapstructure:"namespace"`
	Exporter  struct {
		Addr string `json:"addr" yaml:"addr" mapstructure:"addr" validate:"required"`
	} `json:"exporter" yaml:"exporter" mapstructure:"exporter" validate:"required"`
}

func (conf Config) Validate() error {
	if !conf.Enable {
		return nil
	}
	return validator.New().Struct(conf)
}
