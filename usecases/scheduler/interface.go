package scheduler

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
)

type Scheduler interface {
	patterns.Connectable
	ArrangeRequests(ctx context.Context, req *ArrangeRequestsReq) (*ArrangeRequestsRes, error)
}

type ArrangeRequestsReq struct {
	Message entities.Message
}

type ArrangeRequestsRes struct {
	Entities    []structure.BulkRes[entities.Request]
	FailKeys    []string
	SuccessKeys []string
}
