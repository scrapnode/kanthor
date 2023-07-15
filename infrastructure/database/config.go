package database

import "github.com/go-playground/validator/v10"

type Config struct {
	Uri string `json:"uri" yaml:"uri" mapstructure:"uri" validate:"required,uri"`
}

func (conf *Config) Validate() error {
	if err := validator.New().Struct(conf); err != nil {
		return err
	}
	return nil
}
