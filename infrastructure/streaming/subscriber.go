package streaming

import (
	"context"

	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"github.com/scrapnode/kanthor/pkg/validator"
)

func NewSubscriber(conf *SubscriberConfig, logger logging.Logger) (Subscriber, error) {
	return NewNatsSubscriber(conf, logger), nil
}

type Subscriber interface {
	patterns.Connectable
	Sub(ctx context.Context, handler SubHandler) error
}

type SubHandler func(events []*Event) map[string]error

type SubscriberConfig struct {
	Connection ConnectionConfig `json:"connection" yaml:"connection" mapstructure:"connection"`

	// common config
	Name string `json:"name" yaml:"name" mapstructure:"name"`
	// only consume matched event with given subject
	FilterSubject string `json:"filter_subject" yaml:"filter_subject" mapstructure:"filter_subject"`

	// advance config, don't change it until you know what you are doing

	// MaxRetry is how many times we should try to re-deliver message if we get any error
	MaxRetry int `json:"max_retry" yaml:"max_retry" mapstructure:"max_retry"`
	// if MaxWaiting is 1, and more than one sub.Fetch actions, we will get an error
	MaxWaiting int `json:"max_waiting" yaml:"max_waiting" mapstructure:"max_waiting"`
	// if MaxAckPending is 1, and we are processing 1 message already
	// then we are going to request 2, we will only get 1
	MaxAckPending         int `json:"max_ack_pending" yaml:"max_ack_pending" mapstructure:"max_ack_pending"`
	MaxAckWaitingDuration int `json:"max_ack_wating_duration" yaml:"max_ack_wating_duration" mapstructure:"max_ack_wating_duration"`
	// if MaxRequestBatch is 10, that means for each instance you can only fetch maximum 10 messages
	// for instance, MaxRequestBatch=10 instances=5 -> each instance will receive 10 msgs -> total msgs were transfered is 50 msgs
	MaxRequestBatch int `json:"max_request_batch" yaml:"max_request_batch" mapstructure:"max_request_batch"`
}

func (conf *SubscriberConfig) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringRequired("streaming.subscriber.config.name", conf.Name),
		validator.StringRequired("streaming.subscriber.config.filter_subject", conf.FilterSubject),
		validator.NumberGreaterThanOrEqual("streaming.subscriber.config.max_deliver", conf.MaxRetry, 1),
		validator.NumberGreaterThanOrEqual("streaming.subscriber.config.max_waiting", conf.MaxWaiting, 1),
		validator.NumberGreaterThanOrEqual("streaming.subscriber.config.max_ack_pending", conf.MaxAckPending, 1),
		validator.NumberGreaterThanOrEqual("streaming.subscriber.config.max_ack_wating_duration", conf.MaxAckWaitingDuration, 1),
		validator.NumberGreaterThanOrEqual("streaming.subscriber.config.max_request_batch", conf.MaxRequestBatch, 1),
	)
}
