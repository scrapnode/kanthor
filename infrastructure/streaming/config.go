package streaming

import (
	"fmt"

	"github.com/scrapnode/kanthor/pkg/validator"
)

type Config struct {
	Name       string           `json:"name" yaml:"name" mapstructure:"name"`
	Uri        string           `json:"uri" yaml:"uri" mapstructure:"uri"`
	Nats       NatsConfig       `json:"stream" yaml:"stream" mapstructure:"stream"`
	Publisher  PublisherConfig  `json:"publisher" yaml:"publisher" mapstructure:"publisher"`
	Subscriber SubscriberConfig `json:"subscriber" yaml:"subscriber" mapstructure:"subscriber"`
}

func (conf *Config) Validate() error {
	err := validator.Validate(
		validator.DefaultConfig,
		validator.StringRequired("streaming.conf.name", conf.Name),
		validator.StringUri("streaming.conf.uri", conf.Uri),
	)
	if err != nil {
		return err
	}

	if err := conf.Nats.Validate(); err != nil {
		return err
	}

	if err := conf.Publisher.Validate(); err != nil {
		return err
	}

	if err := conf.Subscriber.Validate(); err != nil {
		return err
	}

	return nil
}

type NatsConfig struct {
	Replicas int      `json:"replicas" yaml:"replicas" mapstructure:"replicas"`
	Subjects []string `json:"subjects" yaml:"subjects" mapstructure:"subjects"`
	Limits   struct {
		Msgs     int64 `json:"msgs" yaml:"msgs" mapstructure:"msgs"`
		MsgBytes int32 `json:"msg_bytes" yaml:"msg_bytes" mapstructure:"msg_bytes"`
		Bytes    int64 `json:"bytes" yaml:"bytes" mapstructure:"bytes"`
		Age      int64 `json:"age" yaml:"age" mapstructure:"age"`
	} `json:"limits" yaml:"limits" mapstructure:"limits"`
}

func (conf *NatsConfig) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThanOrEqual("streaming.conf.nats.replicas", conf.Replicas, 0),
		validator.SliceRequired("streaming.conf.nats.subjects", conf.Subjects),
		validator.Array(conf.Subjects, func(i int, item *string) error {
			return validator.Validate(validator.DefaultConfig, validator.StringRequired(fmt.Sprintf("streaming.conf.nats.subjects[%d]", i), *item))
		}),
		validator.NumberGreaterThanOrEqual("streaming.conf.nats.limits.msgs", conf.Limits.Msgs, 0),
		validator.NumberGreaterThanOrEqual("streaming.conf.nats.limits.msg_bytes", int(conf.Limits.MsgBytes), 0),
		validator.NumberGreaterThanOrEqual("streaming.conf.nats.limits.bytes", int(conf.Limits.Bytes), 0),
		validator.NumberGreaterThanOrEqual("streaming.conf.nats.limits.age", int(conf.Limits.Age), 0),
	)
}

type PublisherConfig struct {
	Timeout   int64 `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
	RateLimit int   `json:"rate_limit" yaml:"rate_limit" mapstructure:"rate_limit"`
}

func (conf *PublisherConfig) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThan("streaming.conf.publisher.timeout", conf.Timeout, 1000),
		validator.NumberGreaterThan("streaming.conf.publisher.rate_limit", conf.RateLimit, 0),
	)
}

type SubscriberConfig struct {
	// MaxRetry is how many times we should try to re-deliver message if we get any error
	MaxRetry    int   `json:"max_retry" yaml:"max_retry" mapstructure:"max_retry"`
	Timeout     int64 `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
	Concurrency int   `json:"concurrency" yaml:"concurrency" mapstructure:"concurrency"`
}

func (conf *SubscriberConfig) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThan("streaming.conf.publisher.concurrency", conf.Concurrency, 0),
		validator.NumberGreaterThan("streaming.conf.publisher.timeout", conf.Timeout, 1000),
		validator.NumberGreaterThanOrEqual("streaming.subscriber.config.max_deliver", conf.MaxRetry, 1),
	)
}
