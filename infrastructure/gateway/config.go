package gateway

import "github.com/go-playground/validator/v10"

type Config struct {
	Protocol string     `json:"protocol" yaml:"protocol" mapstructure:"protocol" validate:"required,oneof=grpc"`
	GRPC     GRPCConfig `json:"grpc" yaml:"grpc" mapstructure:"grpc" validate:"-"`
}

func (conf *Config) Validate() error {
	if err := validator.New().Struct(conf); err != nil {
		return err
	}

	if conf.Protocol == "grpc" {
		if err := conf.GRPC.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type GRPCConfig struct {
	Addr string `json:"addr" yaml:"addr" mapstructure:"addr" validate:"required"`
}

func (conf *GRPCConfig) Validate() error {
	if err := validator.New().Struct(conf); err != nil {
		return err
	}
	return nil
}
