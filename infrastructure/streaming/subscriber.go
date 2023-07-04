package streaming

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
)

func NewSubscriber(conf *SubscriberConfig, logger logging.Logger) Subscriber {
	return NewNatsSubscriber(conf, logger)
}

type Subscriber interface {
	patterns.Connectable
	Sub(ctx context.Context, handler SubHandler) error
}

type SubHandler func(event *Event) error

type SubscriberConfig struct {
	*ConnectionConfig
	Name  string `json:"name" mapstructure:"name"`
	Topic string `json:"topic" mapstructure:"topic" validate:"required"`
	Group string `json:"group" mapstructure:"group" validate:"required"`
	// only consume matched event with given subject
	FilterSubject string `json:"filter_subject" mapstructure:"filter_subject"`

	// Advance configuration, don't change it until you know what you are doing
	// must set it to TRUE explicitly to avoid misconfiguration
	Temporary bool `json:"temporary" mapstructure:"temporary" validate:"boolean"`
	// or we can call this option by MaxRetry
	MaxDeliver int `json:"max_deliver" mapstructure:"max_deliver" validate:"number,gte=0"`
	// @TODO: consider apply RateLimit
}

func (conf SubscriberConfig) Validate() error {
	return validator.New().Struct(conf)
}
