package interceptors

import (
	"github.com/scrapnode/kanthor/infrastructure/enforcer"
	gatewayinterceptors "github.com/scrapnode/kanthor/infrastructure/gateway/grpc/interceptors"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/services/controlplane/grpc/interceptors/authz"
	usecase "github.com/scrapnode/kanthor/usecases/controlplane"
	grpccore "google.golang.org/grpc"
)

func Interceptors(logger logging.Logger, uc usecase.Controlplane, engine enforcer.Enforcer) gatewayinterceptors.Interceptors {
	mw := gatewayinterceptors.Interceptors{
		Unary: []grpccore.UnaryServerInterceptor{
			authz.UnaryServerInterceptor(logger, uc, engine),
		},
		Stream: []grpccore.StreamServerInterceptor{
			authz.StreamServerInterceptor(logger, uc, engine),
		},
	}
	return mw
}
