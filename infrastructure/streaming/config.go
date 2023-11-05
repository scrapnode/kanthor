package streaming

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/scrapnode/kanthor/pkg/validator"
)

type Config struct {
	Name       string           `json:"name" yaml:"name" mapstructure:"name"`
	Uri        string           `json:"uri" yaml:"uri" mapstructure:"uri"`
	Nats       NatsConfig       `json:"nats" yaml:"nats" mapstructure:"nats"`
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

	uri, err := url.Parse(conf.Uri)
	if err != nil {
		return fmt.Errorf("streaming.conf.uri: unable to parse uri | %s", err.Error())
	}

	if strings.HasPrefix(uri.Scheme, "nats") {
		if err := conf.Nats.Validate(); err != nil {
			return err
		}
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
	Replicas int `json:"replicas" yaml:"replicas" mapstructure:"replicas"`
	Limits   struct {
		Size     int64 `json:"size" yaml:"size" mapstructure:"size"`
		MsgSize  int32 `json:"msg_size" yaml:"msg_size" mapstructure:"msg_size"`
		MsgCount int64 `json:"msg_count" yaml:"msg_count" mapstructure:"msg_count"`
		MsgAge   int64 `json:"msg_age" yaml:"msg_age" mapstructure:"msg_age"`
	} `json:"limits" yaml:"limits" mapstructure:"limits"`
}

func (conf *NatsConfig) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThanOrEqual("streaming.conf.nats.replicas", conf.Replicas, 0),
		validator.NumberGreaterThanOrEqual("streaming.conf.nats.limits.size", conf.Limits.Size, 0),
		validator.NumberGreaterThanOrEqual("streaming.conf.nats.limits.msg_size", conf.Limits.MsgSize, 0),
		validator.NumberGreaterThanOrEqual("streaming.conf.nats.limits.msg_bytes", conf.Limits.MsgCount, 0),
		validator.NumberGreaterThanOrEqual("streaming.conf.nats.limits.age", conf.Limits.MsgAge, 0),
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
	Throughput  int   `json:"throughput" yaml:"throughput" mapstructure:"throughput"`
}

func (conf *SubscriberConfig) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThanOrEqual("streaming.subscriber.config.max_deliver", conf.MaxRetry, 1),
		validator.NumberGreaterThan("streaming.conf.publisher.timeout", conf.Timeout, 1000),
		validator.NumberGreaterThan("streaming.conf.publisher.concurrency", conf.Concurrency, 0),
		validator.NumberGreaterThanOrEqual("streaming.conf.publisher.throughput", conf.Throughput, conf.Concurrency),
	)
}
