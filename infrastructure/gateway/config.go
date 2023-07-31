package gateway

import (
	"errors"
	"github.com/go-playground/validator/v10"
)

const EngineGrpc = "grpc"

type Config struct {
	Engine string      `json:"engine" yaml:"engine" mapstructure:"engine" validate:"required,oneof=grpc"`
	Grpc   *GrpcConfig `json:"grpc" yaml:"grpc" mapstructure:"grpc" validate:"-"`
}

func (conf *Config) Validate() error {
	if err := validator.New().Struct(conf); err != nil {
		return err
	}

	if conf.Engine == EngineGrpc {
		if conf.Grpc == nil {
			return errors.New("gateway.config.grpc: null value")
		}
		if err := conf.Grpc.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type GrpcConfig struct {
	Addr string `json:"addr" yaml:"addr" validate:"required"`
}

func (conf *GrpcConfig) Validate() error {
	if err := validator.New().Struct(conf); err != nil {
		return err
	}
	return nil
}
