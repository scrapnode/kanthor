package exporter

import "github.com/scrapnode/kanthor/infrastructure/validator"

type Config struct {
	Addr string `json:"addr" yaml:"addr" validate:"required"`
}

func (conf *Config) Validate() error {
	if err := validator.New().Struct(conf); err != nil {
		return err
	}
	return nil
}
