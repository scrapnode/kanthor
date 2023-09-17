package transformation

import (
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
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
		streaming.Namespace,
		msg.Tier,
		streaming.TopicMsg,
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
		streaming.Namespace,
		req.Tier,
		streaming.TopicReq,
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
		streaming.Namespace,
		res.Tier,
		streaming.TopicRes,
		event.AppId,
		event.Type,
	)

	return event, nil
}
