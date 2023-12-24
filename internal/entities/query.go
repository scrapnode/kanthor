package entities

import (
	"github.com/scrapnode/kanthor/gateway"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type Query struct {
	Search string
	Limit  int
	Page   int
	Ids    []string
}

func (q *Query) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringLenIfNotEmpty("search", q.Search, 3, 100),
		validator.NumberInRange("limit", q.Limit, 5, 100),
		validator.NumberInRange("page", q.Page, 0, 100),
	)
}

func QueryFromGateWay(query *gateway.Query) *Query {
	limit := utils.MinInt(utils.MaxInt(query.Limit, 5), 100)
	page := utils.MinInt(utils.MaxInt(query.Page, 0), 100)
	return &Query{Search: query.Search, Limit: limit, Page: page, Ids: query.Id}
}
