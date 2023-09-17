package validation

func New() Validator {
	return NewPlayaround()
}

type Validator interface {
	Struct(s any) error
}
