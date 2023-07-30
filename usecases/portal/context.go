package portal

type ctxkey string

const (
	CtxAcc ctxkey = "usecases.portal.account"
	CtxWs  ctxkey = "usecases.portal.workspace"
	CtxWst ctxkey = "usecases.portal.workspace.tier"
)
