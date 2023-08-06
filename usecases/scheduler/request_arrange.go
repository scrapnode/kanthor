package scheduler

import (
	"context"
	"errors"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/usecases/scheduler/repos"
	"github.com/sourcegraph/conc/pool"
	"regexp"
	"strings"
	"time"
)

func (uc *request) Arrange(ctx context.Context, req *RequestArrangeReq) (*RequestArrangeRes, error) {
	key := utils.Key("APP_WITH_ENDPOINTS", req.Message.AppId)
	app, err := cache.Warp(uc.cache, ctx, key, time.Hour, func() (*repos.ApplicationWithEndpointsAndRules, error) {
		return uc.repos.Application().ListEndpointsWithRules(ctx, req.Message.AppId)
	})
	if err != nil {
		return nil, err
	}

	res := &RequestArrangeRes{
		Entities:    []structure.BulkRes[entities.Request]{},
		FailKeys:    []string{},
		SuccessKeys: []string{},
	}

	requests := uc.generateRequestsFromEndpoints(req.Message, app.Endpoints)
	if len(requests) == 0 {
		uc.logger.Warnw("no request was generated", "app_id", req.Message.AppId, "message_id", req.Message.Id)
		return res, nil
	}

	p := pool.New().WithMaxGoroutines(uc.conf.Scheduler.Request.Arrange.Concurrency)
	for _, r := range requests {
		ent := r
		p.Go(func() {
			entKey := utils.Key(
				req.Message.AppId,
				ent.Metadata.Get(entities.MetaEpId),
				ent.Metadata.Get(entities.MetaEprId),
				ent.Metadata.Get(entities.MetaMsgId),
				ent.Id,
			)

			event, err := transformRequest2Event(&ent)
			if err == nil {
				err = uc.publisher.Pub(ctx, event)
			}

			res.Entities = append(res.Entities, structure.BulkRes[entities.Request]{Entity: ent, Error: err})
			if err == nil {
				res.SuccessKeys = append(res.SuccessKeys, entKey)
			} else {
				res.FailKeys = append(res.FailKeys, entKey)
				uc.logger.Errorw(err.Error(), "app_id", req.Message.AppId, "key", entKey)
			}
		})
	}
	p.Wait()

	return res, nil
}

func (uc *request) generateRequestsFromEndpoints(
	msg entities.Message,
	endpoints []repos.EndpointWithRules,
) []entities.Request {
	requests := []entities.Request{}

	for _, ep := range endpoints {
		// with this for loop, we enforce endpoint must have at least one rule to construct scheduled request
		for _, epr := range ep.Rules {
			subLogger := uc.logger.With(
				"epr_id", epr.Id,
				"epr_condition_source", epr.ConditionSource,
				"epr_condition_expression", epr.ConditionExpression,
			)
			source := resolveConditionSource(epr, msg)
			if source == "" {
				subLogger.Errorw("arrange: unable to get data source to compare rule")
				continue
			}

			express, err := resolveConditionExpression(epr)
			if err != nil {
				subLogger.Errorw("arrange: unable resolve rule expression", "error", err.Error())
				continue
			}

			matched := express(source)
			// once we got exclusionary rule, ignore the rest
			if epr.Exclusionary && matched {
				break
			}
			// otherwise continue express another condition
			if !matched {
				continue
			}

			ent := entities.Request{
				Tier:     msg.Tier,
				AppId:    msg.AppId,
				Type:     msg.Type,
				Uri:      ep.Uri,
				Method:   ep.Method,
				Headers:  msg.Headers,
				Body:     msg.Body,
				Metadata: entities.Metadata{},
			}
			ent.GenId()
			ent.SetTS(uc.timer.Now(), uc.conf.Bucket.Layout)
			ent.Metadata.Merge(msg.Metadata)
			ent.Metadata.Set(entities.MetaEpId, ep.Id)
			ent.Metadata.Set(entities.MetaEprId, epr.Id)
			ent.Metadata.Set(entities.MetaReqId, ent.Id)

			requests = append(requests, ent)
		}
	}

	return requests
}

func resolveConditionSource(rule entities.EndpointRule, msg entities.Message) string {
	if rule.ConditionSource == "app_id" {
		return msg.AppId
	}

	if rule.ConditionSource == "type" {
		return msg.Type
	}

	if rule.ConditionSource == "body" {
		return string(msg.Body)
	}

	if strings.HasPrefix(rule.ConditionSource, "metadata") {
		kv := strings.Split(rule.ConditionSource, ".")
		if meta, ok := msg.Metadata[strings.Join(kv[1:], ".")]; ok {
			return meta
		}
	}

	return ""
}

func resolveConditionExpression(rule entities.EndpointRule) (func(source string) bool, error) {
	expression := strings.Split(rule.ConditionExpression, "::")
	if len(expression) != 2 {
		return nil, errors.New("invalid rule expression")
	}

	if expression[0] == "regex" {
		r, err := regexp.Compile(expression[1])
		if err != nil {
			return nil, err
		}
		return func(source string) bool { return r.MatchString(source) }, nil
	}

	if expression[0] == "equal" {
		return func(source string) bool { return expression[1] == source }, nil
	}

	return nil, errors.New("unknown rule expression")
}

func transformRequest2Event(req *entities.Request) (*streaming.Event, error) {
	data, err := req.Marshal()
	if err != nil {
		return nil, err
	}

	event := &streaming.Event{
		AppId:    req.AppId,
		Type:     req.Type,
		Id:       req.Id,
		Data:     data,
		Metadata: map[string]string{},
	}
	event.Subject = streaming.Subject(
		streaming.Namespace,
		req.Tier,
		streaming.TopicReq,
		event.AppId,
		event.Type,
	)

	return event, nil
}
