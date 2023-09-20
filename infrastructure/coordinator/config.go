package coordinator

import (
	"errors"

	"github.com/scrapnode/kanthor/pkg/validator"
)

var (
	EngineNats   = "nats"
	HeaderNodeId = "kanthor-coordinator-node-id"
	HeaderCmd    = "kanthor-coordinator-cmd"
)

type Config struct {
	Engine string      `json:"engine" yaml:"engine" mapstructure:"engine"`
	Nats   *NatsConfig `json:"nats" yaml:"nats" mapstructure:"nats"`
}

func (conf *Config) Validate() error {
	err := validator.Validate(
		validator.DefaultConfig,
		validator.StringOneOf("coordinator.config.engine", conf.Engine, []string{EngineNats}),
	)
	if err != nil {
		return err
	}

	if conf.Engine == EngineNats {
		if conf.Nats == nil {
			return errors.New("coordinator.config.nats: null value")
		}
		if err := conf.Nats.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type NatsConfig struct {
	Uri     string `json:"uri" yaml:"uri" mapstructure:"uri"`
	Subject string `json:"subject" yaml:"subject" mapstructure:"subject"`
}

func (conf *NatsConfig) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringUri("coordinator.conf.nats.uri", conf.Uri),
		validator.StringRequired("coordinator.conf.nats.subject", conf.Subject),
	)
}
