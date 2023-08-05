package streaming

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
)

func NewSubscriber(conf *SubscriberConfig, logger logging.Logger) (Subscriber, error) {
	return NewNatsSubscriber(conf, logger), nil
}

type Subscriber interface {
	patterns.Connectable
	Sub(ctx context.Context, handler SubHandler) error
}

type SubHandler func(event *Event) error

type SubscriberConfig struct {
	Connection ConnectionConfig `json:"connection" yaml:"connection" mapstructure:"connection" validate:"required"`
	Name       string           `json:"name" yaml:"name" mapstructure:"name"`
	Topic      string           `json:"topic" yaml:"topic" mapstructure:"topic" validate:"required"`
	Group      string           `json:"group" yaml:"group" mapstructure:"group" validate:"required"`
	// only consume matched event with given subject
	FilterSubject string `json:"filter_subject" yaml:"filter_subject" mapstructure:"filter_subject" validate:"required"`

	// Advance configuration, don't change it until you know what you are doing
	// must set it to TRUE explicitly to avoid misconfiguration
	Temporary bool `json:"temporary" yaml:"temporary" mapstructure:"temporary" validate:"required,boolean"`
	// or we can call this option by MaxRetry
	MaxDeliver int `json:"max_deliver" yaml:"max_deliver" mapstructure:"max_deliver" validate:"required,number,gte=0"`
	// @TODO: consider apply RateLimit
}

func (conf *SubscriberConfig) Validate() error {
	return validator.New().Struct(conf)
}
