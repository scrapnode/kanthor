package coordinator

import (
	"fmt"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
)

func New(conf *Config, logger logging.Logger) (Coordinator, error) {
	if conf.Engine == EngineNats {
		return NewNats(conf, logger), nil
	}

	return nil, fmt.Errorf("coordinator: unknown engine")

}

type Coordinator interface {
	patterns.Connectable

	Send(cmd string, req Request) error
	Receive(handler func(cmd string, req []byte) error) error
}

type Request interface {
	Marshal() ([]byte, error)
	Unmarshal(data []byte) error
}
