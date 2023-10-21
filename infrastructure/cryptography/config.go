package cryptography

import (
	"github.com/scrapnode/kanthor/pkg/validator"
)

type Config struct {
	KDF       KDFConfig       `json:"kdf" yaml:"kdf" mapstructure:"kdf"`
	Symmetric SymmetricConfig `json:"symmetric" yaml:"symmetric" mapstructure:"symmetric"`
}

func (conf *Config) Validate() error {
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
	return nil
}

type SymmetricConfig struct {
	Key string `json:"key" yaml:"key" mapstructure:"key"`
}

func (conf *SymmetricConfig) Validate() error {
	return validator.Validate(validator.DefaultConfig, validator.StringLen("cryptography.symmetric.key", conf.Key, 32, 32))
}
