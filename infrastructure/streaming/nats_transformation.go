package streaming

import (
	"strings"

	natscore "github.com/nats-io/nats.go"
)

func NatsMsgToEvent(msg *natscore.Msg) *Event {
	event := &Event{
		Subject:  msg.Subject,
		Id:       msg.Header.Get(MetaId),
		Data:     msg.Data,
		Metadata: map[string]string{},
	}
	for key, value := range msg.Header {
		if strings.HasPrefix(key, "Nats") {
			continue
		}
		if key == MetaId {
			continue
		}
		event.Metadata[key] = value[0]
	}
	return event
}

func NatsMsgFromEvent(subject string, event *Event) *natscore.Msg {
	msg := &natscore.Msg{
		Subject: subject,
		Header: natscore.Header{
			MetaId: []string{event.Id},
		},
		Data: event.Data,
	}
	for key, value := range event.Metadata {
		msg.Header.Set(key, value)
	}

	return msg
}
