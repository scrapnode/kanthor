package cryptography

import "github.com/go-playground/validator/v10"

type SymmetricConfig struct {
	Key string `json:"key" yaml:"key" mapstructure:"key" validate:"required,len=32"`
}

func (conf *SymmetricConfig) Validate() error {
	if err := validator.New().Struct(conf); err != nil {
		return err
	}
	return nil
}
