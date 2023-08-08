package transformation

import (
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/infrastructure/streaming"
)

func EventToRequest(event streaming.Event) (*entities.Request, error) {
	var req entities.Request
	if err := req.Unmarshal(event.Data); err != nil {
		return nil, err
	}
	return &req, nil
}
