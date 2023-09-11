package sdk

import (
	"context"
	"net/http"
	"time"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/authorizator"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/usecases/transformation"
)

func (uc *message) Put(ctx context.Context, req *MessagePutReq) (*MessagePutRes, error) {
	ws := ctx.Value(authorizator.CtxWs).(*entities.Workspace)
	key := CacheKeyApp(ws.Id, req.AppId)
	app, err := cache.Warp(uc.cache, ctx, key, time.Hour*24, func() (*entities.Application, error) {
		uc.metrics.Count(ctx, "cache_miss_total", 1)

		return uc.repos.Application().Get(ctx, ws.Id, req.AppId)
	})
	if err != nil {
		return nil, err
	}

	wst := ctx.Value(authorizator.CtxWst).(*entities.WorkspaceTier)
	msg := &entities.Message{
		AttId:    utils.ID("att"),
		Tier:     wst.Name,
		AppId:    app.Id,
		Type:     req.Type,
		Body:     req.Body,
		Headers:  entities.Header{Header: http.Header{}},
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
