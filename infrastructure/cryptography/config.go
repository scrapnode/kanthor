package cryptography

import (
	"github.com/go-playground/validator/v10"
)

type Config struct {
	KDF       KDFConfig       `json:"kdf" yaml:"kdf" mapstructure:"kdf"`
	Symmetric SymmetricConfig `json:"symmetric" yaml:"symmetric" mapstructure:"symmetric" validate:"required"`
}

func (conf *Config) Validate() error {
	if err := validator.New().Struct(conf); err != nil {
		return err
	}

	if err := conf.KDF.Validate(); err != nil {
		return err
	}

	if err := conf.Symmetric.Validate(); err != nil {
		return err
	}

	return nil
}

type KDFConfig struct {
}

func (conf *KDFConfig) Validate() error {
	if err := validator.New().Struct(conf); err != nil {
		return err
	}
	return nil
}

type SymmetricConfig struct {
	Key string `json:"key" yaml:"key" mapstructure:"key" validate:"required,len=32"`
}

func (conf *SymmetricConfig) Validate() error {
	if err := validator.New().Struct(conf); err != nil {
		return err
	}
	return nil
}
