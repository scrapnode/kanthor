package transformation

import (
	"fmt"

	"github.com/scrapnode/kanthor/domain/constants"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/namespace"
)

func EventToMessage(event *streaming.Event) (*entities.Message, error) {
	var msg entities.Message
	if err := msg.Unmarshal(event.Data); err != nil {
		return nil, err
	}
	return &msg, nil
}

func EventToRequest(event *streaming.Event) (*entities.Request, error) {
	var req entities.Request
	if err := req.Unmarshal(event.Data); err != nil {
		return nil, err
	}
	return &req, nil
}

func EventToResponse(event *streaming.Event) (*entities.Response, error) {
	var req entities.Response
	if err := req.Unmarshal(event.Data); err != nil {
		return nil, err
	}
	return &req, nil
}

func EventFromMessage(msg *entities.Message) (*streaming.Event, error) {
	data, err := msg.Marshal()
	if err != nil {
		return nil, err
	}

	event := &streaming.Event{
		AppId:    msg.AppId,
		Type:     msg.Type,
		Id:       msg.Id,
		Data:     data,
		Metadata: map[string]string{},
	}
	event.Subject = streaming.Subject(
		namespace.Namespace(),
		msg.Tier,
		constants.TopicMessage,
		event.AppId,
		event.Type,
	)

	return event, nil
}

func EventFromRequest(req *entities.Request) (*streaming.Event, error) {
	data, err := req.Marshal()
	if err != nil {
		return nil, err
	}

	event := &streaming.Event{
		AppId:    req.AppId,
		Type:     req.Type,
		Id:       req.Id,
		Data:     data,
		Metadata: map[string]string{},
	}
	event.Subject = streaming.Subject(
		namespace.Namespace(),
		req.Tier,
		constants.TopicRequest,
		event.AppId,
		event.Type,
	)

	return event, nil
}

func EventFromResponse(res *entities.Response) (*streaming.Event, error) {
	data, err := res.Marshal()
	if err != nil {
		return nil, err
	}

	event := &streaming.Event{
		AppId:    res.AppId,
		Type:     res.Type,
		Id:       res.Id,
		Data:     data,
		Metadata: map[string]string{},
	}
	event.Subject = streaming.Subject(
		namespace.Namespace(),
		res.Tier,
		constants.TopicResponse,
		event.AppId,
		event.Type,
	)

	return event, nil
}

func EventFromTrigger(trigger *entities.AttemptTrigger) (*streaming.Event, error) {
	data, err := trigger.Marshal()
	if err != nil {
		return nil, err
	}

	event := &streaming.Event{
		AppId:    trigger.AppId,
		Type:     constants.TypeInternal,
		Id:       fmt.Sprintf("%s/%d/%d", trigger.AppId, trigger.From, trigger.To),
		Data:     data,
		Metadata: map[string]string{},
	}
	event.Subject = streaming.Subject(
		namespace.Namespace(),
		trigger.Tier,
		constants.TopicTrigger,
		event.AppId,
		event.Type,
	)

	return event, nil
}

func EventToTrigger(event *streaming.Event) (*entities.AttemptTrigger, error) {
	var noti entities.AttemptTrigger
	if err := noti.Unmarshal(event.Data); err != nil {
		return nil, err
	}

	return &noti, nil
}
