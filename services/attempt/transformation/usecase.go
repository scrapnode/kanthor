package transformation

import (
	"time"

	"github.com/scrapnode/kanthor/domain/entities"
	usecase "github.com/scrapnode/kanthor/usecases/attempt"
	"github.com/scrapnode/kanthor/usecases/attempt/repos"
)

func ApplicationTriggerReq(scanSize, publishSize int) *usecase.ApplicationTriggerReq {
	return &usecase.ApplicationTriggerReq{ScanSize: scanSize, PublishSize: publishSize}
}

func ApplicationToTriggerScanReq(app *entities.Application, from, to time.Time) *usecase.TriggerScanReq {
	return &usecase.TriggerScanReq{AppId: app.Id, From: from, To: to}
}

func MsgIdsToTriggerScheduleReq(app *entities.Application, msgIds []string) *usecase.TriggerScheduleReq {
	return &usecase.TriggerScheduleReq{AppId: app.Id, MsgIds: msgIds}
}

func RequestsToTriggerScheduleReq(requests []repos.Req) *usecase.TriggerCreateReq {
	return &usecase.TriggerCreateReq{Requests: requests}
}
