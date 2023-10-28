package config

import (
	"fmt"

	"github.com/scrapnode/kanthor/configuration"
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
	Warehouse StorageWarehouse `json:"warehouse" yaml:"warehouse" mapstructure:"warehouse"`
}

func (conf *Config) Validate() error {
	if err := conf.Warehouse.Validate(); err != nil {
		return fmt.Errorf("config.storage.warehouse: %v", err)
	}
	return nil
}

type StorageWarehouse struct {
	Put StorageWarehousePut `json:"put" yaml:"put" mapstructure:"put"`
}

func (conf *StorageWarehouse) Validate() error {
	if err := conf.Put.Validate(); err != nil {
		return fmt.Errorf("config.storage.warehouse: %v", err)
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
		validator.NumberGreaterThanOrEqual("config.storage.forwarder.send.timeout", conf.Timeout, 1000),
		validator.NumberGreaterThan("config.storage.forwarder.send.size", conf.Size, 0),
	)
}
