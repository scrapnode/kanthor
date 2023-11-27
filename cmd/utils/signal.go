package utils

import (
	"context"
	"os"
	"os/signal"
)

type SignalCtx struct {
	context.Context

	cancel  context.CancelFunc
	signals []os.Signal
	ch      chan os.Signal
	Fired   os.Signal // <- add this
}

func (c *SignalCtx) stop() {
	c.cancel()
	signal.Stop(c.ch)
}

func NotifyContext(parent context.Context, signals ...os.Signal) (ctx context.Context, stop context.CancelFunc) {
	ctx, cancel := context.WithCancel(parent)
	c := &SignalCtx{
		Context: ctx,
		cancel:  cancel,
		signals: signals,
	}
	c.ch = make(chan os.Signal, 1)
	signal.Notify(c.ch, c.signals...)
	if ctx.Err() == nil {
		go func() {
			select {
			case fired := <-c.ch:
				c.Fired = fired // <- add this
				c.cancel()
			case <-c.Done():
			}
		}()
	}
	return c, c.stop
}
