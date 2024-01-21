package streaming

import (
	"errors"
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
		validator.StringRequired("INFRASTRUCTURE.STREAMING.CONFIG.NAME", conf.Name),
		validator.StringUri("INFRASTRUCTURE.STREAMING.CONFIG.URI", conf.Uri),
	)
	if err != nil {
		return err
	}

	uri, err := url.Parse(conf.Uri)
	if err != nil {
		return errors.New("INFRASTRUCTURE.STREAMING.CONFIG.URI_PARSE.ERROR")
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
		validator.NumberGreaterThanOrEqual("INFRASTRUCTURE.STREAMING.CONFIG.NATS.REPLICAS", conf.Replicas, 0),
		validator.NumberGreaterThanOrEqual("INFRASTRUCTURE.STREAMING.CONFIG.NATS.LIMITS.SIZE", conf.Limits.Size, 0),
		validator.NumberGreaterThanOrEqual("INFRASTRUCTURE.STREAMING.CONFIG.NATS.LIMITS.MSG_SIZE", conf.Limits.MsgSize, 0),
		validator.NumberGreaterThanOrEqual("INFRASTRUCTURE.STREAMING.CONFIG.NATS.LIMITS.MSG_COUNT", conf.Limits.MsgCount, 0),
		validator.NumberGreaterThanOrEqual("INFRASTRUCTURE.STREAMING.CONFIG.NATS.LIMITS.AGE", conf.Limits.MsgAge, 0),
	)
}

type PublisherConfig struct {
	RateLimit int `json:"rate_limit" yaml:"rate_limit" mapstructure:"rate_limit"`
}

func (conf *PublisherConfig) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThan("INFRASTRUCTURE.STREAMING.CONFIG.PUBLISHER.RATE_LIMIT", conf.RateLimit, 0),
	)
}

type SubscriberConfig struct {
	// MaxRetry is how many times we should try to re-deliver message if we get any error
	MaxRetry    int `json:"max_retry" yaml:"max_retry" mapstructure:"max_retry"`
	Concurrency int `json:"concurrency" yaml:"concurrency" mapstructure:"concurrency"`
	Throughput  int `json:"throughput" yaml:"throughput" mapstructure:"throughput"`
}

func (conf *SubscriberConfig) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThanOrEqual("INFRASTRUCTURE.STREAMING.CONFIG.SUBSCRIBER.MAX_RETRY", conf.MaxRetry, 1),
		validator.NumberGreaterThan("INFRASTRUCTURE.STREAMING.CONFIG.SUBSCRIBER.CONCURRENCY", conf.Concurrency, 0),
		validator.NumberGreaterThanOrEqual("INFRASTRUCTURE.STREAMING.CONFIG.SUBSCRIBER.THOUGHPUT", conf.Throughput, conf.Concurrency),
	)
}
