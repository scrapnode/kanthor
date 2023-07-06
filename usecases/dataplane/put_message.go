package dataplane

import (
	"context"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/domain/constants"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/repositories"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/timer"
	"time"
)

func (usecase *dataplane) PutMessage(ctx context.Context, req *PutMessageReq) (*PutMessageRes, error) {
	cacheKey := cache.Key("APP_WITH_WORKSPACE", req.AppId)
	app, err := cache.Warp(usecase.cache, cacheKey, time.Hour, func() (*repositories.ApplicationWithWorkspace, error) {
		// @TODO: cache miss
		return usecase.repos.Application().GetWithWorkspace(ctx, req.AppId)
	})
	if err != nil {
		return nil, err
	}

	msg := transformPutMessageReq2Message(app.Workspace.Tier.Name, req, usecase.timer, usecase.conf)
	msg.Metadata[entities.MetaTier] = app.Workspace.Tier.Name

	event, err := transformMessage2Event(msg)
	if err != nil {
		return nil, err
	}

	if err := usecase.publisher.Pub(ctx, event); err != nil {
		return nil, err
	}

	res := transformMessage2PutMessageRes(msg)
	return res, nil
}

func transformPutMessageReq2Message(tier string, req *PutMessageReq, timer timer.Timer, conf *config.Config) *entities.Message {
	msg := &entities.Message{
		Tier:     tier,
		AppId:    req.AppId,
		Type:     req.Type,
		Headers:  req.Headers,
		Body:     req.Body,
		Metadata: req.Metadata,
	}

	msg.GenId()
	msg.SetTS(timer.Now(), conf.Bucket.Layout)

	return msg
}

func transformMessage2Event(msg *entities.Message) (*streaming.Event, error) {
	data, err := msg.Marshal()
	if err != nil {
		return nil, err
	}

	event := &streaming.Event{
		AppId:    msg.AppId,
		Type:     msg.Type,
		Data:     data,
		Metadata: map[string]string{},
	}
	event.GenId()
	event.Subject = streaming.Subject(
		constants.TopicNamespace,
		msg.Tier,
		constants.TopicMessage,
		event.AppId,
		event.Type,
	)

	return event, nil
}

func transformMessage2PutMessageRes(msg *entities.Message) *PutMessageRes {
	return &PutMessageRes{
		Id:        msg.Id,
		Timestamp: msg.Timestamp,
		Bucket:    msg.Bucket,
	}
}
