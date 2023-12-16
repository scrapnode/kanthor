package ds

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
)

type Request interface {
	Scan(ctx context.Context, appId string, msgIds []string) (map[string]*entities.Request, error)
	ListByIds(ctx context.Context, maps map[string]map[string][]string) (map[string]*entities.Request, error)
}
