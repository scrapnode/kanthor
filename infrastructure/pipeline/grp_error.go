package pipeline

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/scrapnode/kanthor/domain/structure"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

func UseGRPCError(logger logging.Logger) Middleware {
	return func(next Pipeline) Pipeline {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			response, err = next(ctx, request)
			if err == nil {
				return response, nil
			}
			logger.Errorw(err.Error(), "request", utils.Stringify(request))

			// validation error
			if errs, ok := err.(validator.ValidationErrors); ok {
				s := structure.M{"error": http.StatusText(http.StatusBadRequest)}
				props := structure.M{}
				for _, err := range errs {
					props[err.Field()] = err.ActualTag()
				}
				s["props"] = props
				return nil, status.Error(codes.InvalidArgument, utils.Stringify(s))
			}

			s := structure.M{"error": http.StatusText(http.StatusInternalServerError)}
			return nil, status.Error(codes.Internal, utils.Stringify(s))
		}
	}
}
