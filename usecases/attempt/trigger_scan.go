package attempt

import (
	"context"
	"fmt"
	"time"

	"github.com/samber/lo"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/internal/planner"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/pkg/validator"
	"github.com/scrapnode/kanthor/usecases/attempt/repos"
)

type TriggerScanReq struct {
	AppId string
	From  time.Time
	To    time.Time
}

func (req *TriggerScanReq) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("app_id", req.AppId, entities.IdNsApp),
	)
}

type TriggerScanRes struct {
	ToScheduleMsgIds []string
	ToAttemptReqs    []repos.Req
}

var (
	reqKeyFmt = "%s/%s/request"
	resKeyFmt = "%s/%s/response"
)

func (uc *trigger) Scan(ctx context.Context, req *TriggerScanReq) (*TriggerScanRes, error) {
	// @TODO: use internal cursor

	messages, err := uc.repos.Message().Scan(ctx, req.AppId, req.From, req.To)
	if err != nil {
		return nil, err
	}
	msgIds := lo.Keys(messages)

	requests, err := uc.repos.Request().Scan(ctx, req.AppId, msgIds, req.From, req.To)
	if err != nil {
		return nil, err
	}

	responses, err := uc.repos.Response().Scan(ctx, req.AppId, msgIds, req.From, req.To)
	if err != nil {
		return nil, err
	}

	app, err := uc.applicable(ctx, req.AppId)
	if err != nil {
		return nil, err
	}

	res := &TriggerScanRes{ToScheduleMsgIds: []string{}, ToAttemptReqs: []repos.Req{}}
	status := uc.hash(requests, responses)
	for _, message := range messages {
		for _, ep := range app.EndpointMap {
			reqKey := fmt.Sprintf(reqKeyFmt, message.Id, ep.Id)
			reqId, hasReq := status[reqKey]
			if !hasReq {
				// no request -> must schedule message again -> don't create any attempt
				res.ToScheduleMsgIds = append(res.ToScheduleMsgIds, message.Id)
				continue
			}

			resKey := fmt.Sprintf(resKeyFmt, message.Id, ep.Id)
			_, hasRes := status[resKey]
			if !hasRes {
				// has request + no success response -> create an attempt
				res.ToAttemptReqs = append(res.ToAttemptReqs, requests[reqId])
				continue
			}

			// has success response, ignore
		}
	}

	return res, nil
}

func (uc *trigger) applicable(ctx context.Context, appId string) (*planner.Applicable, error) {
	key := utils.Key("scheduler", appId)
	return cache.Warp(uc.cache, ctx, key, time.Hour, func() (*planner.Applicable, error) {
		uc.metrics.Count(ctx, "cache_miss_total", 1)

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
}

func (uc *trigger) hash(requests map[string]repos.Req, responses map[string]repos.Res) map[string]string {
	returning := map[string]string{}

	for _, request := range requests {
		// for checking whether we have scheduled a request for an endpoint or not
		// if no request was scheduled, we should schedule it instead of create an attempt
		key := fmt.Sprintf(reqKeyFmt, request.MsgId, request.EpId)
		returning[key] = request.Id
	}

	for _, response := range responses {
		key := fmt.Sprintf(resKeyFmt, response.MsgId, response.EpId)

		// we already recognized that the endpoint had success status, don't need to check any more
		if _, has := returning[key]; has {
			continue
		}

		// status is ok, saved the success response id
		if entities.Is2xx(response.Status) {
			returning[key] = response.Id
		}
	}

	return returning
}
