package controlplane

import (
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/enforcer"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/services"
	"github.com/scrapnode/kanthor/services/controlplane/grpc"
	usecase "github.com/scrapnode/kanthor/usecases/controlplane"
)

func New(
	conf *config.Config,
	logger logging.Logger,
	authenticator authenticator.Authenticator,
	meter metric.Meter,
	uc usecase.Controlplane,
	enforcer enforcer.Enforcer,
) services.Service {
	logger = logger.With("service", "controlplane")
	return grpc.New(conf, logger, authenticator, meter, enforcer, uc)
}
