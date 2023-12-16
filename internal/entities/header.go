package entities

import (
	"encoding/json"
	"net/http"
	"net/textproto"
)

var (
	HeaderIdempotencyKey = "idempotency-key"
	HeaderWebhookId      = "webhook-id"
	HeaderWebhookTs      = "webhook-timestamp"
	HeaderWebhookSig     = "webhook-signature"
	HeaderWebhookRef     = "webhook-ref"
)

type Header map[string][]string

func (h Header) Add(key, value string) {
	key = textproto.CanonicalMIMEHeaderKey(key)
	h[key] = append(h[key], value)
}

func (h Header) Set(key, value string) {
	h[textproto.CanonicalMIMEHeaderKey(key)] = []string{value}
}

func (h Header) Get(key string) string {
	if h == nil {
		return ""
	}
	v := h[textproto.CanonicalMIMEHeaderKey(key)]
	if len(v) == 0 {
		return ""
	}
	return v[0]
}

func (h Header) Values(key string) []string {
	if h == nil {
		return nil
	}
	return h[textproto.CanonicalMIMEHeaderKey(key)]
}

func (h Header) Del(key string) {
	delete(h, textproto.CanonicalMIMEHeaderKey(key))
}

func (h Header) Merge(src Header) {
	if len(src) > 0 {
		for key, values := range src {
			for _, v := range values {
				h.Add(key, v)
			}
		}
	}
}

func (h Header) ToHTTP() http.Header {
	returning := http.Header{}
	for key, values := range h {
		for _, value := range values {
			returning.Add(key, value)
		}
	}
	return returning
}

func (h Header) FromHTTP(headers http.Header) {
	for key, values := range headers {
		for _, value := range values {
			h.Add(key, value)
		}
	}
}

func (h Header) String() string {
	data, _ := json.Marshal(h)
	return string(data)
}
