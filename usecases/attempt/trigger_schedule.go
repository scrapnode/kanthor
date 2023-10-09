package attempt

import (
	"context"

	"github.com/samber/lo"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/internal/planner"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/pkg/validator"
	"github.com/scrapnode/kanthor/usecases/transformation"
	"github.com/sourcegraph/conc"
)

type TriggerScheduleReq struct {
	AppId  string
	MsgIds []string
}

func (req *TriggerScheduleReq) Validate() error {
	return validator.Validate(
		validator.DefaultConfig,
		validator.StringStartsWith("app_id", req.AppId, entities.IdNsApp),
		validator.SliceRequired("msg_ids", req.MsgIds),
	)
}

type TriggerScheduleRes struct {
	Success []string
	Error   map[string]error
}

func (uc *trigger) Schedule(ctx context.Context, req *TriggerScheduleReq) (*TriggerScheduleRes, error) {
	app, err := uc.applicable(ctx, req.AppId)
	if err != nil {
		return nil, err
	}

	res := &TriggerScheduleRes{Success: []string{}, Error: map[string]error{}}
	chunks := lo.Chunk(req.MsgIds, uc.conf.Attempt.Trigger.Schedule.Concurrency)

	for _, msgIds := range chunks {
		resp, err := uc.schedule(ctx, app, msgIds)
		if err != nil {
			for _, msgId := range msgIds {
				res.Error[msgId] = err
			}
			continue
		}

		res.Success = append(res.Success, resp.Success...)
		utils.SliceMerge[error](res.Error, resp.Error)
	}

	return res, nil
}

func (uc *trigger) schedule(ctx context.Context, app *planner.Applicable, msgIds []string) (*TriggerScheduleRes, error) {
	messages, err := uc.repos.Message().ListByIds(ctx, msgIds)
	if err != nil {
		return nil, err
	}

	requests := []entities.Request{}
	for _, message := range messages {
		reqs, logs := planner.Requests(&message, app, uc.timer)
		if len(logs) > 0 {
			for _, l := range logs {
				uc.logger.Warnw(l[0].(string), l[1:]...)
			}
		}
		requests = append(requests, reqs...)
	}

	res := &TriggerScheduleRes{Success: []string{}, Error: map[string]error{}}
	if len(requests) == 0 {
		return res, nil
	}

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
