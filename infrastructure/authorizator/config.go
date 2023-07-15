package authorizator

import (
	"errors"
	"github.com/go-playground/validator/v10"
)

type Config struct {
	Engine string        `json:"engine" yaml:"engine" mapstructure:"engine" validate:"required,oneof=casbin"`
	Casbin *CasbinConfig `json:"casbin" yaml:"casbin" mapstructure:"casbin" validate:"-"`
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

type CasbinConfig struct {
	ModelUri  string              `json:"model_uri" yaml:"model_uri" mapstructure:"model_uri" validate:"required,uri"`
	PolicyUri string              `json:"policy_uri" yaml:"policy_uri" mapstructure:"policy_uri" validate:"required,uri"`
	Watcher   CasbinWatcherConfig `json:"watcher" yaml:"watcher" mapstructure:"watcher" validate:"-"`
}

func (conf *CasbinConfig) Validate() error {
	if err := validator.New().Struct(conf); err != nil {
		return err
	}

	if err := conf.Watcher.Validate(); err != nil {
		return err
	}

	return nil
}

type CasbinWatcherConfig struct {
	Uri string `json:"uri" yaml:"uri" mapstructure:"uri" validate:"required,uri"`
}

func (conf *CasbinWatcherConfig) Validate() error {
	if err := validator.New().Struct(conf); err != nil {
		return err
	}

	return nil
}
