package entities

import "net/http"

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
