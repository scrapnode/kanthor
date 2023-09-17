package validation

import govalidator "github.com/go-playground/validator/v10"

func NewPlayaround() Validator {
	return govalidator.New(govalidator.WithRequiredStructEnabled())
}
