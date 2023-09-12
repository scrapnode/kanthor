package scheduler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
	"github.com/scrapnode/kanthor/infrastructure/cache"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/usecases/transformation"
	"github.com/sourcegraph/conc/pool"
)

type applicable struct {
	Endpoints map[string]entities.Endpoint
	Rules     []entities.EndpointRule
}

func (uc *request) Arrange(ctx context.Context, req *RequestArrangeReq) (*RequestArrangeRes, error) {
	key := utils.Key("scheduler", req.Message.AppId)
	// @TODO: find a way to notify attempt that message is not able to schedule
	app, err := cache.Warp(uc.cache, ctx, key, time.Hour, func() (*applicable, error) {
		uc.metrics.Count(ctx, "cache_miss_total", 1)

		endpoints, err := uc.repos.Endpoint().List(ctx, req.Message.AppId)
		if err != nil {
			return nil, err
		}

		returning := &applicable{Endpoints: map[string]entities.Endpoint{}}
		ids := []string{}
		for _, endpoint := range endpoints {
			ids = append(ids, endpoint.Id)
			returning.Endpoints[endpoint.Id] = endpoint
		}

		rules, err := uc.repos.EndpointRule().List(ctx, ids)
		if err != nil {
			return nil, err
		}
		returning.Rules = rules

		return returning, nil
	})
	if err != nil {
		return nil, err
	}

	res := &RequestArrangeRes{
		Entities:    []structure.BulkRes[entities.Request]{},
		FailKeys:    []string{},
		SuccessKeys: []string{},
	}

	requests := uc.generateRequestsFromEndpoints(req.Message, app)
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
				req.Message.AttId,
				ent.Id,
			)

			event, err := transformation.EventFromRequest(&ent)
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
	app *applicable,
) []entities.Request {
	requests := []entities.Request{}
	seen := map[string]bool{}

	for _, epr := range app.Rules {
		// already evaluated rules of this endpoint, ignore
		if ignore, ok := seen[epr.EndpointId]; ok && ignore {
			continue
		}

		slogger := uc.logger.With(
			"epr_id", epr.Id,
			"epr_condition_source", epr.ConditionSource,
			"epr_condition_expression", epr.ConditionExpression,
		)

		source := resolveConditionSource(epr, msg)
		if source == "" {
			slogger.Errorw("arrange: unable to get data source to compare rule")
			continue
		}

		express, err := resolveConditionExpression(epr)
		if err != nil {
			slogger.Errorw("arrange: unable resolve rule expression", "error", err.Error())
			continue
		}

		matched := express(source)
		// once we got exclusionary rule, ignore all other rules of current endpoint
		if epr.Exclusionary && matched {
			seen[epr.EndpointId] = true
			slogger.Warn("arrange: matched exclusionary rule")
			continue
		}

		if !matched {
			seen[epr.EndpointId] = false
			slogger.Debugw("arrange: rule is not matched")
			continue
		}

		// construct request
		ep := app.Endpoints[epr.EndpointId]
		req := entities.Request{
			AttId:    msg.AttId,
			Tier:     msg.Tier,
			AppId:    msg.AppId,
			Type:     msg.Type,
			Metadata: entities.Metadata{},
			Headers:  entities.Header{Header: http.Header{}},
			Body:     msg.Body,
			Uri:      ep.Uri,
			Method:   ep.Method,
		}
		// must use merge function otherwise you will edit the original data
		req.Headers.Merge(msg.Headers)
		req.Metadata.Merge(msg.Metadata)
		req.GenId()
		req.SetTS(uc.timer.Now(), uc.conf.Bucket.Layout)

		req.Headers.Set("idempotency-key", req.AttId)
		req.Headers.Set("kanthor-msg-id", msg.Id)
		req.Headers.Set("kanthor-req-ts", fmt.Sprintf("%d", req.Timestamp))

		sign := fmt.Sprintf("%s.%d.%s", msg.Id, req.Timestamp, string(msg.Body))
		signed := uc.signature.Sign(sign, ep.SecretKey)
		req.Headers.Set("kanthor-req-signature", signed)

		req.Metadata.Set(entities.MetaEpId, ep.Id)
		req.Metadata.Set(entities.MetaEprId, epr.Id)

		requests = append(requests, req)
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
