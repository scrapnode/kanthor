package idempotency

import "github.com/go-playground/validator/v10"

type Config struct {
	Namespace  string `json:"namespace" yaml:"namespace" validate:"required"`
	Uri        string `json:"uri" yaml:"uri" mapstructure:"uri" validate:"required,uri"`
	TimeToLive int    `json:"time_to_live" yaml:"timeToLive" mapstructure:"time_to_live" validate:"required,number,gte=1000"`
}

func (conf *Config) Validate() error {
	if err := validator.New().Struct(conf); err != nil {
		return err
	}

	return nil
}
