package db

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
)

type Workspace interface {
	Get(ctx context.Context, id string) (*entities.Workspace, error)
}
