package dataplane

import (
	"context"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/usecases/dataplane/repos"
	"time"
)

func (usecase *message) Put(ctx context.Context, req *MessagePutReq) (*MessagePutRes, error) {
	cacheKey := cache.Key("APP_WITH_WORKSPACE", req.AppId)
	app, err := cache.Warp(usecase.cache, cacheKey, time.Hour, func() (*repos.ApplicationWithWorkspace, error) {
		usecase.meter.Count("cache_miss_total", 1, metric.Label("source", "dataplane_message_put"))
		return usecase.repos.Application().GetWithWorkspace(ctx, req.AppId)
	})
	if err != nil {
		return nil, err
	}

	msg := transformMessagePutReq2Message(app.Workspace.Tier.Name, req, usecase.timer, usecase.conf)
	msg.Metadata[entities.MetaTier] = app.Workspace.Tier.Name

	event, err := transformMessage2Event(msg)
	if err != nil {
		return nil, err
	}

	if err := usecase.publisher.Pub(ctx, event); err != nil {
		return nil, err
	}

	res := transformMessage2MessagePutRes(msg)
	return res, nil
}

func transformMessagePutReq2Message(tier string, req *MessagePutReq, timer timer.Timer, conf *config.Config) *entities.Message {
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
		streaming.Namespace,
		msg.Tier,
		streaming.TopicMsg,
		event.AppId,
		event.Type,
	)

	return event, nil
}

func transformMessage2MessagePutRes(msg *entities.Message) *MessagePutRes {
	return &MessagePutRes{
		Id:        msg.Id,
		Timestamp: msg.Timestamp,
		Bucket:    msg.Bucket,
	}
}
