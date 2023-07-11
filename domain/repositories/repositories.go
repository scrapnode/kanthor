package repositories

import (
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/database"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"github.com/scrapnode/kanthor/pkg/timer"
)

func New(conf *database.Config, logger logging.Logger, timer timer.Timer) Repositories {
	return NewSql(conf, logger, timer)
}

type Repositories interface {
	patterns.Connectable
	Workspace() Workspace
	Application() Application
	Endpoint() Endpoint
	EndpointRule() EndpointRule
}

type ListReq struct {
	Cursor string
	Search string
}

type ListOps func(req *ListReq)

func WithListCursor(cursor string) ListOps {
	return func(req *ListReq) {
		req.Cursor = cursor
	}
}

func WithListSearch(search string) ListOps {
	return func(req *ListReq) {
		req.Search = search
	}
}

func ListReqBuild(opts []ListOps) ListReq {
	req := ListReq{}
	for _, opt := range opts {
		opt(&req)
	}
	return req
}

func ListResBuild[T any](res *ListRes[T]) *ListRes[T] {
	if len(res.Data) == 0 {
		return res
	}

	latest := any(res.Data[len(res.Data)-1])
	res.Cursor = latest.(entities.Entity).Id
	return res
}

type ListRes[T any] struct {
	Cursor string
	Data   []T
}
