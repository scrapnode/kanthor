package entities

import (
	"encoding/json"
	"net/http"
)

type Header struct {
	http.Header
}

func (h Header) Merge(src Header) {
	for key, values := range src.Header {
		for _, v := range values {
			h.Add(key, v)
		}
	}
}

func (h Header) String() string {
	data, _ := json.Marshal(h.Header)
	return string(data)
}
