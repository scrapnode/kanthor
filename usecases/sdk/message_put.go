package sdk

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/utils"
	"time"
)

func (uc *message) Put(ctx context.Context, req *MessagePutReq) (*MessagePutRes, error) {
	ws := ctx.Value(CtxWs).(*entities.Workspace)
	key := utils.Key(ws.Id, req.AppId)
	app, err := cache.Warp(uc.cache, ctx, key, time.Hour*24, func() (*entities.Application, error) {
		return uc.repos.Application().Get(ctx, ws.Id, req.AppId)
	})
	if err != nil {
		return nil, err
	}

	wst := ctx.Value(CtxWst).(*entities.WorkspaceTier)
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

	data, err := msg.Marshal()
	if err != nil {
		return nil, err
	}

	event := &streaming.Event{
		AppId:    msg.AppId,
		Type:     msg.Type,
		Id:       msg.Id,
		Data:     data,
		Metadata: map[string]string{},
	}
	event.Subject = streaming.Subject(
		streaming.Namespace,
		msg.Tier,
		streaming.TopicMsg,
		event.AppId,
		event.Type,
	)
	if err := uc.publisher.Pub(ctx, event); err != nil {
		return nil, err
	}

	res := &MessagePutRes{Msg: msg}
	return res, nil
}
