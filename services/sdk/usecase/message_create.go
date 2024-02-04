package usecase

import (
	"context"

	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/internal/transformation"
	"github.com/scrapnode/kanthor/pkg/identifier"
	"github.com/scrapnode/kanthor/pkg/validator"
	"github.com/scrapnode/kanthor/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type MessageCreateIn struct {
	WsId  string
	Tier  string
	AppId string
	Type  string

	Body     string
	Headers  entities.Header
	Metadata entities.Metadata
}

func (in *MessageCreateIn) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("ws_id", in.WsId, entities.IdNsWs),
		validator.StringRequired("tier", in.Tier),
		validator.StringStartsWith("app_id", in.AppId, entities.IdNsApp),
		validator.StringRequired("type", in.Type),
		validator.StringRequired("body", in.Body),
	)
}

type MessageCreateOut struct {
	EventId string `json:"event_id"`
	Message *entities.Message
}

func (uc *message) Create(ctx context.Context, in *MessageCreateIn) (*MessageCreateOut, error) {
	attributes := trace.WithAttributes(attribute.String("app.id", in.AppId))
	subctx, span := ctx.Value(telemetry.CtxTracer).(trace.Tracer).Start(ctx, "usecase.message.create", attributes)
	defer func() {
		span.End()
	}()

	app, err := uc.repositories.Database().Application().Get(subctx, in.WsId, in.AppId)
	if err != nil {
		return nil, err
	}

	msg := &entities.Message{
		Tier:     in.Tier,
		AppId:    app.Id,
		Type:     in.Type,
		Body:     in.Body,
		Headers:  entities.Header{},
		Metadata: entities.Metadata{},
	}
	// must use merge function otherwise you will edit the original data
	msg.Headers.Merge(in.Headers)
	msg.Metadata.Merge(in.Metadata)
	msg.Id = identifier.New(entities.IdNsMsg)
	msg.SetTS(uc.infra.Timer.Now())

	event, err := transformation.EventFromMessage(msg)
	if err != nil {
		return nil, err
	}

	events := map[string]*streaming.Event{}
	events[event.Id] = event
	spanner := &telemetry.Spanner{
		Tracer:   ctx.Value(telemetry.CtxTracer).(trace.Tracer),
		Contexts: make(map[string]context.Context),
	}
	spanner.Contexts[event.Id] = subctx
	pubctx := context.WithValue(subctx, telemetry.CtxSpanner, spanner)
	if errs := uc.publisher.Pub(pubctx, events); len(errs) > 0 {
		return nil, errs[event.Id]
	}

	return &MessageCreateOut{Message: msg}, nil
}
