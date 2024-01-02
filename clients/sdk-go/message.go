package kanthor

import (
	"context"

	"github.com/scrapnode/kanthor/clients/sdk-go/internal/openapi"
)

type Message struct {
	api *openapi.APIClient
}

type (
	MessageCreateReq = openapi.MessageCreateReq
	MessageCreateRes = openapi.MessageCreateRes
)

func (instance *Message) Create(ctx context.Context, req *MessageCreateReq) (*MessageCreateRes, error) {
	request := instance.api.MessageAPI.MessagePost(ctx)
	response, res, err := request.Payload(*req).Execute()
	if err != nil {
		return nil, errorify(err, res)
	}

	return response, nil
}
