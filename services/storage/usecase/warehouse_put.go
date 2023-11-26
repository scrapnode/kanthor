package usecase

import (
	"context"
	"fmt"

	"github.com/scrapnode/kanthor/internal/domain/entities"
	"github.com/scrapnode/kanthor/pkg/safe"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/pkg/validator"
	"github.com/sourcegraph/conc/pool"
)

type WarehousePutIn struct {
	Size      int
	Messages  []entities.Message
	Requests  []entities.Request
	Responses []entities.Response
}

func (in *WarehousePutIn) Validate() error {
	err := validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThan("size", in.Size, 0),
	)
	if err != nil {
		return err
	}

	err = validator.Validate(
		validator.DefaultConfig,
		validator.Slice(in.Messages, func(i int, item *entities.Message) error {
			prefix := fmt.Sprintf("messages[%d]", i)
			return ValidateWarehousePutInMessage(prefix, item)
		}),
	)
	if err != nil {
		return err
	}

	err = validator.Validate(
		validator.DefaultConfig,
		validator.Slice(in.Requests, func(i int, item *entities.Request) error {
			prefix := fmt.Sprintf("requests[%d]", i)
			return ValidateWarehousePutInRequest(prefix, item)
		}),
	)
	if err != nil {
		return err
	}

	err = validator.Validate(
		validator.DefaultConfig,
		validator.Slice(in.Responses, func(i int, item *entities.Response) error {
			prefix := fmt.Sprintf("responses[%d]", i)
			return ValidateWarehousePutInResponse(prefix, item)
		}),
	)
	if err != nil {
		return err
	}

	return nil
}

func ValidateWarehousePutInMessage(prefix string, message *entities.Message) error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith(prefix+".id", message.Id, entities.IdNsMsg),
		validator.NumberGreaterThan(prefix+".timestamp", message.Timestamp, 0),
		validator.StringRequired(prefix+".tier", message.Tier),
		validator.StringStartsWith(prefix+".app_id", message.AppId, entities.IdNsApp),
		validator.StringRequired(prefix+".type", message.Type),
		validator.StringRequired(prefix+".body", message.Body),
	)
}

func ValidateWarehousePutInRequest(prefix string, request *entities.Request) error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith(prefix+".id", request.Id, entities.IdNsReq),
		validator.NumberGreaterThan(prefix+".timestamp", request.Timestamp, 0),
		validator.StringStartsWith(prefix+".msg_id", request.MsgId, entities.IdNsMsg),
		validator.StringStartsWith(prefix+".ep_id", request.EpId, entities.IdNsEp),
		validator.StringRequired(prefix+".tier", request.Tier),
		validator.StringStartsWith(prefix+".app_id", request.AppId, entities.IdNsApp),
		validator.StringRequired(prefix+".type", request.Type),
		validator.StringRequired(prefix+".body", request.Body),
		validator.StringRequired(prefix+".uri", request.Uri),
		validator.StringRequired(prefix+".method", request.Method),
	)
}

func ValidateWarehousePutInResponse(prefix string, response *entities.Response) error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith(prefix+".id", response.Id, entities.IdNsRes),
		validator.NumberGreaterThan(prefix+".timestamp", response.Timestamp, 0),
		validator.StringStartsWith(prefix+".msg_id", response.MsgId, entities.IdNsMsg),
		validator.StringStartsWith(prefix+".ep_id", response.EpId, entities.IdNsEp),
		validator.StringStartsWith(prefix+".req_id", response.ReqId, entities.IdNsReq),
		validator.StringRequired(prefix+".tier", response.Tier),
		validator.StringStartsWith(prefix+".app_id", response.AppId, entities.IdNsApp),
		validator.StringRequired(prefix+".type", response.Type),
	)
}

type WarehousePutOut struct {
	Success []string
	Error   map[string]error
}

func (uc *warehose) Put(ctx context.Context, in *WarehousePutIn) (*WarehousePutOut, error) {
	count := len(in.Messages) + len(in.Requests) + len(in.Responses)
	if count == 0 {
		return &WarehousePutOut{Success: []string{}, Error: map[string]error{}}, nil
	}

	ok := safe.Slice[string]{}
	ko := safe.Map[error]{}

	// hardcode the go routine to 1 because we are expecting stable throughput of database inserting
	p := pool.New().WithMaxGoroutines(1)
	for i := 0; i < len(in.Messages); i += in.Size {
		j := utils.ChunkNext(i, len(in.Messages), in.Size)

		messages := in.Messages[i:j]
		p.Go(func() {
			ids, err := uc.repositories.Message().Create(ctx, messages)
			if err != nil {
				for _, message := range messages {
					ko.Set(message.Id, err)
				}
				return
			}

			ok.Append(ids...)
		})
	}

	for i := 0; i < len(in.Requests); i += in.Size {
		j := utils.ChunkNext(i, len(in.Requests), in.Size)

		requests := in.Requests[i:j]
		p.Go(func() {
			ids, err := uc.repositories.Request().Create(ctx, requests)
			if err != nil {
				for _, request := range requests {
					ko.Set(request.Id, err)
				}
				return
			}

			ok.Append(ids...)
		})
	}

	for i := 0; i < len(in.Responses); i += in.Size {
		j := utils.ChunkNext(i, len(in.Responses), in.Size)

		responses := in.Responses[i:j]
		p.Go(func() {
			ids, err := uc.repositories.Response().Create(ctx, responses)
			if err != nil {
				for _, response := range responses {
					ko.Set(response.Id, err)
				}
				return
			}

			ok.Append(ids...)
		})
	}

	c := make(chan bool)
	defer close(c)

	go func() {
		p.Wait()
		c <- true
	}()

	select {
	case <-c:
		return &WarehousePutOut{Success: ok.Data(), Error: ko.Data()}, nil
	case <-ctx.Done():
		// context deadline exceeded, should set that error to remain messages
		for _, message := range in.Messages {
			// no error, should add context deadline error
			if _, has := ko.Get(message.Id); !has {
				ko.Set(message.Id, ctx.Err())
			}
		}

		// context deadline exceeded, should set that error to remain requests
		for _, request := range in.Requests {
			// no error, should add context deadline error
			if _, has := ko.Get(request.Id); !has {
				ko.Set(request.Id, ctx.Err())
			}
		}

		// context deadline exceeded, should set that error to remain responses
		for _, response := range in.Responses {
			// no error, should add context deadline error
			if _, has := ko.Get(response.Id); !has {
				ko.Set(response.Id, ctx.Err())
			}
		}

		return &WarehousePutOut{Success: []string{}, Error: ko.Data()}, nil
	}
}
