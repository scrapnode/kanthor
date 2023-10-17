package config

import (
	"fmt"

	"github.com/scrapnode/kanthor/infrastructure/gateway"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
)

type SdkApi struct {
	Gateway   gateway.Config            `json:"gateway" yaml:"gateway" mapstructure:"gateway"`
	Publisher streaming.PublisherConfig `json:"publisher" yaml:"publisher" mapstructure:"publisher"`
}

func (conf *SdkApi) Validate() error {
	if err := conf.Gateway.Validate(); err != nil {
		return fmt.Errorf("config.sdkapi.gateway: %v", err)
	}
	return nil
}
