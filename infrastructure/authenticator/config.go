package authenticator

import "github.com/scrapnode/kanthor/pkg/validator"

type Config struct {
	Engine  string   `json:"engine" yaml:"engine" mapstructure:"engine"`
	Ask     *Ask     `json:"ask" yaml:"ask" mapstructure:"ask"`
	Forward *Forward `json:"forward" yaml:"forward" mapstructure:"forward"`
}

func (conf *Config) Validate(prefix string) error {
	err := validator.Validate(
		validator.DefaultConfig,
		validator.StringOneOf(prefix+".AUTHENTICATOR.SCHEME", conf.Engine, []string{EngineAsk, EngineForward}),
	)
	if err != nil {
		return err
	}

	if conf.Engine == EngineAsk {
		return conf.Ask.Validate(prefix)
	}

	if conf.Engine == EngineForward {
		return conf.Forward.Validate(prefix)
	}

	return nil
}

type Ask struct {
	AccessKey string `json:"access_key" yaml:"access_key" mapstructure:"access_key"`
	SecretKey string `json:"secret_key" yaml:"secret_key" mapstructure:"secret_key"`
}

func (conf *Ask) Validate(prefix string) error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringRequired(prefix+"AUTHENTICATOR.ASK.ACCESS_KEY", conf.AccessKey),
		validator.StringRequired(prefix+"AUTHENTICATOR.ASK.SECRET_KEY", conf.SecretKey),
	)
}

type Forward struct {
	Uri     string   `json:"uri" yaml:"uri" mapstructure:"uri"`
	Headers []string `json:"headers" yaml:"headers" mapstructure:"headers"`
}

func (conf *Forward) Validate(prefix string) error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringUri(prefix+"AUTHENTICATOR.FORWARD.URI", conf.Uri),
	)
}
