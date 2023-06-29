package logging

type Config struct {
	Pretty bool              `json:"pretty" mapstructure:"pretty"`
	Level  string            `json:"level" mapstructure:"level"`
	With   map[string]string `json:"with" mapstructure:"with"`
}
