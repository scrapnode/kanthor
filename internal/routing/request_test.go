package routing_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/internal/routing"
	"github.com/scrapnode/kanthor/internal/tester"
	"github.com/scrapnode/kanthor/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNewRequest(t *testing.T) {
	now := time.Now().UTC()
	timer := mocks.NewTimer(t)
	timer.On("Now").Return(now)

	app := tester.Application(timer)
	msg := tester.MessageOfApp(timer, app)

	random := uuid.NewString()
	msg.Headers.Set("x-random-id", random)
	msg.Metadata.Set("source", "test")

	ep := tester.EndpointOfApp(timer, app)
	rule := tester.RuleOfEndpoint(timer, ep)
	req := routing.NewRequest(timer, msg, ep, rule)

	// constructor
	assert.Equal(t, req.MsgId, msg.Id)
	assert.Equal(t, req.Tier, msg.Tier)
	assert.Equal(t, req.AppId, msg.AppId)
	assert.Equal(t, req.Type, msg.Type)
	assert.Equal(t, req.Body, msg.Body)
	assert.Equal(t, req.EpId, ep.Id)
	assert.Equal(t, req.Uri, ep.Uri)
	assert.Equal(t, req.Method, ep.Method)

	// id
	assert.True(t, strings.HasPrefix(req.Id, entities.IdNsReq))

	// data from message
	assert.Equal(t, req.Metadata.Get("source"), "test")
	assert.Equal(t, req.Headers.Get("x-random-id"), random)

	// data that we set
	assert.Equal(t, req.Metadata.Get(entities.MetaEprId), rule.Id)
	assert.Equal(t, req.Headers.Get(entities.HeaderIdempotencyKey), msg.Id)
	assert.Equal(t, req.Headers.Get(entities.HeaderWebhookId), msg.Id)
	assert.Equal(t, req.Headers.Get(entities.HeaderWebhookTs), fmt.Sprintf("%d", req.Timestamp))
	assert.Equal(t, req.Headers.Get(entities.HeaderWebhookRef), fmt.Sprintf("%s/%s", msg.AppId, ep.Id))
	assert.NotEmpty(t, req.Headers.Get(entities.HeaderWebhookSign))
}
