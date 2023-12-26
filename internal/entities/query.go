package entities

import (
	"time"

	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type PagingQuery struct {
	Search string
	Limit  int
	Page   int
	Ids    []string
}

func (q *PagingQuery) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringLenIfNotEmpty("search", q.Search, 3, 100),
		validator.NumberInRange("limit", q.Limit, 5, 100),
		validator.NumberInRange("page", q.Page, 0, 100),
	)
}

func PagingQueryFromGatewayQuery(query *gateway.Query) *PagingQuery {
	if query == nil {
		return &PagingQuery{}
	}

	limit := utils.MinInt(utils.MaxInt(query.Limit, 5), 100)
	page := utils.MinInt(utils.MaxInt(query.Page, 0), 100)
	return &PagingQuery{Search: query.Search, Limit: limit, Page: page, Ids: query.Id}
}

type ScanningQuery struct {
	Limit int
	From  time.Time
	To    time.Time
}

func (q *ScanningQuery) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberInRange("limit", q.Limit, 5, 100),
		validator.NumberGreaterThan("start", q.From.UnixMilli(), 0),
		validator.NumberGreaterThan("end", q.To.UnixMilli(), q.From.UnixMilli()),
	)
}

func ScanningQueryFromGatewayQuery(query *gateway.Query, timer timer.Timer) *ScanningQuery {
	if query == nil {
		return &ScanningQuery{}
	}

	return &ScanningQuery{Limit: query.Limit, From: timer.UnixMilli(query.Start), To: timer.UnixMilli(query.End)}
}
