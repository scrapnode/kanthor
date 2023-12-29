package usecase

import (
	"context"
	"fmt"

	"github.com/samber/lo"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/internal/assessor"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/internal/transformation"
	"github.com/scrapnode/kanthor/pkg/safe"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/pkg/validator"
	"github.com/sourcegraph/conc"
)

type RequestScheduleIn struct {
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
	ok := &safe.Map[[]string]{}
	ko := &safe.Map[error]{}

	// we have to store a ref map of messages.id and the key
	// so if we got any error, we can report back to the call that a key has a error
	eventIdRefs := map[string]string{}
	for eventId, msg := range in.Messages {
		eventIdRefs[msg.Id] = eventId
	}

	errc := make(chan error, 1)
	defer close(errc)

	go func() {
		requests := uc.arrange(ctx, in.Messages)
		if len(requests) == 0 {
			errc <- nil
			return
		}

		msgrefs := map[string]string{}
		events := map[string]*streaming.Event{}
		for _, request := range requests {
			msgRefId := utils.Key(request.AppId, request.MsgId, request.Id)
			event, err := transformation.EventFromRequest(request)
			if err != nil {
				// un-recoverable error
				uc.logger.Errorw("could not transform request to event", "request", request.String())
				continue
			}

			events[msgRefId] = event
			msgrefs[msgRefId] = request.MsgId
		}

		errs := uc.publisher.Pub(ctx, events)
		for msgRefId := range events {
			// map key back to message id
			msgId := msgrefs[msgRefId]
			eventRef := eventIdRefs[msgId]

			if err, ok := errs[msgRefId]; ok {
				ko.Set(eventRef, err)
				continue
			}

			if _, exist := ko.Get(msgId); !exist {
				ids, has := ok.Get(msgId)
				if !has {
					ids = []string{}
				}
				ok.Set(eventRef, append(ids, msgRefId))
			}
		}

		errc <- nil
	}()

	select {
	case err := <-errc:
		return &RequestScheduleOut{Success: ok.Keys(), Error: ko.Data()}, err
	case <-ctx.Done():
		// context deadline exceeded, should set that error to remain messages
		for _, msg := range in.Messages {
			eventRef := eventIdRefs[msg.Id]

			if _, success := ok.Get(msg.Id); success {
				// already success, should not retry it
				continue
			}

			// no error, should add context deadline error
			if _, has := ko.Get(msg.Id); !has {
				ko.Set(eventRef, ctx.Err())
				continue
			}
		}
		return &RequestScheduleOut{Success: ok.Keys(), Error: ko.Data()}, nil
	}
}

func (uc *request) arrange(ctx context.Context, messages map[string]*entities.Message) map[string]*entities.Request {
	apps := map[string]string{}
	for _, message := range messages {
		apps[message.AppId] = message.AppId
	}
	appIds := lo.Keys(apps)
	applicables := uc.applicables(ctx, appIds)

	returning := map[string]*entities.Request{}
	for _, message := range messages {
		app, ok := applicables[message.AppId]
		if !ok {
			uc.logger.Warnw("could not find applicable assets for app", "app_id", message.AppId, "msg_id", message.Id)
			continue
		}

		reqs, traces := assessor.Requests(message, app, uc.infra.Timer)
		if len(traces) > 0 {
			for _, trace := range traces {
				uc.logger.Warnw(trace[0].(string), trace[1:]...)
			}
		}
		for reqId, req := range reqs {
			returning[reqId] = req
		}
	}

	if len(returning) == 0 {
		uc.logger.Warnw("no request was arranged", "app_id", appIds)
	}

	return returning
}

func (uc *request) applicables(ctx context.Context, appIds []string) map[string]*assessor.Assets {
	apps := &safe.Map[*assessor.Assets]{}

	var wg conc.WaitGroup
	for _, id := range appIds {
		appId := id
		wg.Go(func() {
			endpoints, err := uc.repositories.Endpoint().List(ctx, appId)
			if err != nil {
				uc.logger.Errorw("unable to get applicable endpoints", "err", err.Error(), "app_id", appId)
				return
			}
			assets := &assessor.Assets{EndpointMap: map[string]entities.Endpoint{}}
			for _, ep := range endpoints {
				assets.EndpointMap[ep.Id] = ep
			}

			rules, err := uc.repositories.Endpoint().Rules(ctx, appId)
			if err != nil {
				uc.logger.Errorw("unable to get applicable endpoint rules", "err", err.Error(), "app_id", appId)
				return
			}
			assets.Rules = rules

			apps.Set(appId, assets)
		})
	}
	wg.Wait()

	return apps.Data()
}
