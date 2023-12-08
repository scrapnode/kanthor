package validator

type Config struct {
	StopAtFirstError bool
}

var DefaultConfig = &Config{
	StopAtFirstError: false,
}
