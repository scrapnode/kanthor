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
		Id:       msg.Id,
		Data:     data,
		Metadata: map[string]string{},
	}
	event.Subject = project.Subject(project.Topic(constants.TopicMessage, msg.AppId, msg.Type))

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
		Id:       req.Id,
		Data:     data,
		Metadata: map[string]string{},
	}
	event.Subject = project.Subject(project.Topic(constants.TopicRequest, req.AppId, req.Type))

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
	var res entities.Response
	if err := res.Unmarshal(event.Data); err != nil {
		return nil, err
	}
	return &res, nil
}

func EventFromResponse(res *entities.Response) (*streaming.Event, error) {
	data, err := res.Marshal()
	if err != nil {
		return nil, err
	}

	event := &streaming.Event{
		Id:       res.Id,
		Data:     data,
		Metadata: map[string]string{},
	}
	event.Subject = project.Subject(project.Topic(constants.TopicResponse, res.AppId, res.Type))

	return event, nil
}

func EventToRecoveryTask(event *streaming.Event) (*entities.RecoveryTask, error) {
	var task entities.RecoveryTask
	if err := task.Unmarshal(event.Data); err != nil {
		return nil, err
	}
	return &task, nil
}

func EventFromRecoveryTask(task *entities.RecoveryTask) (*streaming.Event, error) {
	data, err := task.Marshal()
	if err != nil {
		return nil, err
	}

	event := &streaming.Event{
		Id:       fmt.Sprintf("recovery_task.%s.%d.%d.%d", task.AppId, task.From, task.To, task.Init),
		Data:     data,
		Metadata: map[string]string{},
	}
	event.Subject = project.Subject(project.Topic(constants.TopicRecoveryTask, task.AppId))

	return event, nil
}

func EventToAttemptTask(event *streaming.Event) (*entities.AttemptTask, error) {
	var task entities.AttemptTask
	if err := task.Unmarshal(event.Data); err != nil {
		return nil, err
	}
	return &task, nil
}

func EventFromAttemptTask(task *entities.AttemptTask) (*streaming.Event, error) {
	data, err := task.Marshal()
	if err != nil {
		return nil, err
	}

	event := &streaming.Event{
		Id:       fmt.Sprintf("attempt_task.%s.%d.%d.%d", task.EpId, task.From, task.To, task.Init),
		Data:     data,
		Metadata: map[string]string{},
	}
	event.Subject = project.Subject(project.Topic(constants.TopicAttemptTask, task.AppId))

	return event, nil
}

func EventToAttemptTrigger(event *streaming.Event) (*entities.AttemptTrigger, error) {
	var task entities.AttemptTrigger
	if err := task.Unmarshal(event.Data); err != nil {
		return nil, err
	}
	return &task, nil
}

func EventFromAttemptTrigger(task *entities.AttemptTrigger) (*streaming.Event, error) {
	data, err := task.Marshal()
	if err != nil {
		return nil, err
	}

	event := &streaming.Event{
		Id:       fmt.Sprintf("attempt_trigger.%d.%d.%d", task.From, task.To, task.Init),
		Data:     data,
		Metadata: map[string]string{},
	}
	event.Subject = project.Subject(project.Topic(constants.TopicAttemptTrigger, constants.TypeEndeavor))

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
		Id:       att.Id(),
		Data:     data,
		Metadata: map[string]string{},
	}
	event.Subject = project.Subject(project.Topic(constants.TopicAttempt, att.AppId))

	return event, nil
}
