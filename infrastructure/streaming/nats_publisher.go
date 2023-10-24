package streaming

import (
	"context"
	"time"

	natscore "github.com/nats-io/nats.go"
	"github.com/scrapnode/kanthor/infrastructure/logging"
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

	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*time.Duration(publisher.conf.Publisher.Timeout))
	defer cancel()

	p := pool.New().WithMaxGoroutines(publisher.conf.Publisher.RateLimit)
	for key, event := range events {
		if err := event.Validate(); err != nil {
			returning.Set(key, err)
			continue
		}

		msg := NatsMsgFromEvent(event.Subject, event)
		p.Go(func() {
			ack, err := publisher.nats.js.PublishMsg(msg, natscore.Context(ctx), natscore.MsgId(event.Id))
			if err != nil {
				returning.Set(key, err)
				return
			}

			publisher.logger.Debugw("published message", "msg_seq", ack.Sequence)
		})
	}
	p.Wait()

	return returning.Data()
}
