package config

import (
	"fmt"

	"github.com/scrapnode/kanthor/infrastructure/streaming"
)

type Storage struct {
	Subscriber streaming.SubscriberConfig `json:"subscriber" yaml:"subscriber" mapstructure:"subscriber"`
}

func (conf *Storage) Validate() error {
	if err := conf.Subscriber.Validate(); err != nil {
		return fmt.Errorf("config.storage.subscriber: %v", err)
	}
	return nil
}
