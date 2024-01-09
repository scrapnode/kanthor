package routing

import (
	"fmt"

	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/identifier"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/pkg/utils"
)

func NewRequest(
	timer timer.Timer,
	msg *entities.Message,
	ep *entities.Endpoint,
	epr *entities.EndpointRule,
) *entities.Request {
	// construct request
	req := &entities.Request{
		MsgId:    msg.Id,
		Tier:     msg.Tier,
		AppId:    msg.AppId,
		Type:     msg.Type,
		EpId:     ep.Id,
		Metadata: entities.Metadata{},
		Headers:  entities.Header{},
		Body:     msg.Body,
		Uri:      ep.Uri,
		Method:   ep.Method,
	}

	// must use merge function otherwise you will edit the original data
	req.Headers.Merge(msg.Headers)
	req.Metadata.Merge(msg.Metadata)
	req.Id = identifier.New(entities.IdNsReq)
	req.SetTS(timer.Now())

	req.Metadata.Set(entities.MetaEprId, epr.Id)
	req.Headers.Set(entities.HeaderIdempotencyKey, msg.Id)

	// https://github.com/standard-webhooks/standard-webhooks/blob/main/spec/standard-webhooks.md
	req.Headers.Set(entities.HeaderWebhookId, msg.Id)
	req.Headers.Set(entities.HeaderWebhookTs, fmt.Sprintf("%d", req.Timestamp))
	signature := utils.SignatureSign(ep.SecretKey, fmt.Sprintf("%s.%d.%s", msg.Id, req.Timestamp, msg.Body))
	req.Headers.Set(entities.HeaderWebhookSign, fmt.Sprintf("v1,%s", signature))

	// custom headers
	req.Headers.Set(entities.HeaderWebhookRef, fmt.Sprintf("%s/%s", msg.AppId, ep.Id))

	return req
}
