package servers

import (
	"github.com/scrapnode/kanthor/dataplane/config"
	"github.com/scrapnode/kanthor/dataplane/servers/grpc"
	"github.com/scrapnode/kanthor/infrastructure/database"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/msgbroker"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
)

type Servers interface {
	patterns.Runnable
}

func New(
	conf *config.Config,
	logger logging.Logger,
	msgbroker msgbroker.MsgBroker,
	database database.Database,
) (Servers, error) {
	return grpc.New(conf, logger)
}
