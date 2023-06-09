package message

import (
	"context"
	"github.com/scrapnode/kanthor/dataplane/config"
	"github.com/scrapnode/kanthor/domain/constants"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/infrastructure/timer"
	"time"
)

func (service *service) Create(ctx context.Context, req *CreateReq) (*CreateRes, error) {
	app, err := service.repos.Application().Get(ctx, req.AppId)
	if err != nil {
		return nil, err
	}
	ws, err := service.repos.Workspace().Get(ctx, app.WorkspaceId)
	if err != nil {
		return nil, err
	}

	msg := TransformCreateReq2Message(req, service.timer, service.conf)
	msg.Metadata[constants.MetaKeyTier] = ws.Tier.Name

	subject, event := TransformMessage2Event(msg, ws.Tier.Name)
	if err := service.publisher.Pub(ctx, subject, event); err != nil {
		return nil, err
	}

	res := TransformMessage2CreateRes(msg)
	return res, nil
}

type CreateReq struct {
	AppId string `json:"app_id"`
	Type  string `json:"type"`
	Body  string `json:"body"`
}

type CreateRes struct {
	Id        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Bucket    string    `json:"bucket"`
}

func TransformCreateReq2Message(req *CreateReq, timer timer.Timer, conf *config.Config) *entities.Message {
	msg := &entities.Message{
		AppId:    req.AppId,
		Type:     req.Type,
		Body:     []byte(req.Body),
		Metadata: map[string]string{},
	}
	msg.GenId()
	msg.SetTS(timer.Now(), conf.Dataplane.Message.BucketLayout)

	return msg
}

func TransformMessage2Event(msg *entities.Message, tierName string) (string, *streaming.Event) {
	event := &streaming.Event{
		AppId:    msg.AppId,
		Type:     msg.Type,
		Data:     msg.Body,
		Metadata: msg.Metadata,
	}
	event.GenId()
	subject := streaming.Subject(
		constants.TopicNamespace,
		tierName,
		TopicMessage,
		event.AppId,
		event.Type,
	)

	return subject, event
}

func TransformMessage2CreateRes(msg *entities.Message) *CreateRes {
	return &CreateRes{
		Id:        msg.Id,
		Timestamp: time.UnixMilli(msg.Timestamp),
		Bucket:    msg.Bucket,
	}
}
