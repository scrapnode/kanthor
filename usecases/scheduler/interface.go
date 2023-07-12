package scheduler

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
)

type Scheduler interface {
	patterns.Connectable
	Request() Request
}

type Request interface {
	Arrange(ctx context.Context, req *RequestArrangeReq) (*RequestArrangeRes, error)
}

type RequestArrangeReq struct {
	Message entities.Message
}

type RequestArrangeRes struct {
	Entities    []structure.BulkRes[entities.Request]
	FailKeys    []string
	SuccessKeys []string
}
