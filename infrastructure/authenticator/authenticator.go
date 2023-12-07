package authenticator

import (
	"context"
	"fmt"
	"sync"

	"github.com/scrapnode/kanthor/logging"
)

var (
	HeaderAuth          = "authorization"
	HeaderAuthScheme    = "x-authorization-scheme"
	HeaderAuthWorkspace = "x-authorization-workspace"
)

type Request struct {
	Credentials string
	Metadata    map[string]string
}

type Authenticate func(ctx context.Context, request *Request) (*Account, error)

type Authenticator interface {
	Register(scheme string, authenticate Authenticate) error
	Authenticate(ctx context.Context, scheme string, request *Request) (*Account, error)
}

func New(conf *Config, logger logging.Logger) (Authenticator, error) {
	return &authenticator{
		conf:       conf,
		logger:     logger,
		strategies: map[string]Authenticate{},
	}, nil
}

type authenticator struct {
	conf   *Config
	logger logging.Logger

	mu         sync.Mutex
	strategies map[string]Authenticate
}

func (instance *authenticator) Register(scheme string, authenticate Authenticate) error {
	instance.mu.Lock()
	defer instance.mu.Unlock()

	if _, has := instance.strategies[scheme]; has {
		return fmt.Errorf("AUTHENTICATOR.SCHEME.ALREADY_REGISTERED")
	}

	instance.strategies[scheme] = authenticate
	return nil
}

func (instance *authenticator) Authenticate(ctx context.Context, scheme string, request *Request) (*Account, error) {
	authenticate, has := instance.strategies[scheme]
	if !has {
		return nil, fmt.Errorf("AUTHENTICATOR.SCHEME.UNKNOWN")
	}

	return authenticate(ctx, request)
}
