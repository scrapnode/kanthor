package dispatcher

import (
	"context"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/services/dispatcher/transformation"
	usecase "github.com/scrapnode/kanthor/usecases/dispatcher"
	"github.com/sourcegraph/conc"
	"time"
)

func Consumer(service *dispatcher) streaming.SubHandler {
	// if you return error here, the event will be retried
	// so, you must test your error before return it
	return func(events []streaming.Event) map[string]error {
		results := map[string]error{}

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		var wg conc.WaitGroup
		for _, e := range events {
			event := e
			wg.Go(func() {
				service.logger.Debugw("received event", "event_id", event.Id)
				req, err := transformation.EventToRequest(event)
				if err != nil {
					service.logger.Error(err)
					// next retry will be failed because of same error of transformation
					return
				}

				request := &usecase.ForwarderSendReq{Request: *req}
				response, err := service.uc.Forwarder().Send(ctx, request)
				if err != nil {
					service.logger.Errorw(err.Error(), "event", event.String())
					results[event.Id] = err
					return
				}
				// custom handler for error
				if response.Response.Error != "" {
					service.logger.Errorw(response.Response.Error, "event", event.String())
					results[event.Id] = err
					return
				}

				service.logger.Debugw("received response", "response_id", response.Response.Id, "response_status", response.Response.Status)
			})
		}

		c := make(chan bool)
		go func() {
			defer close(c)
			wg.Wait()
		}()

		select {
		case <-c:
			return results
		case <-ctx.Done():
			for _, event := range events {
				results[event.Id] = ctx.Err()
			}
			return results
		}
	}
}
