package gateway

type ctxkey string

var (
	KeyContext        = "kanthor.gateway.context"
	CtxWs      ctxkey = "kanthor.gateway.context.workspace"
	CtxAuhzOk  ctxkey = "kanthor.gateway.context.authz.ok"
)
