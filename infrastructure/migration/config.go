package migration

import "github.com/go-playground/validator/v10"

type Config struct {
	Source string `json:"source" yaml:"source"  mapstructure:"source" validate:"required,uri"`
}

func (conf *Config) Validate() error {
	if err := validator.New().Struct(conf); err != nil {
		return err
	}
	return nil
}
