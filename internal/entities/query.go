package entities

import (
	"time"

	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/pkg/validator"
)

var DefaultPagingLimitMin = 5
var DefaultPagingLimitMax = 100
var DefaultPagingPageMin = 1
var DefaultPagingPageMax = 100

var DefaultPagingQuery = &PagingQuery{
	Limit: DefaultPagingLimitMin,
	Page:  DefaultPagingPageMin,
}

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
		validator.NumberInRange("limit", q.Limit, DefaultPagingLimitMin, DefaultPagingLimitMax),
		validator.NumberInRange("page", q.Page, DefaultPagingPageMin, DefaultPagingPageMax),
	)
}

func PagingQueryFromGatewayQuery(query *gateway.Query) *PagingQuery {
	if query == nil {
		return &PagingQuery{}
	}

	limit := utils.Min(utils.Max(query.Limit, DefaultPagingLimitMin), DefaultPagingLimitMax)
	page := utils.Min(utils.Max(query.Page, DefaultPagingPageMin), DefaultPagingPageMax)
	return &PagingQuery{Search: query.Search, Limit: limit, Page: page, Ids: query.Id}
}

type ScanningQuery struct {
	Search string
	Limit  int
	From   time.Time
	To     time.Time
}

func (q *ScanningQuery) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringLenIfNotEmpty("search", q.Search, 27, 100),
		validator.NumberInRange("limit", q.Limit, 5, 100),
		validator.NumberGreaterThan("start", q.From.UnixMilli(), 0),
	)
}

func ScanningQueryFromGatewayQuery(query *gateway.Query, timer timer.Timer) *ScanningQuery {
	if query == nil {
		return &ScanningQuery{}
	}

	q := &ScanningQuery{
		Search: query.Search,
		Limit:  query.Limit,
		From:   timer.UnixMilli(query.Start),
		To:     timer.Now(),
	}
	if timer.Now().UnixMilli() > query.End && query.End > 0 {
		q.To = timer.UnixMilli(query.End)
	}

	return q
}

type ScanningResult[T any] struct {
	Data  T
	Error error
}
