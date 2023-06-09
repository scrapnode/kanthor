package entities

import "time"

type Entity struct {
	Id string `json:"id"`
}

type AuditTime struct {
	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

type SoftDelete struct {
	DeletedAt int64 `json:"deleted_at"`
}

type TimeSeries struct {
	Timestamp int64  `json:"timestamp"`
	Bucket    string `json:"bucket"`
}

func (ts *TimeSeries) SetTS(now time.Time, layout string) {
	ts.Timestamp = now.UnixMilli()
	ts.Bucket = now.Format(layout)
}
