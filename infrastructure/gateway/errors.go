package gateway

func NewError(msg string) *Error {
	return &Error{Error: msg}
}

type Error struct {
	Error string `json:"error"`
	Code  string `json:"code,omitempty"`
}
