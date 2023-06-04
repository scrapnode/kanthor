package logging

type Config struct {
	Debug bool              `json:"debug"`
	Level string            `json:"level"`
	With  map[string]string `json:"with"`
}
