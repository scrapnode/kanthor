package streaming

import (
	"fmt"

	"github.com/scrapnode/kanthor/pkg/validator"
)

type ConnectionConfig struct {
	Uri    string       `json:"uri" yaml:"uri" mapstructure:"uri"`
	Stream StreamConfig `json:"stream" yaml:"stream" mapstructure:"stream"`
}

func (conf *ConnectionConfig) Validate() error {
	err := validator.Validate(validator.DefaultConfig, validator.StringUri("streaming.conf.uri", conf.Uri))
	if err != nil {
		return err
	}

	if err := conf.Stream.Validate(); err != nil {
		return err
	}

	return nil
}

type StreamConfig struct {
	Name     string   `json:"name" yaml:"name" mapstructure:"name"`
	Replicas int      `json:"replicas" yaml:"replicas" mapstructure:"replicas"`
	Subjects []string `json:"subjects" yaml:"subjects" mapstructure:"subjects"`
	Limits   struct {
		Msgs     int64 `json:"msgs" yaml:"msgs" mapstructure:"msgs"`
		MsgBytes int32 `json:"msg_bytes" yaml:"msg_bytes" mapstructure:"msg_bytes"`
		Bytes    int64 `json:"bytes" yaml:"bytes" mapstructure:"bytes"`
		Age      int64 `json:"age" yaml:"age" mapstructure:"age"`
	} `json:"limits" yaml:"limits" mapstructure:"limits"`
}

func (conf *StreamConfig) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringRequired("streaming.conf.stream.name", conf.Name),
		validator.NumberGreaterThanOrEqual("streaming.conf.stream.replicas", conf.Replicas, 0),
		validator.SliceRequired("streaming.conf.stream.subjects", conf.Subjects),
		validator.Array(conf.Subjects, func(i int, item *string) error {
			return validator.Validate(validator.DefaultConfig, validator.StringRequired(fmt.Sprintf("streaming.conf.stream.subjects[%d]", i), *item))
		}),
		validator.NumberGreaterThanOrEqual("streaming.conf.stream.limits.msgs", conf.Limits.Msgs, 0),
		validator.NumberGreaterThanOrEqual("streaming.conf.stream.limits.msg_bytes", int(conf.Limits.MsgBytes), 0),
		validator.NumberGreaterThanOrEqual("streaming.conf.stream.limits.bytes", int(conf.Limits.Bytes), 0),
		validator.NumberGreaterThanOrEqual("streaming.conf.stream.limits.age", int(conf.Limits.Age), 0),
	)
}
