package entities

import "time"

type Entity struct {
	Id string `json:"id" validate:"required"`
}

type AuditTime struct {
	ModifiedBy string `json:"modified_by" validate:"required"`
	CreatedAt  int64  `json:"created_at" validate:"required"`
	UpdatedAt  int64  `json:"updated_at" validate:"required"`
}

func (at *AuditTime) SetAT(now time.Time) {
	if at.CreatedAt == 0 {
		at.CreatedAt = now.UnixMilli()
	}
	if at.UpdatedAt == 0 {
		at.UpdatedAt = now.UnixMilli()
	}
}

type TimeSeries struct {
	Timestamp int64  `json:"timestamp" validate:"required"`
	Bucket    string `json:"bucket" validate:"required"`
}

func (ts *TimeSeries) SetTS(now time.Time, layout string) {
	if ts.Timestamp == 0 {
		ts.Timestamp = now.UnixMilli()
	}
	ts.Bucket = now.Format(layout)
}
