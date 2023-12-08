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
		if ferr := fn(); ferr != nil {
			if !conf.StopAtFirstError {
				err = errors.Join(err, ferr)
				continue
			}

			return err
		}
	}
	return
}
