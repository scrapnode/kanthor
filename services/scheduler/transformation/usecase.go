package transformation

import (
	"github.com/scrapnode/kanthor/domain/entities"
	usecase "github.com/scrapnode/kanthor/usecases/scheduler"
)

func MessageToRequestArrangeReq(msg *entities.Message) *usecase.RequestArrangeReq {
	return &usecase.RequestArrangeReq{Message: msg}
}

func RequestToRequestScheduleReq(reqs []entities.Request) *usecase.RequestScheduleReq {
	return &usecase.RequestScheduleReq{Requests: reqs}
}
