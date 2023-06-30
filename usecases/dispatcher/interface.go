package dispatcher

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
)

type Dispatcher interface {
	patterns.Connectable
	SendRequest(ctx context.Context, req *SendRequestsReq) (*SendRequestsRes, error)
}

type SendRequestsReq struct {
	Request entities.Request
}

type SendRequestsRes struct {
	Response entities.Response
}
