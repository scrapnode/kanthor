package interceptors

import (
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	gatewayinterceptors "github.com/scrapnode/kanthor/infrastructure/gateway/grpc/interceptors"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/services/controlplane/grpc/interceptors/authz"
)

func WithAuthz(logger logging.Logger, engine authorizator.Authorizator) gatewayinterceptors.InterceptorOps {
	return func(ics *gatewayinterceptors.Interceptors) {
		ics.Unary = append(ics.Unary, authz.UnaryServerInterceptor(logger, engine))
		ics.Stream = append(ics.Stream, authz.StreamServerInterceptor(logger, engine))
	}
}
