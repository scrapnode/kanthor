package dispatcher

import (
	"context"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	usecase "github.com/scrapnode/kanthor/usecases/dispatcher"
)

func Consumer(logger logging.Logger, uc usecase.Dispatcher) streaming.SubHandler {
	// if you return error here, the event will be retried
	// so, you must test your error before return it
	return func(event *streaming.Event) error {
		logger.Debugw("received event", "event_id", event.Id)
		req, err := transformEventToRequest(event)
		if err != nil {
			logger.Error(err)
			return nil
		}

		request := &usecase.SendRequestsReq{Request: *req}
		response, err := uc.SendRequest(context.TODO(), request)
		if err != nil {
			logger.Error(err)
			return nil
		}

		logger.Debugw("received response", "response_id", response.Response.Id, "response_status", response.Response.Status)
		return nil
	}
}

func transformEventToRequest(event *streaming.Event) (*entities.Request, error) {
	var req entities.Request
	if err := req.Unmarshal(event.Data); err != nil {
		return nil, err
	}
	return &req, nil
}
