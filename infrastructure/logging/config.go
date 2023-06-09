package logging

type Config struct {
	Pretty bool              `json:"pretty"`
	Level  string            `json:"level"`
	With   map[string]string `json:"with"`
}
