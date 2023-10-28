package config

import (
	"fmt"

	"github.com/scrapnode/kanthor/infrastructure/configuration"
	"github.com/scrapnode/kanthor/pkg/validator"
)

// @TODO: mapstructure with env
func New(provider configuration.Provider) (*Config, error) {
	var conf Config
	return &conf, provider.Unmarshal(&conf)
}

type Config struct {
	Storage Storage `json:"storage" yaml:"storage" mapstructure:"storage"`
}

func (conf *Config) Validate() error {
	if err := conf.Storage.Validate(); err != nil {
		return err
	}
	return nil
}

type Storage struct {
	Warehouse StorageWarehouse `json:"warehouse" yaml:"warehouse" mapstructure:"warehouse"`
}

func (conf *Storage) Validate() error {
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
