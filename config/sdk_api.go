package config

import (
	"fmt"

	"github.com/scrapnode/kanthor/infrastructure/gateway"
)

type SdkApi struct {
	Gateway gateway.Config `json:"gateway" yaml:"gateway" mapstructure:"gateway"`
}

func (conf *SdkApi) Validate() error {
	if err := conf.Gateway.Validate(); err != nil {
		return fmt.Errorf("config.sdkapi.gateway: %v", err)
	}
	return nil
}
