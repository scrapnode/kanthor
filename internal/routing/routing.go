package routing

import (
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/safe"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/sourcegraph/conc"
)

type Route struct {
	Endpoint *entities.Endpoint
	Rules    []entities.EndpointRule
}

func PlanRequests(
	timer timer.Timer,
	msg *entities.Message,
	routes []Route,
) (map[string]*entities.Request, [][]any) {
	requests := &safe.Map[*entities.Request]{}
	traces := &safe.Slice[[]any]{}

	var wg conc.WaitGroup
	for i := range routes {
		route := routes[i]
		wg.Go(func() {
			request, trace := PlanRequest(timer, msg, &route)
			if request != nil {
				requests.Set(request.Id, request)
			}
			if len(trace) > 0 {
				traces.Append(trace)
			}
		})
	}
	wg.Wait()

	return requests.Data(), traces.Data()
}

func PlanRequest(
	timer timer.Timer,
	msg *entities.Message,
	route *Route,
) (*entities.Request, []any) {
	for i := range route.Rules {
		trace := []any{
			"msg_id", msg.Id,
			"ep_id", route.Endpoint.Id,
			"epr_id", route.Rules[i].Id,
			"epr_cs", route.Rules[i].ConditionSource,
			"epr_ce", route.Rules[i].ConditionExpression,
		}

		check, err := ConditionExpression(&route.Rules[i])
		if err != nil {
			return nil, append([]any{"ROUTING.PLAN.RULE.CE.ERROR", "error", err.Error()}, trace...)
		}

		source := ConditionSource(&route.Rules[i], msg)
		if source == "" {
			return nil, append([]any{"ROUTING.PLAN.RULE.CS.EMPTY.ERROR"}, trace...)
		}

		matched := check(source)
		if route.Rules[i].Exclusionary && matched {
			return nil, append([]any{"ROUTING.PLAN.RULE.EXCLUSIONARY.ERROR"}, trace...)
		}

		if matched {
			return NewRequest(timer, msg, route.Endpoint, &route.Rules[i]), nil
		}
	}

	return nil, []any{"ROUTING.PLAN.NOT_MATCH.ERROR", "msg_id", msg.Id, "ep_id", route.Endpoint.Id}
}