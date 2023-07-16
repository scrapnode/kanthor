package repos

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
)

type Application interface {
	Create(ctx context.Context, entity *entities.Application) (*entities.Application, error)
}
