package message

import (
	"context"
	"github.com/scrapnode/kanthor/infrastructure/database/entities"
)

type Service interface {
	Put(ctx context.Context, message *entities.Message) error
}

type service struct {
}
