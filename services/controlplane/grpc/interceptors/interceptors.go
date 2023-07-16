package interceptors

import (
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	gatewayinterceptors "github.com/scrapnode/kanthor/infrastructure/gateway/grpc/interceptors"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/services/controlplane/grpc/interceptors/authz"
	"github.com/scrapnode/kanthor/services/controlplane/grpc/interceptors/workspace"
	usecase "github.com/scrapnode/kanthor/usecases/controlplane"
)

func WithWorkspace(logger logging.Logger, uc usecase.Controlplane) gatewayinterceptors.InterceptorOps {
	return func(ics *gatewayinterceptors.Interceptors) {
		ics.Unary = append(ics.Unary, workspace.UnaryServerInterceptor(logger, uc))
		ics.Stream = append(ics.Stream, workspace.StreamServerInterceptor(logger, uc))
	}
}

func WithAuthz(logger logging.Logger, engine authorizator.Authorizator, protected map[string]bool) gatewayinterceptors.InterceptorOps {
	return func(ics *gatewayinterceptors.Interceptors) {
		ics.Unary = append(ics.Unary, authz.UnaryServerInterceptor(logger, engine, protected))
		ics.Stream = append(ics.Stream, authz.StreamServerInterceptor(logger, engine, protected))
	}
}
