package sdkapi

func NewHttpError(message string) *HttpError {
	return &HttpError{Error: message}
}

type HttpError struct {
	Error string `json:"error"`
}
