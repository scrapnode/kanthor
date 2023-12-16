package repositories

import (
	"context"

	"github.com/scrapnode/kanthor/internal/entities"
)

type Message interface {
	Create(ctx context.Context, docs []*entities.Message) ([]string, error)
}
