package scheduler

import (
	"context"
	"fmt"

	"github.com/samber/lo"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/pkg/validator"
	"github.com/scrapnode/kanthor/usecases/transformation"
	"github.com/sourcegraph/conc"
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
	res := &RequestScheduleRes{Success: []string{}, Error: map[string]error{}}
	chunks := lo.Chunk(req.Requests, uc.conf.Scheduler.Request.Schedule.Concurrency)

	for _, requests := range chunks {
		resp, err := uc.schedule(ctx, requests)
		if err != nil {
			for _, request := range requests {
				res.Error[request.Id] = err
			}
			continue
		}

		res.Success = append(res.Success, resp.Success...)
		utils.SliceMerge[error](res.Error, resp.Error)
	}

	return res, nil
}

func (uc *request) schedule(ctx context.Context, requests []entities.Request) (*RequestScheduleRes, error) {
	res := &RequestScheduleRes{Success: []string{}, Error: map[string]error{}}
	var wg conc.WaitGroup
	for _, entity := range requests {
		r := entity
		wg.Go(func() {
			key := utils.Key(r.AppId, r.Id, r.Id)

			event, err := transformation.EventFromRequest(&r)
			if err == nil {
				err = uc.publisher.Pub(ctx, event)
			}

			if err == nil {
				res.Success = append(res.Success, key)
			} else {
				res.Error[key] = err
				uc.logger.Errorw(err.Error(), "key", key)
			}
		})
	}
	wg.Wait()

	return res, nil
}
