package gateway

type ctxkey string

var (
	AccessPublicable  ctxkey = "kanthor.gateway.access.publicable"
	AccessProtectable ctxkey = "kanthor.gateway.access.protectable"
)
