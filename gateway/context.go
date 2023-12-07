package gateway

var Ctx = "kanthor.gateway.context"

type ctxkey string

const (
	CtxAccount   ctxkey = "kanthor.gateway.context.account"
	CtxWorkspace ctxkey = "kanthor.gateway.context.workspace"
)
