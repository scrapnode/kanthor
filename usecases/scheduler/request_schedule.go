package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/samber/lo"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/internal/planner"
	"github.com/scrapnode/kanthor/pkg/safe"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/pkg/validator"
	"github.com/scrapnode/kanthor/usecases/transformation"
	"github.com/sourcegraph/conc"
	"github.com/sourcegraph/conc/pool"
)

type RequestScheduleReq struct {
	Timeout   int64
	RateLimit int
	Messages  []entities.Message
}

func ValidateRequestScheduleReqMessage(prefix string, message *entities.Message) error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith(prefix+".id", message.Id, entities.IdNsMsg),
		validator.StringRequired(prefix+".tier", message.Tier),
		validator.StringStartsWith(prefix+".app_id", message.AppId, entities.IdNsApp),
		validator.StringRequired(prefix+".type", message.Type),
		validator.MapNotNil[string, string](prefix+".metadata", message.Metadata),
		validator.SliceRequired(prefix+".body", message.Body),
	)
}

func (req *RequestScheduleReq) Validate() error {
	err := validator.Validate(
		validator.DefaultConfig,
		validator.SliceRequired("messages", req.Messages),
		validator.NumberGreaterThan("timeout", req.Timeout, 1000),
		validator.NumberGreaterThan("rate_limit", req.Timeout, 1),
	)
	if err != nil {
		return err
	}

	return validator.Validate(
		validator.DefaultConfig,
		validator.Array(req.Messages, func(i int, item *entities.Message) error {
			prefix := fmt.Sprintf("messages[%d]", i)
			return ValidateRequestScheduleReqMessage(prefix, item)
		}),
	)
}

type RequestScheduleRes struct {
	Success []string
	Error   map[string]error
}

func (uc *request) Schedule(ctx context.Context, req *RequestScheduleReq) (*RequestScheduleRes, error) {
	requests := uc.arrange(ctx, req.Messages)
	if len(requests) == 0 {
		return &RequestScheduleRes{Success: []string{}, Error: map[string]error{}}, nil
	}

	ok := &safe.Map[string]{}
	ko := &safe.Map[error]{}

	// timeout duration will be scaled based on how many requests you have
	duration := time.Duration(req.Timeout * int64(len(requests)+1))
	timeout, cancel := context.WithTimeout(ctx, time.Millisecond*duration)
	defer cancel()

	// we must limit how many publish action could be performed at the same time
	// otherwise our system will be lag
	p := pool.New().WithMaxGoroutines(req.RateLimit)
	for _, r := range requests {
		request := r
		p.Go(func() {
			event, err := transformation.EventFromRequest(&request)
			if err != nil {
				// un-recoverable error
				uc.logger.Errorw("could not transform request to event", "request", request.String())
				return
			}

			if err := uc.publisher.Pub(ctx, event); err != nil {
				ko.Set(request.MsgId, err)
				return
			}

			key := utils.Key(request.AppId, request.MsgId, request.EpId, request.Id)
			ok.Set(key, request.MsgId)
		})
	}

	c := make(chan bool)
	defer close(c)

	go func() {
		p.Wait()
		c <- true
	}()

	select {
	case <-c:
		return &RequestScheduleRes{Success: ok.Keys(), Error: ko.Data()}, nil
	case <-timeout.Done():
		// context deadline exceeded, should set that error to remain messages
		for _, message := range req.Messages {
			if _, success := ok.Get(message.Id); success {
				// already success, should not retry it
				continue
			}

			// no error, should add context deadline error
			if _, has := ko.Get(message.Id); !has {
				ko.Set(message.Id, ctx.Err())
			}
		}
		return &RequestScheduleRes{Success: ok.Keys(), Error: ko.Data()}, nil
	}
}

func (uc *request) arrange(ctx context.Context, messages []entities.Message) []entities.Request {
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
			continue
		}

		reqs, logs := planner.Requests(&message, &app, uc.infra.Timer)
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

func (uc *request) applicables(ctx context.Context, appIds []string) map[string]planner.Applicable {
	apps := &safe.Map[planner.Applicable]{}

	var wg conc.WaitGroup
	for _, id := range appIds {
		appId := id
		wg.Go(func() {
			key := utils.Key("scheduler", appId)
			app, err := cache.Warp(uc.infra.Cache, ctx, key, time.Hour, func() (*planner.Applicable, error) {
				endpoints, err := uc.repos.Endpoint().List(ctx, appId)
				if err != nil {
					return nil, err
				}
				returning := &planner.Applicable{EndpointMap: map[string]entities.Endpoint{}}
				for _, ep := range endpoints {
					returning.EndpointMap[ep.Id] = ep
				}

				rules, err := uc.repos.Endpoint().Rules(ctx, appId)
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
			apps.Set(appId, *app)
		})
	}

	return apps.Data()
}
