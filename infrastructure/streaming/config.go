package streaming

import "github.com/go-playground/validator/v10"

type ConnectionConfig struct {
	Uri    string       `json:"uri" mapstructure:"uri" validate:"required,uri"`
	Stream StreamConfig `json:"stream" mapstructure:"stream" validate:"required"`
}

func (cconf *ConnectionConfig) Validate() error {
	return validator.New().Struct(cconf)
}

type StreamConfig struct {
	Name     string   `json:"name" mapstructure:"name" validate:"required,alphanumunicode"`
	Replicas int      `json:"replicas" mapstructure:"replicas" validate:"number,gte=0"`
	Subjects []string `json:"subjects" mapstructure:"subjects" validate:"required,gt=0,dive,required"`
	Limits   struct {
		Msgs     int64 `json:"msgs" mapstructure:"msgs" validate:"required,number,gte=0"`
		MsgBytes int32 `json:"msg_bytes" mapstructure:"msg_bytes" validate:"required,number,gte=0"`
		Bytes    int64 `json:"bytes" mapstructure:"bytes" validate:"required,number,gte=0"`
		Age      int64 `json:"age" mapstructure:"age" validate:"required,number,gte=0"`
	} `json:"limits" mapstructure:"limits" validate:"required"`
}

func (sconf *StreamConfig) Validate() error {
	return validator.New().Struct(sconf)
}
