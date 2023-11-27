package streaming

import (
	"context"

	natscore "github.com/nats-io/nats.go"
	"github.com/scrapnode/kanthor/logging"
	"github.com/scrapnode/kanthor/pkg/safe"
	"github.com/sourcegraph/conc/pool"
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
	returning := safe.Map[error]{}

	p := pool.New().WithMaxGoroutines(publisher.conf.Publisher.RateLimit)
	for refId, event := range events {
		if err := event.Validate(); err != nil {
			publisher.logger.Errorw("invalid event", "subject", event.Subject, "event_id", event.Id, "event", event.String())
			returning.Set(refId, err)
			continue
		}

		msg := NatsMsgFromEvent(event.Subject, event)
		p.Go(func() {
			ack, err := publisher.nats.js.PublishMsg(msg, natscore.Context(ctx), natscore.MsgId(event.Id))
			if err != nil {
				publisher.logger.Errorw("unable to publish message", "subject", event.Subject, "event_id", event.Id)
				returning.Set(refId, err)
				return
			}

			publisher.logger.Debugw("published message", "subject", event.Subject, "event_id", event.Id, "msg_seq", ack.Sequence)
		})
	}
	p.Wait()

	return returning.Data()
}
