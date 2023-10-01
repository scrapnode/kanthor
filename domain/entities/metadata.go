package entities

import "encoding/json"

var (
	MetaIdempotencyKey = "kanthor.idempotency_key"
	MetaEprId          = "kanthor.epr.id"
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
	data, _ := json.Marshal(meta)
	return string(data)
}
