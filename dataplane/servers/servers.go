package servers

import (
	"github.com/scrapnode/kanthor/dataplane/config"
	"github.com/scrapnode/kanthor/dataplane/servers/grpc"
	"github.com/scrapnode/kanthor/dataplane/usecases/message"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
)

type Servers interface {
	patterns.Runnable
}

func New(
	conf *config.Config,
	logger logging.Logger,
	message message.Service,
) Servers {
	return grpc.New(conf, logger, message)
}
