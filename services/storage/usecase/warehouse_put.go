package usecase

import (
	"context"
	"fmt"

	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/safe"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/pkg/validator"
	"github.com/sourcegraph/conc/pool"
)

type WarehousePutIn struct {
	BatchSize int
	Messages  map[string]*entities.Message
	Requests  map[string]*entities.Request
	Responses map[string]*entities.Response
}

func (in *WarehousePutIn) Validate() error {
	err := validator.Validate(
		validator.DefaultConfig,
		validator.NumberGreaterThan("batch_size", in.BatchSize, 0),
	)
	if err != nil {
		return err
	}

	err = validator.Validate(
		validator.DefaultConfig,
		validator.Map(in.Messages, func(refId string, item *entities.Message) error {
			prefix := fmt.Sprintf("messages.%s", item.Id)
			return ValidateWarehousePutInMessage(prefix, item)
		}),
	)
	if err != nil {
		return err
	}

	err = validator.Validate(
		validator.DefaultConfig,
		validator.Map(in.Requests, func(refId string, item *entities.Request) error {
			prefix := fmt.Sprintf("requests.%s", item.Id)
			return ValidateWarehousePutInRequest(prefix, item)
		}),
	)
	if err != nil {
		return err
	}

	err = validator.Validate(
		validator.DefaultConfig,
		validator.Map(in.Responses, func(refId string, item *entities.Response) error {
			prefix := fmt.Sprintf("responses.%s", item.Id)
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
	ok := safe.Slice[string]{}
	ko := safe.Map[error]{}
	count := len(in.Messages) + len(in.Requests) + len(in.Responses)
	if count == 0 {
		return &WarehousePutOut{Success: ok.Data(), Error: ko.Data()}, nil
	}

	refs, messages, requests, responses := uc.references(in)

	// hardcode the go routine to 1 because we are expecting stable throughput of database inserting
	p := pool.New().WithMaxGoroutines(1)

	for i := 0; i < len(messages); i += in.BatchSize {
		j := utils.ChunkNext(i, len(messages), in.BatchSize)

		msgs := messages[i:j]
		p.Go(func() {
			ids, err := uc.repositories.Datastore().Message().Create(ctx, msgs)
			if err != nil {
				for _, msg := range msgs {
					ko.Set(refs[msg.Id], err)
				}
				return
			}

			for _, msgId := range ids {
				ok.Append(refs[msgId])
			}
		})
	}

	for i := 0; i < len(requests); i += in.BatchSize {
		j := utils.ChunkNext(i, len(requests), in.BatchSize)

		reqs := requests[i:j]
		p.Go(func() {
			ids, err := uc.repositories.Datastore().Request().Create(ctx, reqs)
			if err != nil {
				for _, req := range reqs {
					ko.Set(refs[req.Id], err)
				}
				return
			}

			for _, reqId := range ids {
				ok.Append(refs[reqId])
			}
		})
	}

	for i := 0; i < len(responses); i += in.BatchSize {
		j := utils.ChunkNext(i, len(responses), in.BatchSize)

		resps := responses[i:j]
		p.Go(func() {
			ids, err := uc.repositories.Datastore().Response().Create(ctx, resps)
			if err != nil {
				for _, resp := range resps {
					ko.Set(refs[resp.Id], err)
				}
				return
			}

			for _, resp := range ids {
				ok.Append(refs[resp])
			}
		})
	}

	c := make(chan bool, 1)
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

func (uc *warehose) references(in *WarehousePutIn) (map[string]string, []*entities.Message, []*entities.Request, []*entities.Response) {
	refs := map[string]string{}

	var messages []*entities.Message
	var requests []*entities.Request
	var responses []*entities.Response

	for eventId, msg := range in.Messages {
		refs[msg.Id] = eventId
		messages = append(messages, msg)
	}

	for eventId, req := range in.Requests {
		refs[req.Id] = eventId
		requests = append(requests, req)
	}

	for eventId, res := range in.Responses {
		refs[res.Id] = eventId
		responses = append(responses, res)
	}

	return refs, messages, requests, responses
}
