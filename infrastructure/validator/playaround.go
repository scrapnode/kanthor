package validator

import govalidator "github.com/go-playground/validator/v10"

func NewPlayaround() Validator {
	return &playaround{}
}

type playaround struct {
	v *govalidator.Validate
}

func (validator *playaround) Struct(s any) error {
	validator.init()
	return validator.v.Struct(s)
}

func (validator *playaround) init() {
	if validator.v == nil {
		validator.v = govalidator.New()
	}
}
