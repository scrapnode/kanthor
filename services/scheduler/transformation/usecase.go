package transformation

import (
	"github.com/scrapnode/kanthor/domain/entities"
	usecase "github.com/scrapnode/kanthor/usecases/scheduler"
)

func MsgToArrangeReq(msg *entities.Message) *usecase.RequestArrangeReq {
	return &usecase.RequestArrangeReq{
		Message: usecase.RequestArrangeReqMessage{
			Id:       msg.Id,
			AttId:    msg.AttId,
			Tier:     msg.Tier,
			AppId:    msg.AppId,
			Type:     msg.Type,
			Metadata: msg.Metadata,
			Headers:  msg.Headers,
			Body:     msg.Body,
		},
	}
}
