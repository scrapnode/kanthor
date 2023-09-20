package validator

// Check https://github.com/asaskevich/govalidator/blob/master/validator.go
// for more validation method if you need to implement others

type Fn func() error

func Validate(conf *Config, fns ...Fn) error {
	for _, fn := range fns {
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}
