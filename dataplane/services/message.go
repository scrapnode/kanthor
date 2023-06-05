package services

import (
	"context"
	"github.com/scrapnode/kanthor/dataplane/config"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"strings"
	"time"
)

func NewMessage(conf *config.Config, logger logging.Logger, publisher streaming.Publisher) Message {
	logger = logger.With("component", "dataplane.services.message")
	return &message{conf: conf, logger: logger, publisher: publisher}
}

type Message interface {
	patterns.Connectable
	Create(ctx context.Context, req *MessageCreateReq) (*MessageCreateRes, error)
}

type MessageCreateReq struct {
	AppId      string `json:"app_id"`
	Type       string `json:"type"`
	Body       string `json:"body"`
	Persistent bool   `json:"persistent"`
}

type MessageCreateRes struct {
	Id        string     `json:"id"`
	Timestamp *time.Time `json:"timestamp"`
	Bucket    string     `json:"bucket"`
}

type message struct {
	conf      *config.Config
	logger    logging.Logger
	publisher streaming.Publisher
}

func (service *message) Connect(ctx context.Context) error {
	if err := service.publisher.Connect(ctx); err != nil {
		return err
	}

	service.logger.Info("connected")
	return nil
}

func (service *message) Disconnect(ctx context.Context) error {
	service.logger.Info("disconnected")

	if err := service.publisher.Disconnect(ctx); err != nil {
		return err
	}

	return nil
}

func (service *message) Create(ctx context.Context, req *MessageCreateReq) (*MessageCreateRes, error) {
	message := &entities.Message{
		AppId:    req.AppId,
		Type:     req.Type,
		Body:     []byte(req.Body),
		Metadata: map[string]string{},
	}
	message.GenId()
	message.GenBucket(service.conf.Dataplane.Message.BucketLayout)

	if req.Persistent {
		// @TODO: implement persistent logic here
	}

	event := &streaming.Event{
		AppId:    message.AppId,
		Type:     message.Type,
		Data:     message.Body,
		Metadata: message.Metadata,
	}
	event.GenId()
	subject := strings.Join([]string{
		// @TODO: don't hardcode here
		"kanthor",
		"default",
		"message",
		event.AppId,
		event.Type,
	}, ".")
	if err := service.publisher.Pub(ctx, subject, event); err != nil {
		return nil, err
	}

	res := &MessageCreateRes{
		Id:        message.Id,
		Timestamp: message.Timestamp,
		Bucket:    message.Bucket,
	}
	return res, nil
}
