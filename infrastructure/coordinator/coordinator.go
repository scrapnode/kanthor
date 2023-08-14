package coordinator

import (
	"encoding/json"
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

	Send(cmd *Command) error
	Receive(handler func(cmd *Command) error) error
}

var (
	CmdAuthzRefresh = "kanthor.coordinator.authz.refresh"
)

type Command struct {
	Name    string `json:"name"`
	Target  string `json:"target"`
	Request string `json:"request"`
}

func (cmd *Command) String() string {
	bytes, _ := json.Marshal(cmd)
	return string(bytes)
}
