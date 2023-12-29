package authenticator

import (
	"github.com/scrapnode/kanthor/pkg/validator"
)

type Config struct {
	Engine   string    `json:"engine" yaml:"engine" mapstructure:"engine"`
	Ask      *Ask      `json:"ask" yaml:"ask" mapstructure:"ask"`
	External *External `json:"external" yaml:"external" mapstructure:"external"`
}

func (conf *Config) Validate() error {
	err := validator.Validate(
		validator.DefaultConfig,
		validator.StringOneOf("AUTHENTICATOR.SCHEME", conf.Engine, []string{EngineAsk, EngineExternal}),
	)
	if err != nil {
		return err
	}

	if conf.Engine == EngineAsk {
		return conf.Ask.Validate()
	}

	if conf.Engine == EngineExternal {
		return conf.External.Validate()
	}

	return nil
}

type Ask struct {
	AccessKey string `json:"access_key" yaml:"access_key" mapstructure:"access_key"`
	SecretKey string `json:"secret_key" yaml:"secret_key" mapstructure:"secret_key"`
}

func (conf *Ask) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringRequired("AUTHENTICATOR.ASK.ACCESS_KEY", conf.AccessKey),
		validator.StringRequired("AUTHENTICATOR.ASK.SECRET_KEY", conf.SecretKey),
	)
}

type External struct {
	Uri     string          `json:"uri" yaml:"uri" mapstructure:"uri"`
	Headers []string        `json:"headers" yaml:"headers" mapstructure:"headers"`
	Mapper  *ExternalMapper `json:"mapper" yaml:"mapper" mapstructure:"mapper"`
}

func (conf *External) Validate() error {
	err := validator.Validate(
		validator.DefaultConfig,
		validator.StringUri("AUTHENTICATOR.EXTERNAL.URI", conf.Uri),
	)
	if err != nil {
		return err
	}

	if conf.Mapper != nil {
		if err := conf.Mapper.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type ExternalMapper struct {
	Uri string `json:"uri" yaml:"uri" mapstructure:"uri"`
}

func (conf *ExternalMapper) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringUri("AUTHENTICATOR.EXTERNAL.MAPPER.URI", conf.Uri),
		validator.StringStartsWithOneOf("AUTHENTICATOR.EXTERNAL.MAPPER.URI", conf.Uri, []string{"base64"}),
	)
}
