package entities

import "github.com/scrapnode/kanthor/pkg/validator"

type Query struct {
	Search string
	Limit  int
	Page   int
}

func (q *Query) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringLenIfNotEmpty("search", q.Search, 3, 100),
		validator.NumberInRange("limit", q.Limit, 10, 100),
		validator.NumberInRange("page", q.Page, 0, 100),
	)
}
