package planner

import (
	"fmt"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/pkg/signature"
	"github.com/scrapnode/kanthor/pkg/timer"
)

type Applicable struct {
	EndpointMap map[string]entities.Endpoint
	Rules       []entities.EndpointRule
}

func Requests(msg *entities.Message, applicable *Applicable, timer timer.Timer) ([]entities.Request, [][]interface{}) {
	requests := []entities.Request{}
	logs := [][]interface{}{}
	seen := map[string]bool{}

	for _, epr := range applicable.Rules {
		// already evaluated rules of this endpoint, ignore
		if ignore, ok := seen[epr.EpId]; ok && ignore {
			continue
		}

		log := []interface{}{
			"msg_id", msg.Id,
			"ep_id", epr.EpId,
			"epr_id", epr.Id,
			"epr_condition_source", epr.ConditionSource,
			"epr_condition_expression", epr.ConditionExpression,
		}

		check, err := ConditionExpression(&epr)
		if err != nil {
			logs = append(logs, append([]interface{}{"PLANNER.PLAN_REQUEST.RULE.CONDITION_EXRESSION.ERROR", "err", err.Error()}, log...))
			// once we got error of evaludation, ignore request scheduling for this endpoint
			seen[epr.EpId] = true
			continue
		}
		source := ConditionSource(&epr, msg)
		if source == "" {
			logs = append(logs, append([]interface{}{"PLANNER.PLAN_REQUEST.RULE.CONDITION_SOURCE.EMPTY"}, log...))
			// once we got error of evaludation, ignore request scheduling for this endpoint
			seen[epr.EpId] = true
			continue
		}

		matched := check(source)

		// once we got exclusionary rule, ignore all other rules of current endpoint
		if epr.Exclusionary && matched {
			logs = append(logs, append([]interface{}{"PLANNER.PLAN_REQUEST.RULE.EXCLUSIONARY"}, log...))
			seen[epr.EpId] = true
			continue
		}

		if !matched {
			logs = append(logs, append([]interface{}{"PLANNER.PLAN_REQUEST.RULE.NOT_MATCHED"}, log...))
			seen[epr.EpId] = false
			continue
		}

		ep := applicable.EndpointMap[epr.EpId]
		req := Request(msg, &ep, &epr, timer)
		requests = append(requests, req)
	}

	return requests, logs
}

func Request(msg *entities.Message, ep *entities.Endpoint, epr *entities.EndpointRule, timer timer.Timer) entities.Request {
	// construct request
	req := entities.Request{
		MsgId:    msg.Id,
		EpId:     ep.Id,
		Tier:     msg.Tier,
		AppId:    msg.AppId,
		Type:     msg.Type,
		Metadata: entities.Metadata{},
		Headers:  entities.NewHeader(),
		Body:     msg.Body,
		Uri:      ep.Uri,
		Method:   ep.Method,
	}
	// must use merge function otherwise you will edit the original data
	req.Headers.Merge(msg.Headers)
	req.Metadata.Merge(msg.Metadata)
	req.GenId()
	req.SetTS(timer.Now())

	req.Metadata.Set(entities.MetaEprId, epr.Id)

	req.Headers.Set(entities.HeaderIdempotencyKey, req.Id)
	req.Headers.Set(entities.HeaderMsgRef, fmt.Sprintf("%s/%s", msg.AppId, msg.Id))
	req.Headers.Set(entities.HeaderReqTs, fmt.Sprintf("%d", req.Timestamp))

	// signature
	sign := fmt.Sprintf("%s.%d.%s", msg.Id, req.Timestamp, string(msg.Body))
	signed := signature.Sign(sign, ep.SecretKey)
	req.Headers.Set(entities.HeaderReqSig, fmt.Sprintf("v1=%s", signed))

	return req
}