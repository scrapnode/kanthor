package config

import (
	"github.com/scrapnode/kanthor/configuration"
	"github.com/scrapnode/kanthor/internal/constants"
	"github.com/scrapnode/kanthor/pkg/validator"
)

func New(provider configuration.Provider) (*Config, error) {
	provider.SetDefault("storage.topic", constants.TopicPublic)

	var conf Wrapper
	if err := provider.Unmarshal(&conf); err != nil {
		return nil, err
	}
	return &conf.Storage, nil
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
	Topic     string           `json:"topic" yaml:"topic" mapstructure:"topic"`
	Warehouse StorageWarehouse `json:"warehouse" yaml:"warehouse" mapstructure:"warehouse"`
}

func (conf *Config) Validate() error {
	err := validator.Validate(
		validator.DefaultConfig,
		validator.StringOneOf(
			"CONFIG.STORAGE.TOPIC",
			conf.Topic,
			[]string{
				constants.TopicMessage,
				constants.TopicRequest,
				constants.TopicResponse,
				constants.TopicPublic,
			},
		),
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
	Timeout   int64 `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
	BatchSize int   `json:"batch_size" yaml:"batch_size" mapstructure:"batch_size"`
}

func (conf *StorageWarehousePut) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThanOrEqual("CONFIG.STORAGE.WAREHOUSE.PUT.TIMEOUT", conf.Timeout, 1000),
		validator.NumberGreaterThan("CONFIG.STORAGE.WAREHOUSE.PUT.BATCH_SIZE", conf.BatchSize, 0),
	)
}
