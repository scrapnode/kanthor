package kanthor

import (
	"net/http"

	"github.com/scrapnode/kanthor/clients/sdk-go/internal/openapi"
)

type Error struct {
	status int
	body   []byte
	error  string
}

func (e Error) Error() string {
	return e.error
}

func (e Error) Body() []byte {
	return e.body
}

func (e Error) Status() int {
	return e.status
}

func errorify(err error, res *http.Response) error {
	if openapiError, ok := err.(openapi.GenericOpenAPIError); ok {
		e := &Error{
			body:  openapiError.Body(),
			error: openapiError.Error(),
		}
		if res != nil {
			e.status = res.StatusCode
		}
		return e
	}
	return err
}
