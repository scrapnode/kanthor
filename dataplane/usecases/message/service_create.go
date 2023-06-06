package message

import (
	"context"
	"github.com/scrapnode/kanthor/domain/constants"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"time"
)

func (service *service) Create(ctx context.Context, req *CreateReq) (*CreateRes, error) {
	message := &entities.Message{
		AppId:    req.AppId,
		Type:     req.Type,
		Body:     []byte(req.Body),
		Metadata: map[string]string{},
	}
	message.GenId()
	message.GenBucket(service.conf.Dataplane.Message.BucketLayout)

	if req.Persistent {
		if _, err := service.repo.Put(ctx, message); err != nil {
			return nil, err
		}
	}

	event := &streaming.Event{
		AppId:    message.AppId,
		Type:     message.Type,
		Data:     message.Body,
		Metadata: message.Metadata,
	}
	event.GenId()
	subject := streaming.Subject(
		constants.Namespace,
		service.auth.Tier(),
		TopicMessage,
		event.AppId,
		event.Type,
	)
	if err := service.publisher.Pub(ctx, subject, event); err != nil {
		return nil, err
	}

	res := &CreateRes{
		Id:        message.Id,
		Timestamp: message.Timestamp,
		Bucket:    message.Bucket,
	}
	return res, nil
}

type CreateReq struct {
	AppId      string `json:"app_id"`
	Type       string `json:"type"`
	Body       string `json:"body"`
	Persistent bool   `json:"persistent"`
}

type CreateRes struct {
	Id        string     `json:"id"`
	Timestamp *time.Time `json:"timestamp"`
	Bucket    string     `json:"bucket"`
}
