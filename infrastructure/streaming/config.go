package streaming

import "github.com/go-playground/validator/v10"

type ConnectionConfig struct {
	Uri    string       `json:"uri" yaml:"uri" mapstructure:"uri" validate:"required,uri"`
	Stream StreamConfig `json:"stream" yaml:"stream" mapstructure:"stream" validate:"required"`
}

func (conf *ConnectionConfig) Validate() error {
	return validator.New().Struct(conf)
}

type StreamConfig struct {
	Name     string   `json:"name" yaml:"name" mapstructure:"name" validate:"required,alphanumunicode"`
	Replicas int      `json:"replicas" yaml:"replicas" mapstructure:"replicas" validate:"number,gte=0"`
	Subjects []string `json:"subjects" yaml:"subjects" mapstructure:"subjects" validate:"required,gt=0,dive,required"`
	Limits   struct {
		Msgs     int64 `json:"msgs" yaml:"msgs" mapstructure:"msgs" validate:"required,number,gte=0"`
		MsgBytes int32 `json:"msg_bytes" yaml:"msg_bytes" mapstructure:"msg_bytes" validate:"required,number,gte=0"`
		Bytes    int64 `json:"bytes" yaml:"bytes" mapstructure:"bytes" validate:"required,number,gte=0"`
		Age      int64 `json:"age" yaml:"age" mapstructure:"age" validate:"required,number,gte=0"`
	} `json:"limits" yaml:"limits" mapstructure:"limits" validate:"required"`
}

func (conf *StreamConfig) Validate() error {
	return validator.New().Struct(conf)
}
