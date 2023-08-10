package authorizator

type ctxkey string

const (
	CtxWs  ctxkey = "usecases.sdk.workspace"
	CtxWst ctxkey = "usecases.sdk.workspace.tier"
)
