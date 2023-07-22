package interceptors

import (
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	gatewayinterceptors "github.com/scrapnode/kanthor/infrastructure/gateway/grpc/interceptors"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/services/dataplane/grpc/interceptors/application"
	"github.com/scrapnode/kanthor/services/dataplane/grpc/interceptors/authz"
	usecase "github.com/scrapnode/kanthor/usecases/dataplane"
)

func WithApplication(logger logging.Logger, uc usecase.Dataplane) gatewayinterceptors.InterceptorOps {
	return func(ics *gatewayinterceptors.Interceptors) {
		ics.Unary = append(ics.Unary, application.UnaryServerInterceptor(logger, uc))
		ics.Stream = append(ics.Stream, application.StreamServerInterceptor(logger, uc))
	}
}

func WithAuthz(logger logging.Logger, engine authorizator.Authorizator, protected map[string]bool) gatewayinterceptors.InterceptorOps {
	return func(ics *gatewayinterceptors.Interceptors) {
		ics.Unary = append(ics.Unary, authz.UnaryServerInterceptor(logger, engine, protected))
		ics.Stream = append(ics.Stream, authz.StreamServerInterceptor(logger, engine, protected))
	}
}
