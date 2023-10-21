package dispatcher

import (
	"context"

	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/infrastructure"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/sender"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
)

type Forwarder interface {
	Send(ctx context.Context, req *ForwarderSendReq) (*ForwarderSendRes, error)
}

type forwarder struct {
	conf      *config.Config
	logger    logging.Logger
	infra     *infrastructure.Infrastructure
	dispatch  sender.Send
	publisher streaming.Publisher
}
