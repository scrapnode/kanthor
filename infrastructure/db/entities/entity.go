package entities

import (
	"time"
)

type Entity struct {
	Id string `json:"id"`
}

type AuditTime struct {
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type SoftDelete struct {
	DeletedAt *time.Time `json:"deleted_at"`
}

type TimeSeries struct {
	Timestamp time.Time `json:"timestamp"`
	Bucket    string    `json:"bucket"`
}

func (entity *TimeSeries) GenBucket(layout string) {
	entity.Timestamp = time.Now().UTC()
	entity.Bucket = entity.Timestamp.Format(layout)
}
