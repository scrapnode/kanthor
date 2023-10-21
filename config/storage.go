package config

import (
	"fmt"

	"github.com/scrapnode/kanthor/pkg/validator"
)

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
		validator.NumberGreaterThanOrEqual("config.dispatcher.forwarder.send.timeout", conf.Timeout, 1000),
		validator.NumberGreaterThan("config.dispatcher.forwarder.send.size", conf.Size, 0),
	)
}
