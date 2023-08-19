package sdk

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/usecases/transformation"
	"time"
)

func (uc *message) Put(ctx context.Context, req *MessagePutReq) (*MessagePutRes, error) {
	ws := ctx.Value(authorizator.CtxWs).(*entities.Workspace)
	key := CacheKeyApp(ws.Id, req.AppId)
	app, err := cache.Warp(uc.cache, ctx, key, time.Hour*24, func() (*entities.Application, error) {
		uc.metrics.Count("cache_miss_total", 1)

		return uc.repos.Application().Get(ctx, ws.Id, req.AppId)
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
		Headers:  req.Headers,
		Metadata: req.Metadata,
	}
	msg.GenId()
	msg.SetTS(uc.timer.Now(), uc.conf.Bucket.Layout)
	msg.Metadata[entities.MetaMsgId] = msg.Id
	msg.Metadata[entities.MetaAttId] = utils.ID("att")

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
