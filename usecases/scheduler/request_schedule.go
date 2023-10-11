package scheduler

import (
	"context"
	"fmt"

	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/pkg/ds"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/pkg/validator"
	"github.com/scrapnode/kanthor/usecases/transformation"
	"github.com/sourcegraph/conc/pool"
)

type RequestScheduleReq struct {
	Requests []entities.Request
}

func (req *RequestScheduleReq) Validate() error {
	err := validator.Validate(validator.DefaultConfig, validator.SliceRequired("requests", req.Requests))
	if err != nil {
		return err
	}

	return validator.Validate(
		validator.DefaultConfig,
		validator.Array(req.Requests, func(i int, item entities.Request) error {
			prefix := fmt.Sprintf("requests[%d]", i)
			return validator.Validate(
				validator.DefaultConfig,
				validator.StringStartsWith(prefix+".id", item.Id, entities.IdNsReq),
				validator.StringStartsWith(prefix+".msg_id", item.MsgId, entities.IdNsMsg),
				validator.StringStartsWith(prefix+".ep_id", item.MsgId, entities.IdNsEp),
				validator.StringRequired(prefix+".tier", item.Tier),
				validator.StringStartsWith(prefix+".app_id", item.MsgId, entities.IdNsApp),
				validator.StringRequired(prefix+".type", item.Type),
				validator.SliceRequired(prefix+".body", item.Body),
				validator.StringRequired(prefix+".uri", item.Uri),
				validator.StringRequired(prefix+".method", item.Method),
			)
		}),
	)
}

type RequestScheduleRes struct {
	Success []string
	Error   map[string]error
}

func (uc *request) Schedule(ctx context.Context, req *RequestScheduleReq) (*RequestScheduleRes, error) {
	ok := &ds.SafeSlice[string]{}
	ko := &ds.SafeMap[error]{}

	p := pool.New().WithMaxGoroutines(uc.conf.Scheduler.Request.Schedule.Concurrency)
	for _, r := range req.Requests {
		request := r
		p.Go(func() {
			key := utils.Key(request.AppId, request.MsgId, request.EpId, request.Id)

			event, err := transformation.EventFromRequest(&request)
			if err != nil {
				ko.Set(request.Id, err)
				return
			}

			if err := uc.publisher.Pub(ctx, event); err != nil {
				ko.Set(request.Id, err)
				return
			}

			ok.Append(key)
		})

	}
	p.Wait()

	res := &RequestScheduleRes{Success: ok.Data(), Error: ko.Data()}
	return res, nil
}
