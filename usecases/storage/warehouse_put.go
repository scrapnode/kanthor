package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/pkg/safe"
	"github.com/scrapnode/kanthor/pkg/validator"
	"github.com/sourcegraph/conc"
)

type WarehousePutReq struct {
	ChunkTimeout int64
	ChunkSize    int
	Messages     []entities.Message
	Requests     []entities.Request
	Responses    []entities.Response
}

func (req *WarehousePutReq) Validate() error {
	err := validator.Validate(
		validator.DefaultConfig,
		validator.SliceRequired("messages", req.Messages),
		validator.SliceRequired("requests", req.Requests),
		validator.SliceRequired("responses", req.Responses),
		validator.NumberGreaterThan("chunk_timeout", req.ChunkTimeout, 1000),
		validator.NumberGreaterThan("chunk_size", req.ChunkTimeout, 1),
	)
	if err != nil {
		return err
	}

	err = validator.Validate(
		validator.DefaultConfig,
		validator.Array(req.Messages, func(i int, item *entities.Message) error {
			prefix := fmt.Sprintf("messages[%d]", i)
			return ValidateWarehousePutReqMessage(prefix, item)
		}),
	)
	if err != nil {
		return err
	}

	err = validator.Validate(
		validator.DefaultConfig,
		validator.Array(req.Requests, func(i int, item *entities.Request) error {
			prefix := fmt.Sprintf("requests[%d]", i)
			return ValidateWarehousePutReqRequest(prefix, item)
		}),
	)
	if err != nil {
		return err
	}

	err = validator.Validate(
		validator.DefaultConfig,
		validator.Array(req.Responses, func(i int, item *entities.Response) error {
			prefix := fmt.Sprintf("responses[%d]", i)
			return ValidateWarehousePutReqResponse(prefix, item)
		}),
	)
	if err != nil {
		return err
	}

	return nil
}

func ValidateWarehousePutReqMessage(prefix string, message *entities.Message) error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith(prefix+".id", message.Id, entities.IdNsMsg),
		validator.NumberGreaterThan(prefix+".timestamp", message.Timestamp, 0),
		validator.StringRequired(prefix+".tier", message.Tier),
		validator.StringStartsWith(prefix+".app_id", message.AppId, entities.IdNsApp),
		validator.StringRequired(prefix+".type", message.Type),
		validator.SliceRequired(prefix+".body", message.Body),
	)
}

func ValidateWarehousePutReqRequest(prefix string, request *entities.Request) error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith(prefix+".id", request.Id, entities.IdNsReq),
		validator.NumberGreaterThan(prefix+".timestamp", request.Timestamp, 0),
		validator.StringStartsWith(prefix+".msg_id", request.MsgId, entities.IdNsMsg),
		validator.StringStartsWith(prefix+".ep_id", request.EpId, entities.IdNsEp),
		validator.StringRequired(prefix+".tier", request.Tier),
		validator.StringStartsWith(prefix+".app_id", request.AppId, entities.IdNsApp),
		validator.StringRequired(prefix+".type", request.Type),
		validator.SliceRequired(prefix+".body", request.Body),
		validator.StringRequired(prefix+".uri", request.Uri),
		validator.StringRequired(prefix+".method", request.Method),
	)
}

func ValidateWarehousePutReqResponse(prefix string, response *entities.Response) error {
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

type WarehousePutRes struct {
	Success []string
	Error   map[string]error
}

func (uc *warehose) Put(ctx context.Context, req *WarehousePutReq) (*WarehousePutRes, error) {
	count := len(req.Messages) + len(req.Requests) + len(req.Responses)
	if count == 0 {
		return &WarehousePutRes{Success: []string{}, Error: map[string]error{}}, nil
	}

	ok := safe.Map[string]{}
	ko := safe.Map[error]{}

	// timeout duration will be scaled based on how many requests you have
	duration := time.Duration(req.ChunkTimeout * int64(count+1))
	timeout, cancel := context.WithTimeout(ctx, time.Millisecond*duration)
	defer cancel()
	var wg conc.WaitGroup

	for i := 0; i < len(req.Messages); i += req.ChunkSize {
		j := i + req.ChunkSize
		messages := req.Messages[i:j]
		wg.Go(func() {
			records, err := uc.repos.Message().Create(ctx, messages)
			if err != nil {
				for _, message := range req.Messages[i:j] {
					ko.Set(message.Id, err)
				}
				return
			}

			for _, record := range records {
				ok.Set(record.Id, record.Id)
			}
		})
	}

	for i := 0; i < len(req.Requests); i += req.ChunkSize {
		j := i + req.ChunkSize
		requests := req.Requests[i:j]
		wg.Go(func() {
			records, err := uc.repos.Request().Create(ctx, requests)
			if err != nil {
				for _, request := range req.Requests[i:j] {
					ko.Set(request.Id, err)
				}
				return
			}

			for _, record := range records {
				ok.Set(record.Id, record.Id)
			}
		})
	}

	for i := 0; i < len(req.Responses); i += req.ChunkSize {
		j := i + req.ChunkSize
		responses := req.Responses[i:j]
		wg.Go(func() {
			records, err := uc.repos.Response().Create(ctx, responses)
			if err != nil {
				for _, response := range req.Responses[i:j] {
					ko.Set(response.Id, err)
				}
				return
			}

			for _, record := range records {
				ok.Set(record.Id, record.Id)
			}
		})
	}

	c := make(chan bool)
	defer close(c)

	go func() {
		wg.Wait()
		c <- true
	}()

	select {
	case <-c:
		return &WarehousePutRes{Success: ok.Keys(), Error: ko.Data()}, nil
	case <-timeout.Done():
		// context deadline exceeded, should set that error to remain messages
		for _, message := range req.Messages {
			if _, success := ok.Get(message.Id); success {
				// already success, should not retry it
				continue
			}

			// no error, should add context deadline error
			if _, has := ko.Get(message.Id); !has {
				ko.Set(message.Id, ctx.Err())
			}
		}

		// context deadline exceeded, should set that error to remain requests
		for _, request := range req.Requests {
			if _, success := ok.Get(request.Id); success {
				// already success, should not retry it
				continue
			}

			// no error, should add context deadline error
			if _, has := ko.Get(request.Id); !has {
				ko.Set(request.Id, ctx.Err())
			}
		}

		// context deadline exceeded, should set that error to remain responses
		for _, response := range req.Responses {
			if _, success := ok.Get(response.Id); success {
				// already success, should not retry it
				continue
			}

			// no error, should add context deadline error
			if _, has := ko.Get(response.Id); !has {
				ko.Set(response.Id, ctx.Err())
			}
		}

		return &WarehousePutRes{Success: ok.Keys(), Error: ko.Data()}, nil
	}
}
