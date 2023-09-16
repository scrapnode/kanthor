package validator

import govalidator "github.com/go-playground/validator/v10"

func NewPlayaround() Validator {
	return &playaround{v: govalidator.New()}
}

type playaround struct {
	v *govalidator.Validate
}

func (validator *playaround) Struct(s any) error {
	return validator.v.Struct(s)
}
