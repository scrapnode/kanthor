package consumer

import (
	"context"
	"time"

	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/internal/transformation"
	"github.com/scrapnode/kanthor/services/recovery/usecase"
)

func Handler(service *consumer) streaming.SubHandler {
	// if you return error here, the event will be retried
	// so, you must test your error before return it
	return func(events map[string]*streaming.Event) map[string]error {
		recovery := map[string]*entities.Recovery{}
		for id, event := range events {
			recover, err := transformation.EventToRecovery(event)
			if err != nil {
				service.logger.Errorw(err.Error(), "event", event.String())
				// unable to parse recovery from event is considered as un-retriable error
				// ignore the error, and we need to check it manually with log
				continue
			}

			if err := usecase.ValidateScannerExecuteRecovery("recover", recover); err != nil {
				service.logger.Errorw(err.Error(), "event", event.String(), "recover", recover.String())
				// got malformed recovery, should ignore and not retry it
				continue
			}

			recovery[id] = recover
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(service.conf.Consumer.Timeout))
		defer cancel()

		in := &usecase.ScannerExecuteIn{
			RecoveryBatchSize: service.conf.Consumer.BatchSize,
			Recovery:          recovery,
		}
		// we alreay validated messages of request, don't need to validate again
		out, err := service.uc.Scanner().Execute(ctx, in)
		if err != nil {
			service.logger.Errorw("unable to execute recovery", "error", err.Error())

			retruning := map[string]error{}
			// got un-coverable error, should retry all event
			for _, event := range events {
				retruning[event.Id] = err
			}
			return retruning
		}

		service.logger.Infow("execute recovery success", "event_count", len(events), "ok_count", len(out.Success), "ko_count", len(out.Error))
		if len(out.Error) > 0 {
			for ref, err := range out.Error {
				service.logger.Errorw("execute got error", "ref", ref, "error", err.Error())
			}
		}

		return out.Error
	}
}
