package sdk

import (
	"context"
	"time"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/validator"
	"github.com/scrapnode/kanthor/usecases/transformation"
)

type MessagePutReq struct {
	WsId          string
	WorkspaceTier string
	AppId         string
	Type          string

	Body     []byte
	Headers  entities.Header
	Metadata entities.Metadata
}

func (req *MessagePutReq) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("ws_id", req.WsId, entities.IdNsWs),
		validator.StringRequired("workspace_tier", req.WorkspaceTier),
		validator.StringStartsWith("app_id", req.AppId, entities.IdNsApp),
		validator.StringRequired("type", req.Type),
		validator.SliceRequired("body", req.Body),
	)
}

type MessagePutRes struct {
	EventId string `json:"event_id"`
	Message *entities.Message
}

func (uc *message) Put(ctx context.Context, req *MessagePutReq) (*MessagePutRes, error) {
	key := CacheKeyApp(req.WsId, req.AppId)
	app, err := cache.Warp(uc.infra.Cache, ctx, key, time.Hour*24, func() (*entities.Application, error) {
		return uc.repos.Application().Get(ctx, req.WsId, req.AppId)
	})
	if err != nil {
		return nil, err
	}

	msg := &entities.Message{
		Tier:     req.WorkspaceTier,
		AppId:    app.Id,
		Type:     req.Type,
		Body:     req.Body,
		Headers:  entities.NewHeader(),
		Metadata: entities.Metadata{},
	}
	// must use merge function otherwise you will edit the original data
	msg.Headers.Merge(req.Headers)
	msg.Metadata.Merge(req.Metadata)

	msg.GenId()
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

	res := &MessagePutRes{Message: msg}
	return res, nil
}
