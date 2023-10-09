package scheduler

import (
	"context"
	"time"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/internal/planner"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/pkg/validator"
)

type RequestArrangeReq struct {
	Message *entities.Message
}

func (req *RequestArrangeReq) Validate() error {
	err := validator.Validate(validator.DefaultConfig, validator.PointerNotNil[entities.Message]("message", req.Message))
	if err != nil {
		return err
	}
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("message.id", req.Message.Id, entities.IdNsMsg),
		validator.StringRequired("message.tier", req.Message.Tier),
		validator.StringStartsWith("message.app_id", req.Message.AppId, entities.IdNsApp),
		validator.StringRequired("message.type", req.Message.Type),
		validator.MapNotNil[string, string]("message.metadata", req.Message.Metadata),
		validator.SliceRequired("message.body", req.Message.Body),
	)
}

type RequestArrangeRes struct {
	Requests []entities.Request
}

func (uc *request) Arrange(ctx context.Context, req *RequestArrangeReq) (*RequestArrangeRes, error) {
	key := utils.Key("scheduler", req.Message.AppId)
	app, err := cache.Warp(uc.cache, ctx, key, time.Hour, func() (*planner.Applicable, error) {
		uc.metrics.Count(ctx, "cache_miss_total", 1)

		endpoints, err := uc.repos.Endpoint().List(ctx, req.Message.AppId)
		if err != nil {
			return nil, err
		}
		returning := &planner.Applicable{EndpointMap: map[string]entities.Endpoint{}}
		for _, ep := range endpoints {
			returning.EndpointMap[ep.Id] = ep
		}

		rules, err := uc.repos.Endpoint().Rules(ctx, req.Message.AppId)
		if err != nil {
			return nil, err
		}
		returning.Rules = rules

		return returning, nil
	})
	if err != nil {
		return nil, err
	}

	requests, logs := planner.Requests(req.Message, app, uc.timer)
	if len(logs) > 0 {
		for _, l := range logs {
			uc.logger.Warnw(l[0].(string), l[1:]...)
		}
	}

	res := &RequestArrangeRes{Requests: requests}
	return res, nil
}
