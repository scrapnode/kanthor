package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/samber/lo"
	"github.com/scrapnode/kanthor/assessor"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/transformation"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/safe"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/pkg/validator"
	"github.com/sourcegraph/conc"
)

type RequestScheduleIn struct {
	Timeout int64

	Messages map[string]*entities.Message
}

func ValidateRequestScheduleInMessage(prefix string, message *entities.Message) error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith(prefix+".id", message.Id, entities.IdNsMsg),
		validator.StringRequired(prefix+".tier", message.Tier),
		validator.StringStartsWith(prefix+".app_id", message.AppId, entities.IdNsApp),
		validator.StringRequired(prefix+".type", message.Type),
		validator.MapNotNil[string, string](prefix+".metadata", message.Metadata),
		validator.StringRequired(prefix+".body", message.Body),
	)
}

func (in *RequestScheduleIn) Validate() error {
	err := validator.Validate(
		validator.DefaultConfig,
		validator.MapRequired("messages", in.Messages),
		validator.NumberGreaterThan("timeout", in.Timeout, 1000),
	)
	if err != nil {
		return err
	}

	return validator.Validate(
		validator.DefaultConfig,
		validator.Map(in.Messages, func(refId string, item *entities.Message) error {
			prefix := fmt.Sprintf("messages.%s", refId)
			return ValidateRequestScheduleInMessage(prefix, item)
		}),
	)
}

type RequestScheduleOut struct {
	Success []string
	Error   map[string]error
}

func (uc *request) Schedule(ctx context.Context, in *RequestScheduleIn) (*RequestScheduleOut, error) {
	timeout, cancel := context.WithTimeout(ctx, time.Millisecond*time.Duration(in.Timeout))
	defer cancel()

	ok := &safe.Map[[]string]{}
	// ko must be map of message id and their error
	// so we can retry it if schedule requests of message got any error
	ko := &safe.Map[error]{}

	errc := make(chan error)
	defer close(errc)

	go func() {
		requests := uc.arrange(ctx, in.Messages)
		if len(requests) == 0 {
			errc <- nil
			return
		}

		maps := map[string]string{}
		events := map[string]*streaming.Event{}
		for _, request := range requests {
			key := utils.Key(request.AppId, request.MsgId, request.EpId, request.Id)
			event, err := transformation.EventFromRequest(&request)
			if err != nil {
				// un-recoverable error
				uc.logger.Errorw("could not transform request to event", "request", request.String())
				continue
			}
			events[key] = event
			maps[key] = request.MsgId
		}

		errs := uc.publisher.Pub(ctx, events)
		for refId := range events {
			// map key back to message id
			// to let system retry the message if any error was happen
			msgId := maps[refId]

			if err, ok := errs[refId]; ok {
				ko.Set(msgId, err)
				continue
			}

			if _, exist := ko.Get(msgId); !exist {
				ids, has := ok.Get(msgId)
				if !has {
					ids = []string{}
				}
				ok.Set(msgId, append(ids, msgId))
			}
		}

		errc <- nil
	}()

	select {
	case err := <-errc:
		return &RequestScheduleOut{Success: ok.Keys(), Error: ko.Data()}, err
	case <-timeout.Done():
		// context deadline exceeded, should set that error to remain messages
		for _, message := range in.Messages {
			if _, success := ok.Get(message.Id); success {
				// already success, should not retry it
				continue
			}

			// no error, should add context deadline error
			if _, has := ko.Get(message.Id); !has {
				ko.Set(message.Id, ctx.Err())
			}
		}
		return &RequestScheduleOut{Success: ok.Keys(), Error: ko.Data()}, nil
	}
}

func (uc *request) arrange(ctx context.Context, messages map[string]*entities.Message) []entities.Request {
	apps := map[string]string{}
	for _, message := range messages {
		apps[message.AppId] = message.AppId
	}
	appIds := lo.Keys(apps)
	applicables := uc.applicables(ctx, appIds)

	requests := []entities.Request{}
	for _, message := range messages {
		app, ok := applicables[message.AppId]
		if !ok {
			uc.logger.Warnw("could not find applicable assets for app", "app_id", message.AppId, "msg_id", message.Id)
			continue
		}

		reqs, logs := assessor.Requests(message, app, uc.infra.Timer)
		if len(logs) > 0 {
			for _, l := range logs {
				uc.logger.Warnw(l[0].(string), l[1:]...)
			}
		}
		requests = append(requests, reqs...)
	}

	if len(requests) == 0 {
		uc.logger.Warnw("no request was arranged", "app_id", appIds)
	}

	return requests
}

func (uc *request) applicables(ctx context.Context, appIds []string) map[string]*assessor.Assets {
	apps := &safe.Map[*assessor.Assets]{}

	var wg conc.WaitGroup
	for _, id := range appIds {
		appId := id
		wg.Go(func() {
			key := utils.Key("scheduler", appId)
			app, err := cache.Warp(uc.infra.Cache, ctx, key, time.Hour, func() (*assessor.Assets, error) {
				endpoints, err := uc.repositories.Endpoint().List(ctx, appId)
				if err != nil {
					return nil, err
				}
				returning := &assessor.Assets{EndpointMap: map[string]entities.Endpoint{}}
				for _, ep := range endpoints {
					returning.EndpointMap[ep.Id] = ep
				}

				rules, err := uc.repositories.Endpoint().Rules(ctx, appId)
				if err != nil {
					return nil, err
				}
				returning.Rules = rules

				return returning, nil
			})
			if err != nil {
				uc.logger.Errorw("unable to get applicable endpoints", "err", err.Error(), "app_id", appId)
				return
			}
			apps.Set(appId, app)
		})
	}
	wg.Wait()

	return apps.Data()
}
