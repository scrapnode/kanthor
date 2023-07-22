package dataplane

type ctxkey string

const (
	CtxWorkspace   ctxkey = "kanthor.usecase.dataplane.workspace"
	CtxApplication ctxkey = "kanthor.usecase.dataplane.application"
)
