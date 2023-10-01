package entities

import (
	"encoding/json"
	"net/http"
)

func NewHeader() Header {
	return Header{Header: http.Header{}}
}

type Header struct {
	http.Header
}

func (h Header) Merge(src Header) {
	if h.Header == nil {
		h.Header = http.Header{}
	}

	if len(src.Header) > 0 {
		for key, values := range src.Header {
			for _, v := range values {
				h.Add(key, v)
			}
		}
	}
}

func (h Header) String() string {
	data, _ := json.Marshal(h.Header)
	return string(data)
}
