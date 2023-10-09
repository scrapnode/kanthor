package attempt

import (
	"context"
	"fmt"
	"time"

	"github.com/samber/lo"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/pkg/validator"
	"github.com/scrapnode/kanthor/usecases/attempt/repos"
)

type TriggerCreateReq struct {
	Requests []repos.Req
}

func (req *TriggerCreateReq) Validate() error {
	err := validator.Validate(validator.DefaultConfig, validator.SliceRequired("requests", req.Requests))
	if err != nil {
		return err
	}
	return validator.Validate(
		validator.DefaultConfig,
		validator.Array(req.Requests, func(i int, item repos.Req) error {
			prefix := fmt.Sprintf("requests[%d]", i)
			return validator.Validate(
				validator.DefaultConfig,
				validator.StringStartsWith(prefix+".app_id", item.AppId, entities.IdNsApp),
				validator.StringStartsWith(prefix+".msg_id", item.MsgId, entities.IdNsMsg),
				validator.StringStartsWith(prefix+".ep_id", item.EpId, entities.IdNsEp),
				validator.StringStartsWith(prefix+".id", item.Id, entities.IdNsReq),
				validator.StringRequired(prefix+".tier", item.Tier),
			)
		}),
	)
}

type TriggerCreateRes struct {
	Success []string
	Error   map[string]error
}

func (uc *trigger) Create(ctx context.Context, req *TriggerCreateReq) (*TriggerCreateRes, error) {
	res := &TriggerCreateRes{Success: []string{}, Error: map[string]error{}}
	chunks := lo.Chunk(req.Requests, uc.conf.Attempt.Trigger.Create.Concurrency)

	for _, requests := range chunks {
		resp, err := uc.create(ctx, requests)
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

func (uc *trigger) create(ctx context.Context, requests []repos.Req) (*TriggerCreateRes, error) {
	res := &TriggerCreateRes{Success: []string{}, Error: map[string]error{}}

	attempts := []entities.Attempt{}
	for _, request := range requests {
		attempts = append(attempts, entities.Attempt{
			ReqId:        request.Id,
			Tier:         request.Tier,
			ScheduleNext: uc.timer.Now().Add(time.Duration(uc.conf.Attempt.Trigger.Create.ScheduleDelay) * time.Second).UnixMilli(),
			ScheduledAt:  uc.timer.Now().UnixMilli(),
		})
	}

	ids, err := uc.repos.Attempt().Create(ctx, attempts)
	if err != nil {
		return nil, err
	}
	res.Success = append(res.Success, ids...)

	return res, nil
}
