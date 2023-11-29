package usecase

import (
	"context"
	"time"

	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/internal/domain/entities"
	"github.com/scrapnode/kanthor/internal/domain/transformation"
	"github.com/scrapnode/kanthor/pkg/suid"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type MessagePutIn struct {
	WsId  string
	Tier  string
	AppId string
	Type  string

	Body     string
	Headers  entities.Header
	Metadata entities.Metadata
}

func (in *MessagePutIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("ws_id", in.WsId, entities.IdNsWs),
		validator.StringRequired("tier", in.Tier),
		validator.StringStartsWith("app_id", in.AppId, entities.IdNsApp),
		validator.StringRequired("type", in.Type),
		validator.StringRequired("body", in.Body),
	)
}

type MessagePutOut struct {
	EventId string `json:"event_id"`
	Message *entities.Message
}

func (uc *message) Put(ctx context.Context, in *MessagePutIn) (*MessagePutOut, error) {
	key := CacheKeyApp(in.WsId, in.AppId)
	app, err := cache.Warp(uc.infra.Cache, ctx, key, time.Hour*24, func() (*entities.Application, error) {
		return uc.repositories.Application().Get(ctx, in.WsId, in.AppId)
	})
	if err != nil {
		return nil, err
	}

	msg := &entities.Message{
		Tier:     in.Tier,
		AppId:    app.Id,
		Type:     in.Type,
		Body:     in.Body,
		Headers:  entities.Header{},
		Metadata: entities.Metadata{},
	}
	// must use merge function otherwise you will edit the original data
	msg.Headers.Merge(in.Headers)
	msg.Metadata.Merge(in.Metadata)
	msg.Id = suid.New(entities.IdNsMsg)
	msg.SetTS(uc.infra.Timer.Now())

	event, err := transformation.EventFromMessage(msg)
	if err != nil {
		return nil, err
	}

	events := map[string]*streaming.Event{}
	events[event.Id] = event
	if errs := uc.publisher.Pub(ctx, events); len(errs) > 0 {
		return nil, errs[event.Id]
	}

	return &MessagePutOut{Message: msg}, nil
}
