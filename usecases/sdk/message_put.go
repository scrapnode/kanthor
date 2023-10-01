package sdk

import (
	"context"
	"time"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/pkg/validator"
	"github.com/scrapnode/kanthor/usecases/transformation"
)

type MessagePutReq struct {
	AppId string
	Type  string

	Body     []byte
	Headers  entities.Header
	Metadata entities.Metadata
}

func (req *MessagePutReq) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("app_id", req.AppId, "app_"),
		validator.StringRequired("type", req.Type),
		validator.SliceRequired("body", req.Body),
	)
}

type MessagePutRes struct {
	Msg *entities.Message
}

func (uc *message) Put(ctx context.Context, req *MessagePutReq) (*MessagePutRes, error) {
	ws := ctx.Value(authorizator.CtxWs).(*entities.Workspace)

	key := CacheKeyApp(ws.Id, req.AppId)
	app, err := cache.Warp(uc.cache, ctx, key, time.Hour*24, func() (*entities.Application, error) {
		uc.metrics.Count(ctx, "cache_miss_total", 1)

		return uc.repos.Application().Get(ctx, ws, req.AppId)
	})
	if err != nil {
		return nil, err
	}

	wst := ctx.Value(authorizator.CtxWst).(*entities.WorkspaceTier)
	msg := &entities.Message{
		Tier:     wst.Name,
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
	msg.SetTS(uc.timer.Now(), uc.conf.Bucket.Layout)

	event, err := transformation.EventFromMessage(msg)
	if err != nil {
		return nil, err
	}

	if err := uc.publisher.Pub(ctx, event); err != nil {
		return nil, err
	}

	res := &MessagePutRes{Msg: msg}
	return res, nil
}
