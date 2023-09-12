package coordinator

import (
	"context"
	"errors"
	"sync"

	natscore "github.com/nats-io/nats.go"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/utils"
)

func NewNats(conf *Config, logger logging.Logger) Coordinator {
	return &nats{conf: conf, logger: logger}
}

type nats struct {
	id     string
	conf   *Config
	logger logging.Logger

	mu           sync.Mutex
	conn         *natscore.Conn
	subscription *natscore.Subscription
}

func (coordinator *nats) Connect(ctx context.Context) error {
	coordinator.mu.Lock()
	defer coordinator.mu.Unlock()

	coordinator.id = utils.ID("coordinator")

	conn, err := streaming.NewNats(streaming.ConnectionConfig{Uri: coordinator.conf.Nats.Uri}, coordinator.logger)
	if err != nil {
		return err
	}
	coordinator.conn = conn

	coordinator.logger.Info("connected")
	return nil
}

func (coordinator *nats) Disconnect(ctx context.Context) error {
	coordinator.mu.Lock()
	defer coordinator.mu.Unlock()

	if coordinator.subscription.IsValid() {
		if err := coordinator.subscription.Unsubscribe(); err != nil {
			return err
		}
	}
	coordinator.subscription = nil

	if !coordinator.conn.IsDraining() {
		if err := coordinator.conn.Drain(); err != nil {
			coordinator.logger.Error(err)
		}
	}
	if !coordinator.conn.IsClosed() {
		coordinator.conn.Close()
	}
	coordinator.conn = nil

	coordinator.logger.Info("disconnected")
	return nil
}

func (coordinator *nats) Send(ctx context.Context, cmd string, req Request) error {
	data, err := req.Marshal()
	if err != nil {
		return err
	}

	msg := &natscore.Msg{
		Header:  natscore.Header{},
		Subject: coordinator.conf.Nats.Subject,
		Data:    data,
	}
	msg.Header.Set(natscore.MsgIdHdr, utils.ID("coord"))
	msg.Header.Set(HeaderNodeId, coordinator.id)
	msg.Header.Set(HeaderCmd, cmd)

	coordinator.logger.Debugw("sending", "msg", utils.Stringify(msg))
	return coordinator.conn.PublishMsg(msg)
}

func (coordinator *nats) Receive(handle func(cmd string, req []byte) error) error {
	subscription, err := coordinator.conn.Subscribe(coordinator.conf.Nats.Subject, func(msg *natscore.Msg) {
		nodeId := msg.Header.Get(HeaderNodeId)
		if nodeId == "" {
			coordinator.logger.Errorw("ignore empty node id message", "msg", utils.Stringify(msg))
			return
		}
		if nodeId == coordinator.id {
			coordinator.logger.Debugw("ignore published node", "msg", utils.Stringify(msg))
			return
		}

		cmd := msg.Header.Get(HeaderCmd)
		if cmd == "" {
			coordinator.logger.Errorw("ignore empty cmd name message", "msg", utils.Stringify(msg))
			return
		}

		if err := handle(cmd, msg.Data); err != nil {
			coordinator.logger.Errorw(err.Error(), "msg", utils.Stringify(msg))
			if err := msg.Nak(); err != nil {
				coordinator.logger.Errorw(err.Error(), "msg", utils.Stringify(msg))
			}
			return
		}

		if err := msg.Ack(); err != nil && !errors.Is(err, natscore.ErrMsgNoReply) {
			coordinator.logger.Errorw(err.Error(), "msg", utils.Stringify(msg))
		}
	})
	if err != nil {
		return err
	}

	coordinator.subscription = subscription
	coordinator.logger.Infow("receiving", "subject", coordinator.conf.Nats.Subject)
	return nil
}
