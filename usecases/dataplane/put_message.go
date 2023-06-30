package dataplane

import (
	"context"
	"github.com/scrapnode/kanthor/config"
	"github.com/scrapnode/kanthor/domain/constants"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/infrastructure/timer"
)

func (usecase *dataplane) PutMessage(ctx context.Context, req *PutMessageReq) (*PutMessageRes, error) {
	app, err := usecase.repos.Application().Get(ctx, req.AppId)
	if err != nil {
		return nil, err
	}
	ws, err := usecase.repos.Workspace().Get(ctx, app.WorkspaceId)
	if err != nil {
		return nil, err
	}

	msg := transformPutMessageReq2Message(req, usecase.timer, usecase.conf)
	msg.Metadata[constants.MetaKeyTier] = ws.Tier.Name

	event, err := transformMessage2Event(msg, ws.Tier.Name)
	if err != nil {
		return nil, err
	}

	if err := usecase.publisher.Pub(ctx, event); err != nil {
		return nil, err
	}

	res := transformMessage2PutMessageRes(msg)
	return res, nil
}

func transformPutMessageReq2Message(req *PutMessageReq, timer timer.Timer, conf *config.Config) *entities.Message {
	msg := &entities.Message{
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

func transformMessage2Event(msg *entities.Message, tierName string) (*streaming.Event, error) {
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
		tierName,
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
