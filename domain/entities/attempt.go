package entities

import (
	"encoding/json"
	"github.com/scrapnode/kanthor/pkg/utils"
)

var (
	AttStatusFail    = -1
	AttStatusPending = 0
	AttStatusSuccess = 1
)

type Attempt struct {
	TSEntity

	MessageId  string `json:"message_id"`
	RequestId  string `json:"request_id"`
	ResponseId string `json:"response_id"`

	Status      int   `json:"status"`
	NextRetryTs int64 `json:"next_retry_ts"`
}

func (entity *Attempt) TableName() string {
	return "kanthor_attempt"
}

func (entity *Attempt) GenId() {
	if entity.Id == "" {
		entity.Id = utils.ID("att")
	}
}

func (entity *Attempt) Marshal() ([]byte, error) {
	return json.Marshal(entity)
}

func (entity *Attempt) Unmarshal(data []byte) error {
	return json.Unmarshal(data, entity)
}
