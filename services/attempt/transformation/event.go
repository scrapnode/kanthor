package transformation

import (
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
)

func EventToApplication(event *streaming.Event) (*entities.Application, error) {
	var app entities.Application
	if err := app.Unmarshal(event.Data); err != nil {
		return nil, err
	}
	return &app, nil
}
