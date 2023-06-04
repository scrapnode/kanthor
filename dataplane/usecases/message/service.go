package message

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
)

type Service interface {
	Put(ctx context.Context, message *entities.Message) (*entities.Message, error)
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

type service struct {
	repo Repository
}

func (s *service) Put(ctx context.Context, message *entities.Message) (*entities.Message, error) {
	message.GenId()
	message.GenBucket("2006010215")
	return s.repo.Put(ctx, message)
}
