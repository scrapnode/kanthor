package config

import (
	"github.com/scrapnode/kanthor/configuration"
	"github.com/scrapnode/kanthor/domain/constants"
	"github.com/scrapnode/kanthor/pkg/validator"
)

// @TODO: mapstructure with env
func New(provider configuration.Provider) (*Config, error) {
	var conf Wrapper
	return &conf.Storage, provider.Unmarshal(&conf)
}

type Wrapper struct {
	Storage Config `json:"storage" yaml:"storage" mapstructure:"storage"`
}

func (conf *Wrapper) Validate() error {
	if err := conf.Storage.Validate(); err != nil {
		return err
	}
	return nil
}

type Config struct {
	Topic     string           `json:"topic" yaml:"topic" mapstructure:"STORAGE_TOPIC"`
	Warehouse StorageWarehouse `json:"warehouse" yaml:"warehouse" mapstructure:"warehouse"`
}

func (conf *Config) Validate() error {
	err := validator.Validate(
		validator.DefaultConfig,
		validator.StringOneOf("CONFIG.STORAGE.TOPIC", conf.Topic, []string{constants.TopicMessage, constants.TopicRequest, constants.TopicResponse}),
	)
	if err != nil {
		return err
	}

	if err := conf.Warehouse.Validate(); err != nil {
		return err
	}
	return nil
}

type StorageWarehouse struct {
	Put StorageWarehousePut `json:"put" yaml:"put" mapstructure:"put"`
}

func (conf *StorageWarehouse) Validate() error {
	if err := conf.Put.Validate(); err != nil {
		return err
	}
	return nil
}

type StorageWarehousePut struct {
	Timeout int64 `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
	Size    int   `json:"size" yaml:"size" mapstructure:"size"`
}

func (conf *StorageWarehousePut) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThanOrEqual("CONFIG.STORAGE.WAREHOUSE.PUT.TIMEOUT", conf.Timeout, 1000),
		validator.NumberGreaterThan("CONFIG.STORAGE.WAREHOUSE.PUT.SIZE", conf.Size, 0),
	)
}
