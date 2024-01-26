package validator

import "errors"

// Check https://github.com/asaskevich/govalidator/blob/master/validator.go
// for more validation method if you need to implement others

type Validator interface {
	Validate() error
}

type Fn func() error

func Validate(conf *Config, fns ...Fn) (err error) {
	for _, fn := range fns {
		if suberr := fn(); suberr != nil {
			if conf.StopAtFirstError {
				return suberr
			}

			err = errors.Join(err, suberr)
		}
	}
	return
}
