package sdk

type ctxkey string

const (
	CtxAcc ctxkey = "usecases.sdk.account"
	CtxWs  ctxkey = "usecases.sdk.workspace"
	CtxWst ctxkey = "usecases.sdk.workspace.tier"
)
