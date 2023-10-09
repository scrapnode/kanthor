package attempt

import (
	"context"
	"fmt"
	"time"

	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/ds"
	"github.com/scrapnode/kanthor/services/attempt/transformation"
	"github.com/sourcegraph/conc"
)

func RegisterTriggerCron(service *attempt) func() {
	key := "kanthor.services.attempt.trigger"
	duration := time.Duration(service.conf.Attempt.Trigger.Cron.LockDuration) * time.Second

	return func() {
		locker := service.locker(key, duration)
		ctx, cancel := context.WithTimeout(context.Background(), duration)
		defer cancel()

		if err := locker.Lock(ctx); err != nil {
			service.logger.Errorw("unable to acquire a lock", "key", key)
			return
		}
		defer func() {
			ctx, cancel := context.WithTimeout(context.Background(), duration)
			defer cancel()
			if err := locker.Unlock(ctx); err != nil {
				service.logger.Errorw("unable to release a lock", "key", key)
			}
		}()

		ucreq := transformation.ApplicationTriggerReq(service.conf.Attempt.Trigger.Cron.ScanSize, service.conf.Attempt.Trigger.Cron.PublishSize)
		ucres, err := service.uc.Application().Trigger(ctx, ucreq)
		if err != nil {
			service.logger.Errorw("unable to get applications", "err", err.Error())
			return
		}

		if len(ucres.Error) > 0 {
			for key, err := range ucres.Error {
				service.logger.Errorw("unable to trigger application", "key", key, "err", err.Error())
			}
		}
	}
}

func RegisterTriggerConsumer(service *attempt) streaming.SubHandler {
	duration := time.Duration(service.conf.Attempt.Trigger.Consumer.LockDuration) * time.Second

	return func(events []*streaming.Event) map[string]error {
		errs := &ds.SafeMap[error]{}

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		from := service.timer.Now().Add(time.Duration(service.conf.Attempt.Trigger.Consumer.ScanFrom) * time.Second)
		to := service.timer.Now().Add(time.Duration(service.conf.Attempt.Trigger.Consumer.ScanTo) * time.Second)

		var wg conc.WaitGroup
		for _, e := range events {
			event := e
			wg.Go(func() {
				service.logger.Debugw("received event", "event_id", event.Id)
				app, err := transformation.EventToApplication(event)
				if err != nil {
					service.metrics.Count(ctx, "attempt_transform_error", 1)
					service.logger.Errorw(err.Error(), "event", event.String())
					// unable to parse message from event is considered as un-retriable error
					// ignore the error, and we need to check it manually with log
					return
				}
				key := fmt.Sprintf("kanthor.services.attempt.trigger.consumer/%s", app.Id)
				locker := service.locker(key, duration)
				ctx, cancel := context.WithTimeout(context.Background(), duration)
				defer cancel()

				if err := locker.Lock(ctx); err != nil {
					service.logger.Errorw("unable to acquire a lock | key:%s", key)
					// if we could not acquire the lock, don't need to retry so don't set error here
					return
				}
				defer func() {
					ctx, cancel := context.WithTimeout(context.Background(), duration)
					defer cancel()
					if err := locker.Unlock(ctx); err != nil {
						service.logger.Errorw("unable to release a lock | key:%s", key)
					}
				}()

				ucreq := transformation.ApplicationToTriggerScanReq(app, from, to)
				if err := ucreq.Validate(); err != nil {
					service.metrics.Count(ctx, "attempt_trigger_scan_error", 1)
					service.logger.Errorw(err.Error(), "event", event.String(), "app", app.String())
					// IMPORTANT: sometime under high-preasure, we got nil pointer
					// should retry this event anyway
					errs.Set(event.Id, err)
					return
				}

				scanned, err := service.uc.Trigger().Scan(ctx, ucreq)
				if err != nil {
					service.metrics.Count(ctx, "attempt_trigger_scan_error", 1)
					service.logger.Errorw(err.Error(), "evt_id", event.Id, "app_id", app.Id)
					errs.Set(event.Id, err)
					return
				}
				service.metrics.Count(ctx, "attempt_trigger_scan_to_schedule_total", int64(len(scanned.ToScheduleMsgIds)))
				service.metrics.Count(ctx, "attempt_trigger_scan_to_attempt_total", int64(len(scanned.ToAttemptReqs)))

				// schedule message again
				scheduled, err := service.uc.Trigger().Schedule(ctx, transformation.MsgIdsToTriggerScheduleReq(app, scanned.ToScheduleMsgIds))
				if err == nil {
					service.metrics.Count(ctx, "attempt_trigger_schedule_success", int64(len(scheduled.Success)))
					service.metrics.Count(ctx, "attempt_trigger_schedule_error", int64(len(scheduled.Error)))

					// @TODO: use dead-letter
					if len(scheduled.Error) > 0 {
						for key, err := range scheduled.Error {
							service.logger.Errorw("schedule error", "key", key, "err", err.Error())
						}
					}
				} else {
					service.metrics.Count(ctx, "attempt_trigger_schedule_error", 1)
					service.logger.Errorw(err.Error(), "evt_id", event.Id, "app_id", app.Id)
					errs.Set(event.Id, err)
				}

				// create attempts
				created, err := service.uc.Trigger().Create(ctx, transformation.RequestsToTriggerScheduleReq(scanned.ToAttemptReqs))
				if err == nil {
					service.metrics.Count(ctx, "attempt_trigger_create_success", int64(len(created.Success)))
					service.metrics.Count(ctx, "attempt_trigger_create_error", int64(len(created.Error)))

					// @TODO: use dead-letter
					if len(created.Error) > 0 {
						for key, err := range created.Error {
							service.logger.Errorw("create error", "key", key, "err", err.Error())
						}
					}
				} else {
					service.metrics.Count(ctx, "attempt_trigger_create_error", 1)
					service.logger.Errorw(err.Error(), "evt_id", event.Id, "app_id", app.Id)
					errs.Set(event.Id, err)
				}
			})
		}

		c := make(chan bool)
		defer close(c)

		go func() {
			wg.Wait()
			c <- true
		}()

		select {
		case <-c:
			if errs.Count() > 0 {
				service.metrics.Count(ctx, "attempt_send_error", int64(errs.Count()))
				service.logger.Errorw("encoutered errors", "error_count", errs.Count(), "error_sample", errs.Sample().Error())
			}
			service.logger.Infow("scheduled requests", "message_count", len(events), "request_count", len(events)-errs.Count())

			return errs.Data()
		case <-ctx.Done():
			// timeout, all events will be considered as failed
			// set non-error event with timeout error
			for _, event := range events {
				if err, ok := errs.Get(event.Id); !ok && err == nil {
					errs.Set(event.Id, ctx.Err())
				}
			}
			service.metrics.Count(ctx, "attempt_send_timeout_error", 1)
			service.metrics.Count(ctx, "attempt_send_error", int64(errs.Count()))
			service.logger.Errorw("encoutered errors", "error_count", errs.Count())
			return errs.Data()
		}
	}
}
