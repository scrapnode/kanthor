package scheduler

import (
	"context"
	"errors"
	"github.com/scrapnode/kanthor/domain/constants"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/repositories"
	"github.com/scrapnode/kanthor/domain/structure"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/infrastructure/utils"
	"github.com/sourcegraph/conc/pool"
	"regexp"
	"strings"
)

func (service *scheduler) ArrangeRequests(ctx context.Context, req *ArrangeRequestsReq) (*ArrangeRequestsRes, error) {
	res := &ArrangeRequestsRes{Entities: []structure.BulkRes[entities.Request]{}, FailKeys: []string{}, SuccessKeys: []string{}}

	app, err := service.repos.Application().Get(ctx, req.Message.AppId)
	if err != nil {
		return nil, err
	}
	ws, err := service.repos.Workspace().Get(ctx, app.WorkspaceId)
	if err != nil {
		return nil, err
	}

	// rules of endpoint are well sorted slice like this
	// IMPORTANT: the order is important
	// rule.priority - rule.exclusionary
	// 			  15 - TRUE
	// 			  15 - FALSE
	// 			  9  - FALSE
	//		  	  70 - TRUE
	// 			  70 - FALSE
	//			  0  - FALSE
	endpoints, err := service.repos.Endpoint().ListWithRules(ctx, app.Id)
	if err != nil {
		return nil, err
	}

	requests := service.generateRequestsFromEndpoints(endpoints, req.Message)
	if len(requests) == 0 {
		service.logger.Warnw("no request was generated", "message_id", req.Message.Id)
	}

	// @TODO: remove hardcode of max goroutines here
	p := pool.New().WithMaxGoroutines(10)
	for _, r := range requests {
		request := r
		p.Go(func() {
			key := utils.Key(req.Message.Id, request.Metadata["endpoint_id"], request.Metadata["rule_id"], request.Id)

			event, err := transformRequest2Event(&request)
			if err == nil {
				subject := streaming.Subject(
					constants.TopicNamespace,
					ws.Tier.Name,
					constants.TopicRequest,
					event.AppId,
					event.Type,
				)
				err = service.publisher.Pub(ctx, subject, event)
			}

			res.Entities = append(res.Entities, structure.BulkRes[entities.Request]{Entity: request, Error: err})
			if err == nil {
				res.SuccessKeys = append(res.SuccessKeys, key)
			} else {
				res.FailKeys = append(res.FailKeys, key)
			}
		})
	}
	p.Wait()

	return res, nil
}

func (service *scheduler) generateRequestsFromEndpoints(endpoints []repositories.EndpointWithRules, msg *entities.Message) []entities.Request {
	var requests []entities.Request

	for _, endpoint := range endpoints {
		for _, rule := range endpoint.Rules {
			subLogger := service.logger.With(
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
				AppId:  msg.AppId,
				Type:   msg.Type,
				Uri:    endpoint.Uri,
				Method: endpoint.Method,
				Body:   msg.Body,
				Metadata: map[string]string{
					"endpoint_id": endpoint.Id,
					"rule_id":     rule.Id,
				},
				Status: entities.StatusScheduled,
			}
			request.GenId()
			request.SetTS(service.timer.Now(), service.conf.Bucket.Layout)

			requests = append(requests, request)
		}
	}

	return requests
}

func resolveConditionSource(rule entities.EndpointRule, msg *entities.Message) string {
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

	return event, nil
}
