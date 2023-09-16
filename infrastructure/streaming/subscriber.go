package streaming

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
)

var (
	SubModelPush = "push"
	SubModelPull = "pull"
)

func NewSubscriber(conf *SubscriberConfig, logger logging.Logger) (Subscriber, error) {
	if conf.BasedModel == SubModelPush {
		return NewNatsSubscriberPushing(conf, logger), nil
	}
	if conf.BasedModel == SubModelPull {
		return NewNatsSubscriberPulling(conf, logger), nil
	}

	return nil, fmt.Errorf("subscriber: unknown based model")
}

type Subscriber interface {
	patterns.Connectable
	Sub(ctx context.Context, handler SubHandler) error
}

type SubHandler func(events []Event) map[string]error

type SubscriberConfig struct {
	Connection ConnectionConfig `json:"connection" yaml:"connection" mapstructure:"connection" validate:"required"`
	BasedModel string           `json:"based_model" yaml:"based_model" mapstructure:"based_model" validate:"required,oneof=push pull"`

	// push-specific
	Push *SubscriberConfigPush `json:"push" yaml:"push" mapstructure:"push" validate:"-"`
	// pull-specific`
	Pull *SubscriberConfigPull `json:"pull" yaml:"pull" mapstructure:"pull" validate:"-"`

	Name string `json:"name" yaml:"name" mapstructure:"name" validate:"required"`
	// only consume matched event with given subject
	FilterSubject string `json:"filter_subject" yaml:"filter_subject" mapstructure:"filter_subject" validate:"required"`

	// Advance configuration, don't change it until you know what you are doing
	// MaxDelivery is how many times we should try to re-deliver message if we get any error
	MaxDeliver int `json:"max_deliver" yaml:"max_deliver" mapstructure:"max_deliver" validate:"required,number,gte=0"`
	// @TODO: consider apply RateLimit
}

func (conf *SubscriberConfig) Validate() error {
	if err := validator.New().Struct(conf); err != nil {
		return err
	}

	if conf.BasedModel == SubModelPush {
		if conf.Push == nil {
			return errors.New("streaming.subscriber: push config could not be nil")
		}
		if err := conf.Push.Validate(); err != nil {
			return err
		}
	}

	if conf.BasedModel == SubModelPull {
		if conf.Pull == nil {
			return errors.New("streaming.subscriber: pull config could not be nil")
		}
		if err := conf.Pull.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type SubscriberConfigPush struct {
	DeliverSubject string `json:"deliver_subject" yaml:"deliver_subject" mapstructure:"deliver_subject" validate:"required"`
	DeliverGroup   string `json:"deliver_group" yaml:"deliver_group" mapstructure:"deliver_group" validate:"required"`
	// Temporary is a config to allow us to create a temporary consumer that will be deleted after disconnected
	// this option is only available for Push-Based Model because Pull-Based Model requires consumer to be a durable one
	// must set it to TRUE explicitly to avoid misconfiguration
	Temporary bool `json:"temporary" yaml:"temporary" mapstructure:"temporary" validate:"boolean"`
}

func (conf *SubscriberConfigPush) Validate() error {
	if err := validator.New().Struct(conf); err != nil {
		return err
	}
	return nil
}

type SubscriberConfigPull struct {
	// if MaxWaiting is 1, and more than one sub.Fetch actions, we will get an error
	MaxWaiting int `json:"max_waiting" yaml:"max_waiting" mapstructure:"max_waiting" validate:"required,gt=0,lte=300"`
	// if MaxAckPending is 1, and we are processing 1 message already
	// then we are going to request 2, we will only get 1
	MaxAckPending int `json:"max_ack_pending" yaml:"max_ack_pending" mapstructure:"max_ack_pending" validate:"required,gt=0,lte=150000"`
	// if MaxRequestBatch is 1, and we are going to request 2, we will get an error
	MaxRequestBatch int `json:"max_request_batch" yaml:"max_request_batch" mapstructure:"max_request_batch" validate:"required,gt=0,lte=500"`
}

func (conf *SubscriberConfigPull) Validate() error {
	if err := validator.New().Struct(conf); err != nil {
		return err
	}
	return nil
}
