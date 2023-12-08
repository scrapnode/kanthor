package authenticator

import (
	"context"
	"fmt"
	"sync"
)

var (
	HeaderAuthCredentials = "authorization"
	HeaderAuthEngine      = "x-authorization-engine"
)

type Request struct {
	Credentials string
	Metadata    map[string]string
}

type Verifier interface {
	Verify(ctx context.Context, request *Request) (*Account, error)
}

type Authenticator interface {
	Register(engine string, verifier Verifier) error
	Authenticate(engine string, ctx context.Context, request *Request) (*Account, error)
}

func New() (Authenticator, error) {
	return &authenticator{strategies: map[string]Verifier{}}, nil
}

type authenticator struct {
	mu         sync.Mutex
	strategies map[string]Verifier
}

func (instance *authenticator) Register(engine string, verifier Verifier) error {
	instance.mu.Lock()
	defer instance.mu.Unlock()

	if _, has := instance.strategies[engine]; has {
		return fmt.Errorf("AUTHENTICATOR.ENGINE.ALREADY_REGISTERED")
	}

	instance.strategies[engine] = verifier
	return nil
}

func (instance *authenticator) Authenticate(engine string, ctx context.Context, request *Request) (*Account, error) {
	verifier, has := instance.strategies[engine]
	if !has {
		return nil, fmt.Errorf("AUTHENTICATOR.ENGINE.UNKNOWN")
	}

	return verifier.Verify(ctx, request)
}
