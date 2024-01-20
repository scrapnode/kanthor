package routing_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/internal/routing"
	"github.com/scrapnode/kanthor/internal/tester"
	"github.com/scrapnode/kanthor/mocks"
	"github.com/stretchr/testify/assert"
)

func TestPlanRequests(t *testing.T) {
	now := time.Now().UTC()
	timer := mocks.NewTimer(t)
	timer.On("Now").Return(now)

	app := tester.Application(timer)
	ep := tester.EndpointOfApp(timer, app)

	t.Run("success", func(st *testing.T) {
		route := &routing.Route{Endpoint: ep, Rules: make([]entities.EndpointRule, 0)}
		route.Rules = append(route.Rules, entities.EndpointRule{
			ConditionSource:     routing.ConditionSourceType,
			ConditionExpression: fmt.Sprintf("%s%s", routing.ConditionExpressionAny, routing.ConditionExpressionDivider),
		})

		msg := tester.MessageOfApp(timer, app)
		req, trace := routing.PlanRequest(timer, msg, route)
		assert.True(st, len(trace) == 0)
		assert.NotNil(st, req)
	})

	t.Run("error of condition expression", func(st *testing.T) {
		route := &routing.Route{Endpoint: ep, Rules: make([]entities.EndpointRule, 0)}
		route.Rules = append(route.Rules, entities.EndpointRule{
			ConditionExpression: string(routing.ConditionExpressionDivider[0]),
		})

		msg := tester.MessageOfApp(timer, app)

		_, trace := routing.PlanRequest(timer, msg, route)
		assert.True(st, len(trace) > 0)
		assert.Equal(st, trace[0], "ERROR.ROUTING.PLAN.RULE.CE")
	})

	t.Run("error of condition source empty", func(st *testing.T) {
		route := &routing.Route{Endpoint: ep, Rules: make([]entities.EndpointRule, 0)}
		route.Rules = append(route.Rules, entities.EndpointRule{
			ConditionSource:     routing.ConditionSourceType,
			ConditionExpression: fmt.Sprintf("%s%s", routing.ConditionExpressionAny, routing.ConditionExpressionDivider),
		})

		msg := tester.MessageOfApp(timer, app)
		msg.Type = ""
		_, trace := routing.PlanRequest(timer, msg, route)
		assert.True(st, len(trace) > 0)
		assert.Equal(st, trace[0], "ERROR.ROUTING.PLAN.RULE.CS.EMPTY")
	})

	t.Run("error of exclustionary check", func(st *testing.T) {
		route := &routing.Route{Endpoint: ep, Rules: make([]entities.EndpointRule, 0)}
		route.Rules = append(route.Rules, entities.EndpointRule{
			ConditionSource:     routing.ConditionSourceType,
			ConditionExpression: fmt.Sprintf("%s%s", routing.ConditionExpressionAny, routing.ConditionExpressionDivider),
			Exclusionary:        true,
		})

		msg := tester.MessageOfApp(timer, app)
		_, trace := routing.PlanRequest(timer, msg, route)
		assert.True(st, len(trace) > 0)
		assert.Equal(st, trace[0], "ROUTING.PLAN.RULE.EXCLUSIONARY")
	})

	t.Run("error of not matched any rule", func(st *testing.T) {
		route := &routing.Route{Endpoint: ep, Rules: make([]entities.EndpointRule, 0)}
		route.Rules = append(
			route.Rules,
			entities.EndpointRule{
				ConditionSource:     routing.ConditionSourceType,
				ConditionExpression: fmt.Sprintf("%s%sunable", routing.ConditionExpressionEqual, routing.ConditionExpressionDivider),
				Exclusionary:        false,
			},
			entities.EndpointRule{
				ConditionSource:     routing.ConditionSourceAppId,
				ConditionExpression: fmt.Sprintf("%s%sunable", routing.ConditionExpressionEqual, routing.ConditionExpressionDivider),
				Exclusionary:        false,
			})

		msg := tester.MessageOfApp(timer, app)
		_, trace := routing.PlanRequest(timer, msg, route)
		assert.True(st, len(trace) > 0)
		assert.Equal(st, trace[0], "ERROR.ROUTING.PLAN.NOT_MATCH")
	})
}

func TestPlanRequest(t *testing.T) {
	now := time.Now().UTC()
	timer := mocks.NewTimer(t)
	timer.On("Now").Return(now)

	app := tester.Application(timer)

	t.Run("success", func(st *testing.T) {
		routes := []routing.Route{
			// match this endpoint
			{
				Endpoint: tester.EndpointOfApp(timer, app),
				Rules: []entities.EndpointRule{
					{
						ConditionSource:     routing.ConditionSourceType,
						ConditionExpression: fmt.Sprintf("%s%s", routing.ConditionExpressionAny, routing.ConditionExpressionDivider),
					},
				},
			},
			// match this endpoint
			{
				Endpoint: tester.EndpointOfApp(timer, app),
				Rules: []entities.EndpointRule{
					{
						ConditionSource:     routing.ConditionSourceType,
						ConditionExpression: fmt.Sprintf("%s%s", routing.ConditionExpressionAny, routing.ConditionExpressionDivider),
					},
				},
			},
			// but NOT this one
			{
				Endpoint: tester.EndpointOfApp(timer, app),
				Rules: []entities.EndpointRule{
					{
						ConditionSource:     routing.ConditionSourceType,
						ConditionExpression: fmt.Sprintf("%s%s", routing.ConditionExpressionAny, routing.ConditionExpressionDivider),
						Exclusionary:        true,
					},
				},
			},
		}

		msg := tester.MessageOfApp(timer, app)
		reqs, traces := routing.PlanRequests(timer, msg, routes)
		assert.True(st, len(traces) == 1)
		assert.True(st, len(reqs) == 2)
	})
}
