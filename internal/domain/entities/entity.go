package entities

import (
	"time"
)

type Entity struct {
	Id string `json:"id"`
}

// TSEntity is time series entity
type TSEntity struct {
	Id        string `json:"id"`
	Timestamp int64  `json:"timestamp"`
}

func (entity *TSEntity) SetTS(now time.Time) {
	if entity.Timestamp == 0 {
		entity.Timestamp = now.UnixMilli()
	}
}

type AuditTime struct {
	// I didn't find a way to disable automatic fields modify yet
	// so, I use a tag to disable this feature here
	// but, we should keep our entities stateless if we can
	CreatedAt int64 `json:"created_at" gorm:"autoCreateTime:false"`
	UpdatedAt int64 `json:"updated_at" gorm:"autoUpdateTime:false"`
}

func (entity *AuditTime) SetAT(now time.Time) {
	if entity.CreatedAt == 0 {
		entity.CreatedAt = now.UnixMilli()
	}
	entity.UpdatedAt = now.UnixMilli()
}