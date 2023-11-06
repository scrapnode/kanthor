package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/services/dispatcher/config"
)

type Forwarder interface {
	Send(ctx context.Context, in *ForwarderSendIn) (*ForwarderSendOut, error)
}

type forwarder struct {
	conf      *config.Config
	logger    logging.Logger
	infra     *infrastructure.Infrastructure
	publisher streaming.Publisher
}
