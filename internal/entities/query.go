package entities

import (
	"github.com/scrapnode/kanthor/gateway"
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
	Start int64
	End   int64
}

func (q *ScanningQuery) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.NumberInRange("limit", q.Limit, 5, 100),
		validator.NumberGreaterThan("start", q.Start, 0),
		validator.NumberGreaterThan("end", q.End, q.Start),
	)
}

func ScanningQueryFromGatewayQuery(query *gateway.Query) *ScanningQuery {
	if query == nil {
		return &ScanningQuery{}
	}

	return &ScanningQuery{Limit: query.Limit, Start: query.Start, End: query.End}
}
