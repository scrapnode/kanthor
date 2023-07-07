package cache

import "github.com/go-playground/validator/v10"

type Config struct {
	Uri        string `json:"uri" yaml:"uri" mapstructure:"uri" validate:"required,uri"`
	TimeToLive int    `json:"time_to_live" yaml:"timeToLive" mapstructure:"time_to_live" validate:"required,number,gte=0"`
}

func (conf Config) Validate() error {
	return validator.New().Struct(conf)
}
