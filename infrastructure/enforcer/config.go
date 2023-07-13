package enforcer

import (
	"errors"
	"github.com/go-playground/validator/v10"
)

type Config struct {
	Engine string        `json:"engine" yaml:"engine" mapstructure:"engine" validate:"required,oneof=casbin"`
	Casbin *CasbinConfig `json:"casbin" yaml:"casbin" mapstructure:"casbin" validate:"-"`
}

type CasbinConfig struct {
	ModelSource  string `json:"model_source" yaml:"model_source" mapstructure:"model_source" validate:"required,uri"`
	PolicySource string `json:"policy_source" yaml:"policy_source" mapstructure:"policy_source" validate:"required,uri"`
}

func (conf *CasbinConfig) Validate() error {
	return validator.New().Struct(conf)
}

func (conf *Config) Validate() error {
	if err := validator.New().Struct(conf); err != nil {
		return err
	}

	if conf.Engine == "casbin" {
		if conf.Casbin == nil {
			return errors.New("enforcer.config.casbin: null value")
		}
		if err := conf.Casbin.Validate(); err != nil {
			return err
		}
	}

	return nil
}
