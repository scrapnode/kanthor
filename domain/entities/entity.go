package entities

import (
	"time"
)

type Entity struct {
	Id        string `json:"id"`
	Timestamp int64  `json:"timestamp"`
}

func (e *Entity) SetTS(now time.Time) {
	if e.Timestamp == 0 {
		e.Timestamp = now.UnixMilli()
	}
}

type AuditTime struct {
	// I didn't find a way to disable automatic fields modify yet
	// so, I use a tag to disable this feature here
	// but, we should keep our entities stateless if we can
	CreatedAt int64 `json:"created_at" gorm:"autoCreateTime:false"`
	UpdatedAt int64 `json:"updated_at" gorm:"autoUpdateTime:false"`
}

func (at *AuditTime) SetAT(now time.Time) {
	if at.CreatedAt == 0 {
		at.CreatedAt = now.UnixMilli()
	}
	at.UpdatedAt = now.UnixMilli()
}
