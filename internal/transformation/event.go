package transformation

import (
	"fmt"

	"github.com/scrapnode/kanthor/infrastructure/streaming"
	"github.com/scrapnode/kanthor/internal/constants"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/project"
)

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
	event.Subject = project.Subject(project.Topic(constants.TopicMessage, event.AppId, event.Type))

	return event, nil
}

func EventToMessage(event *streaming.Event) (*entities.Message, error) {
	var msg entities.Message
	if err := msg.Unmarshal(event.Data); err != nil {
		return nil, err
	}
	return &msg, nil
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
	event.Subject = project.Subject(project.Topic(constants.TopicRequest, event.AppId, event.Type))

	return event, nil
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
	event.Subject = project.Subject(project.Topic(constants.TopicResponse, event.AppId, event.Type))

	return event, nil
}

func EventToRecoveryTask(event *streaming.Event) (*entities.RecoveryTask, error) {
	var rec entities.RecoveryTask
	if err := rec.Unmarshal(event.Data); err != nil {
		return nil, err
	}
	return &rec, nil
}

func EventFromRecoveryTask(task *entities.RecoveryTask) (*streaming.Event, error) {
	data, err := task.Marshal()
	if err != nil {
		return nil, err
	}

	event := &streaming.Event{
		AppId:    task.AppId,
		Type:     constants.TypeScanner,
		Id:       fmt.Sprintf("%s/%d/%d/%d", task.AppId, task.From, task.To, task.Init),
		Data:     data,
		Metadata: map[string]string{},
	}
	event.Subject = project.Subject(project.Topic(constants.TopicRecoveryTask, event.AppId, event.Type))

	return event, nil
}

func EventToAttemptTask(event *streaming.Event) (*entities.AttemptTask, error) {
	var rec entities.AttemptTask
	if err := rec.Unmarshal(event.Data); err != nil {
		return nil, err
	}
	return &rec, nil
}

func EventFromAttemptTask(task *entities.AttemptTask) (*streaming.Event, error) {
	data, err := task.Marshal()
	if err != nil {
		return nil, err
	}

	event := &streaming.Event{
		AppId:    task.AppId,
		Type:     constants.TypeScanner,
		Id:       fmt.Sprintf("%s/%d/%d/%d", task.EpId, task.From, task.To, task.Init),
		Data:     data,
		Metadata: map[string]string{},
	}
	event.Subject = project.Subject(project.Topic(constants.TopicAttemptTask, event.AppId, event.Type))

	return event, nil
}

func EventToAttempt(event *streaming.Event) (*entities.Attempt, error) {
	var rec entities.Attempt
	if err := rec.Unmarshal(event.Data); err != nil {
		return nil, err
	}
	return &rec, nil
}

func EventFromAttempt(att *entities.Attempt) (*streaming.Event, error) {
	data, err := att.Marshal()
	if err != nil {
		return nil, err
	}

	event := &streaming.Event{
		AppId:    att.AppId,
		Type:     constants.TypeEndeavor,
		Id:       att.ReqId,
		Data:     data,
		Metadata: map[string]string{},
	}
	event.Subject = project.Subject(project.Topic(constants.TopicAttempt, event.AppId, event.Type))

	return event, nil
}
