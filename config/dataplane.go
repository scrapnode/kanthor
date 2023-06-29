package config

type Dataplane struct {
	GRPC    Server `json:"grpc" mapstructure:"grpc"`
	Metrics Server `json:"metrics" mapstructure:"metrics"`
}
