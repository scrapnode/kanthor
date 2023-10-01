package transformation

import (
	"github.com/scrapnode/kanthor/domain/entities"
	usecase "github.com/scrapnode/kanthor/usecases/dispatcher"
)

func ReqToSendReq(req *entities.Request) *usecase.ForwarderSendReq {
	return &usecase.ForwarderSendReq{
		Request: usecase.ForwarderSendReqRequest{
			Id:       req.Id,
			MsgId:    req.MsgId,
			EpId:     req.EpId,
			Tier:     req.Tier,
			AppId:    req.AppId,
			Type:     req.Type,
			Metadata: req.Metadata,
			Headers:  req.Headers,
			Body:     req.Body,
			Uri:      req.Uri,
			Method:   req.Method,
		},
	}
}
