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
		validator.StringRequiredOneOf("CONFIG.INFRA.AUTHORIZATOR.ENGINE", conf.Engine, []string{EngineCasbin}),
	)
	if err != nil {
		return err
	}

	if conf.Engine == EngineCasbin {
		if conf.Casbin == nil {
			return errors.New("CONFIG.INFRA.AUTHORIZATOR.CASBIN: nil value")
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
		validator.StringUri("CONFIG.INFRA.AUTHORIZATOR.CASBIN.MODEL_URI", conf.ModelUri),
		validator.StringUri("CONFIG.INFRA.AUTHORIZATOR.CASBIN.POLICY_URI", conf.PolicyUri),
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
