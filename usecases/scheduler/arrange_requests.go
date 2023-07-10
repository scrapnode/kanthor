package scheduler

import (
	"context"
	"errors"
	"fmt"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/repositories"
	"github.com/scrapnode/kanthor/domain/structure"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/infrastructure/monitoring/metric"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/sourcegraph/conc/pool"
	"regexp"
	"strings"
	"time"
)

func (usecase *scheduler) ArrangeRequests(ctx context.Context, req *ArrangeRequestsReq) (*ArrangeRequestsRes, error) {
	cacheKey := cache.Key("APP_WITH_ENDPOINTS", req.Message.AppId)
	app, err := cache.Warp(usecase.cache, cacheKey, time.Hour, func() (*repositories.ApplicationWithEndpointsAndRules, error) {
		usecase.meter.Count("cache_miss_total", 1, metric.Label("source", "scheduler_arrange_requests"))
		return usecase.repos.Application().ListEndpointsWithRules(ctx, req.Message.AppId)
	})
	if err != nil {
		usecase.logger.Errorw(err.Error(), "app_id", req.Message.AppId)
		return nil, fmt.Errorf("unable to find application [%s]", req.Message.AppId)
	}

	res := &ArrangeRequestsRes{
		Entities:    []structure.BulkRes[entities.Request]{},
		FailKeys:    []string{},
		SuccessKeys: []string{},
	}

	requests := usecase.generateRequestsFromEndpoints(app.Endpoints, req.Message)
	if len(requests) == 0 {
		usecase.logger.Warnw("no request was generated", "app_id", req.Message.AppId, "message_id", req.Message.Id)
		return res, nil
	}

	p := pool.New().WithMaxGoroutines(usecase.conf.Scheduler.ArrangeRequests.Concurrency)
	for _, r := range requests {
		request := r
		p.Go(func() {
			key := utils.Key(req.Message.Id, request.EndpointId, request.Metadata[entities.MetaRuleId], request.Id)

			event, err := transformRequest2Event(&request)
			if err == nil {
				err = usecase.publisher.Pub(ctx, event)
			}

			res.Entities = append(res.Entities, structure.BulkRes[entities.Request]{Entity: request, Error: err})
			if err == nil {
				res.SuccessKeys = append(res.SuccessKeys, key)
			} else {
				res.FailKeys = append(res.FailKeys, key)
				usecase.logger.Errorw(err.Error(), "app_id", req.Message.AppId, "key", key)
			}
		})
	}
	p.Wait()

	return res, nil
}

func (usecase *scheduler) generateRequestsFromEndpoints(endpoints []repositories.EndpointWithRules, msg entities.Message) []entities.Request {
	var requests []entities.Request

	for _, endpoint := range endpoints {
		// with this for loop, we enforce endpoint must have at least one rule to construct scheduled request
		for _, rule := range endpoint.Rules {
			subLogger := usecase.logger.With(
				"rule_id", rule.Id,
				"rule_condition_source", rule.ConditionSource,
				"rule_condition_expression", rule.ConditionExpression,
			)
			source := resolveConditionSource(rule, msg)
			if source == "" {
				subLogger.Errorw("arrange: unable to get data source to compare rule")
				continue
			}

			express, err := resolveConditionExpression(rule)
			if err != nil {
				subLogger.Errorw("arrange: unable resolve rule expression", "error", err.Error())
				continue
			}

			expression := strings.Split(rule.ConditionExpression, "::")
			if len(expression) != 2 {
				subLogger.Errorw("arrange: invalid rule")
				continue
			}

			matched := express(source)
			// once we got exclusionary rule, ignore the rest
			if rule.Exclusionary && matched {
				break
			}
			// otherwise continue express another condition
			if !matched {
				continue
			}

			request := entities.Request{
				Tier:       msg.Tier,
				AppId:      msg.AppId,
				Type:       msg.Type,
				EndpointId: endpoint.Id,
				Uri:        endpoint.Uri,
				Method:     endpoint.Method,
				Headers:    msg.Headers,
				Body:       msg.Body,
				Metadata: map[string]string{
					entities.MetaMsgId:  msg.Id,
					entities.MetaRuleId: rule.Id,
				},
			}
			request.GenId()
			request.SetTS(usecase.timer.Now(), usecase.conf.Bucket.Layout)

			requests = append(requests, request)
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
		Data:     data,
		Metadata: map[string]string{},
	}
	event.GenId()
	event.Subject = streaming.Subject(
		streaming.Namespace,
		req.Tier,
		streaming.TopicReq,
		event.AppId,
		event.Type,
	)

	return event, nil
}
