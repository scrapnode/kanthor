package transformation

import (
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
)

func EventToMessage(event streaming.Event) (*entities.Message, error) {
	var msg entities.Message
	if err := msg.Unmarshal(event.Data); err != nil {
		return nil, err
	}
	return &msg, nil
}

func EventToRequest(event streaming.Event) (*entities.Request, error) {
	var req entities.Request
	if err := req.Unmarshal(event.Data); err != nil {
		return nil, err
	}
	return &req, nil
}

func EventToResponse(event streaming.Event) (*entities.Response, error) {
	var req entities.Response
	if err := req.Unmarshal(event.Data); err != nil {
		return nil, err
	}
	return &req, nil
}
