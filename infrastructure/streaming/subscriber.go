package streaming

import (
	"context"
	"errors"
	"fmt"

	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"github.com/scrapnode/kanthor/pkg/validator"
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

	return nil, fmt.Errorf("streaming.subscriber: unknown subscribe model")
}

type Subscriber interface {
	patterns.Connectable
	Sub(ctx context.Context, handler SubHandler) error
}

type SubHandler func(events []*Event) map[string]error

type SubscriberConfig struct {
	Connection ConnectionConfig `json:"connection" yaml:"connection" mapstructure:"connection"`
	BasedModel string           `json:"based_model" yaml:"based_model" mapstructure:"based_model"`

	// push-specific
	Push *SubscriberConfigPush `json:"push" yaml:"push" mapstructure:"push"`
	// pull-specific`
	Pull *SubscriberConfigPull `json:"pull" yaml:"pull" mapstructure:"pull"`

	Name string `json:"name" yaml:"name" mapstructure:"name"`
	// only consume matched event with given subject
	FilterSubject string `json:"filter_subject" yaml:"filter_subject" mapstructure:"filter_subject"`

	// Advance configuration, don't change it until you know what you are doing
	// MaxDelivery is how many times we should try to re-deliver message if we get any error
	MaxDeliver int `json:"max_deliver" yaml:"max_deliver" mapstructure:"max_deliver"`
	// @TODO: consider apply RateLimit
}

func (conf *SubscriberConfig) Validate() error {
	err := validator.Validate(
		validator.DefaultConfig,
		validator.StringOneOf("streaming.subscriber.config.based_model", conf.BasedModel, []string{SubModelPush, SubModelPull}),
		validator.StringRequired("streaming.subscriber.config.name", conf.Name),
		validator.StringRequired("streaming.subscriber.config.filter_subject", conf.FilterSubject),
		validator.NumberGreaterThanOrEqual("streaming.subscriber.config.max_deliver", conf.MaxDeliver, 0),
	)
	if err != nil {
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
	DeliverSubject string `json:"deliver_subject" yaml:"deliver_subject" mapstructure:"deliver_subject"`
	DeliverGroup   string `json:"deliver_group" yaml:"deliver_group" mapstructure:"deliver_group"`
	// Temporary is a config to allow us to create a temporary consumer that will be deleted after disconnected
	// this option is only available for Push-Based Model because Pull-Based Model requires consumer to be a durable one
	// must set it to TRUE explicitly to avoid misconfiguration
	Temporary bool `json:"temporary" yaml:"temporary" mapstructure:"temporary"`
}

func (conf *SubscriberConfigPush) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringRequired("streaming.subscriber.config.push.deliver_subject", conf.DeliverSubject),
		validator.StringRequired("streaming.subscriber.config.push.deliver_group", conf.DeliverGroup),
	)
}

type SubscriberConfigPull struct {
	// if MaxWaiting is 1, and more than one sub.Fetch actions, we will get an error
	MaxWaiting int `json:"max_waiting" yaml:"max_waiting" mapstructure:"max_waiting"`
	// if MaxAckPending is 1, and we are processing 1 message already
	// then we are going to request 2, we will only get 1
	MaxAckPending int `json:"max_ack_pending" yaml:"max_ack_pending" mapstructure:"max_ack_pending"`
	// if MaxRequestBatch is 1, and we are going to request 2, we will get an error
	MaxRequestBatch int `json:"max_request_batch" yaml:"max_request_batch" mapstructure:"max_request_batch"`
}

func (conf *SubscriberConfigPull) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberLessThanOrEqual("streaming.subscriber.config.pull.max_waiting", conf.MaxWaiting, 300),
		validator.NumberLessThanOrEqual("streaming.subscriber.config.pull.max_ack_pending", conf.MaxAckPending, 150000),
		validator.NumberLessThanOrEqual("streaming.subscriber.config.pull.max_request_batch", conf.MaxRequestBatch, 500),
	)
}
