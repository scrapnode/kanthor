package assessor

import (
	"fmt"

	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/signature"
	"github.com/scrapnode/kanthor/pkg/suid"
	"github.com/scrapnode/kanthor/pkg/timer"
)

type Assets struct {
	EndpointMap map[string]entities.Endpoint
	Rules       []entities.EndpointRule
}

func Requests(msg *entities.Message, assets *Assets, timer timer.Timer) (map[string]*entities.Request, [][]interface{}) {
	requests := map[string]*entities.Request{}
	traces := [][]interface{}{}
	seen := map[string]bool{}

	for _, epr := range assets.Rules {
		// already evaluated rules of this endpoint, ignore
		if ignore, ok := seen[epr.EpId]; ok && ignore {
			continue
		}
		seen[epr.EpId] = false

		trace := []interface{}{
			"msg_id", msg.Id,
			"ep_id", epr.EpId,
			"epr_id", epr.Id,
			"epr_condition_source", epr.ConditionSource,
			"epr_condition_expression", epr.ConditionExpression,
		}

		check, err := ConditionExpression(&epr)
		if err != nil {
			traces = append(traces, append([]interface{}{"ASSESSOR.REQUESTS.RULE.CONDITION_EXRESSION.ERROR", "err", err.Error()}, trace...))
			// once we got error of evaludation, ignore request scheduling for this endpoint
			seen[epr.EpId] = true
			continue
		}
		source := ConditionSource(&epr, msg)
		if source == "" {
			traces = append(traces, append([]interface{}{"ASSESSOR.REQUESTS.RULE.CONDITION_SOURCE.EMPTY"}, trace...))
			// once we got error of evaludation, ignore request scheduling for this endpoint
			seen[epr.EpId] = true
			continue
		}

		matched := check(source)

		// once we got exclusionary rule, ignore all other rules of current endpoint
		if epr.Exclusionary && matched {
			traces = append(traces, append([]interface{}{"ASSESSOR.REQUESTS.RULE.EXCLUSIONARY"}, trace...))
			seen[epr.EpId] = true
			continue
		}

		if matched {
			ep := assets.EndpointMap[epr.EpId]
			req := Request(msg, &ep, &epr, timer)
			requests[req.Id] = req
			seen[epr.EpId] = true
			continue
		}

		traces = append(traces, append([]interface{}{"ASSESSOR.REQUESTS.RULE.NOT_MATCHED"}, trace...))
	}

	return requests, traces
}

func Request(msg *entities.Message, ep *entities.Endpoint, epr *entities.EndpointRule, timer timer.Timer) *entities.Request {
	// construct request
	req := &entities.Request{
		MsgId:    msg.Id,
		EpId:     ep.Id,
		Tier:     msg.Tier,
		AppId:    msg.AppId,
		Type:     msg.Type,
		Metadata: entities.Metadata{},
		Headers:  entities.Header{},
		Body:     msg.Body,
		Uri:      ep.Uri,
		Method:   ep.Method,
	}
	// must use merge function otherwise you will edit the original data
	req.Headers.Merge(msg.Headers)
	req.Metadata.Merge(msg.Metadata)
	req.Id = suid.New(entities.IdNsReq)
	req.SetTS(timer.Now())

	req.Metadata.Set(entities.MetaEprId, epr.Id)

	req.Headers.Set(entities.HeaderIdempotencyKey, msg.Id)

	// https://github.com/standard-webhooks/standard-webhooks/blob/main/spec/standard-webhooks.md
	req.Headers.Set(entities.HeaderWebhookId, msg.Id)
	req.Headers.Set(entities.HeaderWebhookTs, fmt.Sprintf("%d", req.Timestamp))
	sign := fmt.Sprintf("%s.%d.%s", msg.Id, req.Timestamp, msg.Body)
	signed := signature.Sign(sign, ep.SecretKey)
	req.Headers.Set(entities.HeaderWebhookSign, fmt.Sprintf("v1,%s", signed))

	// custom headers
	req.Headers.Set(entities.HeaderWebhookRef, fmt.Sprintf("%s/%s", msg.AppId, ep.Id))

	return req
}
