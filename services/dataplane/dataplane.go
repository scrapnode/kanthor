package dataplane

import (
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure/authenticator"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/services"
	"github.com/scrapnode/kanthor/services/dataplane/grpc"
	usecase "github.com/scrapnode/kanthor/usecases/dataplane"
)

func New(
	conf *config.Config,
	logger logging.Logger,
	authenticator authenticator.Authenticator,
	authorizator authorizator.Authorizator,
	meter metric.Meter,
	uc usecase.Dataplane,
) services.Service {
	logger = logger.With("service", "dataplane")
	return grpc.New(conf, logger, authenticator, authorizator, meter, uc)
}
