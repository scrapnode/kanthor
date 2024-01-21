package gateway

func Error(err error) *Err {
	return &Err{Error: err.Error()}
}

func ErrorString(err string) *Err {
	return &Err{Error: err}
}

type Err struct {
	Error string `json:"error" yaml:"error"`
} // @name Error
