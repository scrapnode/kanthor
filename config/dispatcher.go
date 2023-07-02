package config

import (
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/circuitbreaker"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/sender"
)

type Dispatcher struct {
	Publisher      streaming.PublisherConfig  `json:"publisher" mapstructure:"publisher"`
	Subscriber     streaming.SubscriberConfig `json:"subscriber" mapstructure:"subscriber"`
	Sender         sender.Config              `json:"sender" mapstructure:"sender"`
	Cache          *cache.Config              `json:"cache" mapstructure:"cache"`
	CircuitBreaker circuitbreaker.Config      `json:"circuit_breaker" mapstructure:"circuit_breaker"`
	Metrics        Server                     `json:"metrics" mapstructure:"metrics"`
}
