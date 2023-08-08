package repos

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
)

type Message interface {
	Create(ctx context.Context, docs []entities.Message) ([]entities.TSEntity, error)
}
