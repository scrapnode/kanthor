package authorizator

import (
	"errors"

	"github.com/scrapnode/kanthor/pkg/validator"
)

var (
	EngineCasbin = "casbin"
)

type Config struct {
	Engine string        `json:"engine" yaml:"engine" mapstructure:"engine"`
	Casbin *CasbinConfig `json:"casbin" yaml:"casbin" mapstructure:"casbin"`
}

func (conf *Config) Validate() error {
	err := validator.Validate(
		validator.DefaultConfig,
		validator.StringRequiredOneOf("authorizator.config.engine", conf.Engine, []string{EngineCasbin}),
	)
	if err != nil {
		return err
	}

	if conf.Engine == EngineCasbin {
		if conf.Casbin == nil {
			return errors.New("authorizator.config.casbin: null value")
		}
		if err := conf.Casbin.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type CasbinConfig struct {
	ModelUri  string              `json:"model_uri" yaml:"model_uri" mapstructure:"model_uri"`
	PolicyUri string              `json:"policy_uri" yaml:"policy_uri" mapstructure:"policy_uri"`
	Watcher   CasbinWatcherConfig `json:"watcher" yaml:"watcher" mapstructure:"watcher" validate:"-"`
}

func (conf *CasbinConfig) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringUri("authorizator.conf.casbin.model_uri", conf.ModelUri),
		validator.StringUri("authorizator.conf.casbin.policy_uri", conf.PolicyUri),
	)
}

type CasbinWatcherConfig struct {
	Uri string `json:"uri" yaml:"uri" mapstructure:"uri"`
}

func (conf *CasbinWatcherConfig) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringUri("authorizator.conf.casbin.watcher.uri", conf.Uri),
	)
}
