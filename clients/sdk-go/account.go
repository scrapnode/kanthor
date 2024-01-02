package kanthor

import (
	"context"

	"github.com/scrapnode/kanthor/clients/sdk-go/internal/openapi"
)

type Account struct {
	api *openapi.APIClient
}

type (
	AccountGetRes = openapi.AccountGetRes
)

func (instance *Account) Get(ctx context.Context) (*AccountGetRes, error) {
	request := instance.api.AccountAPI.AccountMeGet(ctx)
	response, res, err := request.Execute()
	if err != nil {
		return nil, errorify(err, res)
	}

	return response, nil
}
