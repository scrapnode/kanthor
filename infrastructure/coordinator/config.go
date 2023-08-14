package coordinator

import (
	"errors"
	"github.com/go-playground/validator/v10"
)

var (
	EngineNats   = "nats"
	HeaderNodeId = "kanthor-coordinator-node-id"
)

type Config struct {
	Engine string      `json:"engine" yaml:"engine" mapstructure:"engine" validate:"required,oneof=nats"`
	Nats   *NatsConfig `json:"nats" yaml:"nats" mapstructure:"nats" validate:"-"`
}

func (conf *Config) Validate() error {
	if err := validator.New().Struct(conf); err != nil {
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
	Uri     string `json:"uri" yaml:"uri" mapstructure:"uri" validate:"required,uri"`
	Subject string `json:"subject" yaml:"subject" mapstructure:"subject" validate:"required"`
}

func (conf *NatsConfig) Validate() error {
	if err := validator.New().Struct(conf); err != nil {
		return err
	}

	return nil
}
