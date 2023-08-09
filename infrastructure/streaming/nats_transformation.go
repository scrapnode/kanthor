package streaming

import (
	"github.com/nats-io/nats.go"
	"strings"
)

func natsMsgToEvent(msg *nats.Msg) Event {
	event := Event{
		Subject:  msg.Subject,
		AppId:    msg.Header.Get(MetaAppId),
		Type:     msg.Header.Get(MetaType),
		Id:       msg.Header.Get(MetaId),
		Data:     msg.Data,
		Metadata: map[string]string{},
	}
	for key, value := range msg.Header {
		if strings.HasPrefix(key, "Nats") {
			continue
		}
		if key == MetaAppId || key == MetaType || key == MetaId {
			continue
		}
		event.Metadata[key] = value[0]
	}
	return event
}

func natsMsgFromEvent(subject string, event *Event) *nats.Msg {
	msg := &nats.Msg{
		Subject: subject,
		Header: nats.Header{
			// for deduplicate purpose
			MetaAppId: []string{event.AppId},
			MetaType:  []string{event.Type},
			MetaId:    []string{event.Id},
		},
		Data: event.Data,
	}
	for key, value := range event.Metadata {
		msg.Header.Set(key, value)
	}

	return msg
}
