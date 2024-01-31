package authenticator

import (
	"context"
	"fmt"
	"sync"

	"github.com/scrapnode/kanthor/infrastructure/circuitbreaker"
	"github.com/scrapnode/kanthor/infrastructure/sender"
	"github.com/scrapnode/kanthor/logging"
)

var (
	HeaderAuthnCredentials = "authorization"
	HeaderAuthnEngine      = "x-authorization-engine"
)

type Request struct {
	Credentials string
	Metadata    map[string]string
}

type Verifier interface {
	Verify(ctx context.Context, request *Request) (*Account, error)
}

type Authenticator interface {
	Engines() []string
	Register(engine string, verifier Verifier) error
	Authenticate(engine string, ctx context.Context, request *Request) (*Account, error)
}

func New(conf []Config, logger logging.Logger, send sender.Send, cb circuitbreaker.CircuitBreaker) (Authenticator, error) {
	instance := &authenticator{strategies: map[string]Verifier{}}

	for _, c := range conf {
		if c.Engine == EngineAsk {
			ask, err := NewAsk(c.Ask)
			if err != nil {
				return nil, err
			}
			instance.Register(c.Engine, ask)
		}

		if c.Engine == EngineExternal {
			external, err := NewExternal(c.External, logger, send, cb)
			if err != nil {
				return nil, err
			}
			instance.Register(c.Engine, external)
		}
	}

	return instance, nil
}

type authenticator struct {
	mu         sync.Mutex
	engines    []string
	strategies map[string]Verifier
}

func (instance *authenticator) Engines() []string {
	return instance.engines
}

func (instance *authenticator) Register(engine string, verifier Verifier) error {
	instance.mu.Lock()
	defer instance.mu.Unlock()

	if _, has := instance.strategies[engine]; has {
		return fmt.Errorf("AUTHENTICATOR.ENGINE.ALREADY_REGISTERED")
	}

	instance.engines = append(instance.engines, engine)
	instance.strategies[engine] = verifier
	return nil
}

func (instance *authenticator) Authenticate(engine string, ctx context.Context, request *Request) (*Account, error) {
	verifier, has := instance.strategies[engine]
	if !has {
		return nil, fmt.Errorf("AUTHENTICATOR.ENGINE.UNKNOWN: %s (%v)", engine, instance.engines)
	}

	return verifier.Verify(ctx, request)
}
