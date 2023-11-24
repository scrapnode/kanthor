package ds

import (
	"context"

	"github.com/scrapnode/kanthor/domain/entities"
)

type Response interface {
	Scan(ctx context.Context, appId string, msgIds []string) (map[string]entities.Response, error)
}
