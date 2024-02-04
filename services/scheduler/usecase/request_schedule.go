package usecase

import (
	"context"
	"fmt"

	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/internal/routing"
	"github.com/scrapnode/kanthor/internal/transformation"
	"github.com/scrapnode/kanthor/pkg/safe"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/pkg/validator"
	"github.com/scrapnode/kanthor/telemetry"
	"go.opentelemetry.io/otel/attribute"
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
	return validator.Validate(
		validator.DefaultConfig,
		validator.MapRequired("messages", in.Messages),
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
	spanName := "usecase.request.schedule"
	spanner := ctx.Value(telemetry.CtxSpanner).(*telemetry.Spanner)

	ok := &safe.Map[[]string]{}
	ko := &safe.Map[error]{}

	// we have to store a ref map so if we got any error,
	// we can report back to the caller that a key has a error and should be retry
	refIds := map[string]string{}
	for refId, msg := range in.Messages {
		spanner.StartWithRefId(spanName, refId)
		refIds[msg.Id] = refId
	}
	defer spanner.End(spanName)

	errc := make(chan error, 1)
	defer close(errc)

	go func() {
		requests, err := uc.arrange(context.WithValue(ctx, telemetry.CtxSpanner, spanner.Clone()), in)
		if err != nil {
			errc <- err
			return
		}
		if len(requests) == 0 {
			errc <- nil
			return
		}

		publishSpanner := spanner.Clone()
		publishSpanner.Start("usecase.request.publish")

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
			eventRef := refIds[msgId]

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

		publishSpanner.End("usecase.request.publish")
		errc <- nil
	}()

	select {
	case err := <-errc:
		return &RequestScheduleOut{Success: ok.Keys(), Error: ko.Data()}, err
	case <-ctx.Done():
		// context deadline exceeded, should set that error to remain messages
		for _, msg := range in.Messages {
			eventRef := refIds[msg.Id]

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

func (uc *request) arrange(ctx context.Context, in *RequestScheduleIn) (map[string]*entities.Request, error) {
	spanName := "usecase.request.arrange"
	spanner := ctx.Value(telemetry.CtxSpanner).(*telemetry.Spanner)

	appIds := make([]string, 0)
	for refId := range in.Messages {
		spanner.StartWithRefId(spanName, refId)
		appIds = append(appIds, in.Messages[refId].AppId)
	}
	defer spanner.End(spanName)

	appSpanner := spanner.Clone()
	appSpanner.Start("repositories.db.application.get_routes")
	routes, err := uc.repositories.Database().Application().GetRoutes(ctx, appIds)
	if err != nil {
		appSpanner.End("repositories.db.application.get_routes")
		return nil, err
	}
	appSpanner.End("repositories.db.application.get_routes", attribute.Int("app.route_count", len(routes)))

	planningSpanner := spanner.Clone()
	returning := map[string]*entities.Request{}
	for refId, msg := range in.Messages {
		spanNamePlanning := "usecase.request.arrange.planning"
		if items, has := routes[msg.AppId]; has {
			planningSpanner.StartWithRefId("usecase.request.arrange.planning", refId)

			requests, traces := routing.PlanRequests(uc.infra.Timer, msg, items)
			if len(traces) > 0 {
				for _, trace := range traces {
					uc.logger.Warnw(trace[0].(string), trace[1:]...)
				}
			}

			for id, request := range requests {
				returning[id] = request
			}

			planningSpanner.End(spanNamePlanning, attribute.Int("request_count", len(requests)))
		}
	}

	if len(returning) == 0 {
		uc.logger.Warnw("SCHEDULER.USECASE.REQUEST.SCHEDULE.NO_REQUEST", "apps", appIds)
	}

	return returning, nil
}
