package pipeline

import (
	"context"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
)

func UseValidation() Middleware {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return func(next Pipeline) Pipeline {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if err := validate.Struct(request); err != nil {
				return nil, err
			}

			return next(ctx, request)
		}
	}
}
