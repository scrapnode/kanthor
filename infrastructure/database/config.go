package database

type Config struct {
	Uri           string `json:"uri"`
	RetryAttempts uint   `json:"retry_attempts"`
}
