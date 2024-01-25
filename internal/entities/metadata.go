package entities

import (
	"database/sql/driver"
	"encoding/json"
)

var (
	MetaAttId    = "kanthor.att.id"
	MetaAttState = "kanthor.att.state"
	MetaEprId    = "kanthor.epr.id"
)

type Metadata map[string]string

func (meta Metadata) Get(key string) string {
	return meta[key]
}

func (meta Metadata) Set(key, value string) {
	meta[key] = value
}

func (meta Metadata) Merge(src map[string]string) {
	if len(src) > 0 {
		for key, value := range src {
			meta[key] = value
		}
	}
}

func (meta Metadata) String() string {
	if meta == nil {
		return ""
	}

	data, _ := json.Marshal(meta)
	return string(data)
}

// Scan implements the Scanner interface.
func (meta *Metadata) Scan(value interface{}) error {
	return json.Unmarshal([]byte(value.(string)), meta)
}

// Value implements the driver Valuer interface.
func (meta Metadata) Value() (driver.Value, error) {
	mt, err := json.Marshal(meta)
	return string(mt), err
}
