package attempt

import (
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
)

func EventFromNotification(noti *entities.AttemptNotification) (*streaming.Event, error) {
	data, err := noti.Marshal()
	if err != nil {
		return nil, err
	}

	event := &streaming.Event{
		AppId:    noti.AppId,
		Type:     "internal.attempt.trigger.notification",
		Id:       noti.AppId,
		Data:     data,
		Metadata: map[string]string{},
	}
	event.Subject = streaming.Subject(
		streaming.Namespace,
		noti.Tier,
		streaming.TopicApp,
		event.AppId,
		event.Type,
	)

	return event, nil
}

func EventToNotification(event *streaming.Event) (*entities.AttemptNotification, error) {
	var noti entities.AttemptNotification
	if err := noti.Unmarshal(event.Data); err != nil {
		return nil, err
	}

	return &noti, nil
}
