package streaming

import (
	"context"

	natscore "github.com/nats-io/nats.go"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/pkg/safe"
	"github.com/scrapnode/kanthor/telemetry"
	"github.com/sourcegraph/conc/pool"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type NatsPublisher struct {
	name   string
	conf   *Config
	logger logging.Logger

	nats *nats
}

func (publisher *NatsPublisher) Name() string {
	return publisher.name
}

func (publisher *NatsPublisher) Pub(ctx context.Context, events map[string]*Event) map[string]error {
	tracer := ctx.Value(telemetry.CtxTracer).(trace.Tracer)
	attributes := trace.WithAttributes(
		attribute.String("streaming.publisher.engine", "nats"),
		attribute.Int("streaming.publisher.event_count", len(events)),
	)
	multictx, multispan := tracer.Start(ctx, "streaming.publisher.pub", attributes)
	defer func() {
		multispan.End()
	}()

	datac := make(chan map[string]error, 1)
	defer close(datac)

	go func() {
		returning := safe.Map[error]{}
		p := pool.New().WithMaxGoroutines(publisher.conf.Publisher.RateLimit)
		for refId, e := range events {
			if err := e.Validate(); err != nil {
				publisher.logger.Errorw("INFRASTRUCTURE.STREAMING.PUBLISHER.NATS.EVENT_VALIDATION.ERROR", "subject", e.Subject, "event_id", e.Id, "event", e.String())
				returning.Set(refId, err)
				continue
			}

			// store the value to use in p.Go, otherwise we got the same value
			event := e
			msg := NatsMsgFromEvent(e.Subject, e)
			p.Go(func() {
				_, singlespan := tracer.Start(multictx, "streaming.publisher.pub.event",
					trace.WithAttributes(
						attribute.String("streaming.publisher.event.id", event.Id),
						attribute.String("streaming.publisher.event.subject", event.Subject),
					))
				defer func() {
					singlespan.End()
				}()

				ack, err := publisher.nats.js.PublishMsg(msg, natscore.Context(ctx), natscore.MsgId(event.Id))
				if err != nil {
					publisher.logger.Errorw("INFRASTRUCTURE.STREAMING.PUBLISHER.NATS.EVENT_PUBLISH.ERROR", "subject", event.Subject, "event_id", event.Id)
					returning.Set(refId, err)
					return
				}

				if ack.Duplicate {
					publisher.logger.Errorw("INFRASTRUCTURE.STREAMING.PUBLISHER.NATS.EVENT_DUPLICATED.ERROR", "subject", event.Subject, "event_id", event.Id)
				}
			})
		}
		p.Wait()

		datac <- returning.Data()
	}()

	select {
	case data := <-datac:
		return data
	case <-ctx.Done():
		data := map[string]error{}
		for refId := range events {
			data[refId] = ctx.Err()
		}
		return data
	}
}
