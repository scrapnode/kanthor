package dispatcher

import (
	"context"
	"fmt"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/circuitbreaker"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/pkg/sender"
)

func (usecase *forwarder) Send(ctx context.Context, req *ForwarderSendReq) (*ForwarderSendRes, error) {
	request := &sender.Request{
		Method:  req.Request.Method,
		Headers: req.Request.Headers,
		Uri:     req.Request.Uri,
		Body:    req.Request.Body,
	}
	request.Headers.Set("Idempotency-Key", req.Request.Id)

	response, err := circuitbreaker.Do[sender.Response](
		usecase.cb,
		req.Request.EndpointId,
		func() (interface{}, error) {
			return usecase.dispatch(request)
		},
		func(err error) error {
			return err
		},
	)

	res := &ForwarderSendRes{
		Response: entities.Response{
			Tier:     req.Request.Tier,
			AppId:    req.Request.AppId,
			Type:     req.Request.Type,
			Metadata: req.Request.Metadata,
		},
	}
	res.Response.GenId()
	res.Response.SetTS(usecase.timer.Now(), usecase.conf.Bucket.Layout)
	res.Response.Metadata[entities.MetaReqId] = req.Request.Id
	res.Response.Metadata[entities.MetaReqBucket] = req.Request.Bucket
	res.Response.Metadata[entities.MetaReqTs] = fmt.Sprintf("%d", req.Request.Timestamp)

	// either error was happened or not, we need to publish response event, so we can handle custom logic later
	// example use cases are retry, notification, i.e
	if err == nil {
		res.Response.Status = response.Status
		res.Response.Uri = response.Uri
		res.Response.Headers = response.Headers
		res.Response.Body = response.Body
	} else {
		usecase.logger.Errorw(err.Error(), "ep_id", req.Request.EndpointId, "req_id", req.Request.Id)
		res.Response.Status = entities.ResponseStatusErr
		res.Response.Error = err.Error()
	}

	event, err := transformResponse2Event(&res.Response)
	if err != nil {
		usecase.logger.Errorw(err.Error(), "ep_id", req.Request.EndpointId, "req_id", req.Request.Id)
		return nil, fmt.Errorf("unable transform response to event [%s/%s]", req.Request.EndpointId, req.Request.Id)
	}
	if err := usecase.publisher.Pub(ctx, event); err != nil {
		usecase.logger.Errorw(err.Error(), "ep_id", req.Request.EndpointId, "req_id", req.Request.Id)
		return nil, fmt.Errorf("unable publish event for response of request [%s/%s]", req.Request.EndpointId, req.Request.Id)
	}

	return res, nil
}

func transformResponse2Event(res *entities.Response) (*streaming.Event, error) {
	data, err := res.Marshal()
	if err != nil {
		return nil, err
	}

	event := &streaming.Event{
		AppId:    res.AppId,
		Type:     res.Type,
		Data:     data,
		Metadata: map[string]string{},
	}
	event.GenId()
	event.Subject = streaming.Subject(
		streaming.Namespace,
		res.Tier,
		streaming.TopicRes,
		event.AppId,
		event.Type,
	)

	return event, nil
}