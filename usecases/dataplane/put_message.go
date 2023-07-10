package dataplane

import (
	"context"
	"fmt"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/repositories"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/timer"
	"time"
)

func (usecase *dataplane) PutMessage(ctx context.Context, req *PutMessageReq) (*PutMessageRes, error) {
	cacheKey := cache.Key("APP_WITH_WORKSPACE", req.AppId)
	app, err := cache.Warp(usecase.cache, cacheKey, time.Hour, func() (*repositories.ApplicationWithWorkspace, error) {
		usecase.meter.Count("cache_miss_total", 1, metric.Label("source", "dataplane_put_message"))
		return usecase.repos.Application().GetWithWorkspace(ctx, req.AppId)
	})
	if err != nil {
		usecase.logger.Errorw(err.Error(), "app_id", req.AppId)
		return nil, fmt.Errorf("unable to find application [%s]", req.AppId)
	}

	msg := transformPutMessageReq2Message(app.Workspace.Tier.Name, req, usecase.timer, usecase.conf)
	msg.Metadata[entities.MetaTier] = app.Workspace.Tier.Name

	event, err := transformMessage2Event(msg)
	if err != nil {
		usecase.logger.Errorw(err.Error(), "app_id", req.AppId, "msg_id", msg.Id)
		return nil, fmt.Errorf("unable transform message to event [%s/%s]", req.AppId, msg.Id)
	}

	if err := usecase.publisher.Pub(ctx, event); err != nil {
		usecase.logger.Errorw(err.Error(), "app_id", req.AppId, "msg_id", msg.Id)
		return nil, fmt.Errorf("unable to publish event for message [%s/%s]", req.AppId, msg.Id)
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
		streaming.Namespace,
		msg.Tier,
		streaming.TopicMsg,
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
