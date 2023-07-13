package authenticator

import "context"

type ctxkey string

const (
	CtxAuthAccount    ctxkey = "kanthor.auth.account"
	CtxAuthWsIds      ctxkey = "kanthor.auth.workspace.ids"
	CtxAuthWsSelected ctxkey = "kanthor.auth.workspace.selected"
)

func AccountWithContext(ctx context.Context, account *Account) context.Context {
	return context.WithValue(ctx, CtxAuthAccount, account)
}

func AccountFromContext(ctx context.Context) *Account {
	return ctx.Value(CtxAuthAccount).(*Account)
}

func WorkspaceIdsWithContext(ctx context.Context, ids []string) context.Context {
	return context.WithValue(ctx, CtxAuthWsIds, ids)
}

func WorkspaceIdsFromContext(ctx context.Context) []string {
	return ctx.Value(CtxAuthWsIds).([]string)
}

func WorkspaceSelectedWithContext(ctx context.Context, ids []string) context.Context {
	return context.WithValue(ctx, CtxAuthWsSelected, ids)
}

func WorkspaceSelectedFromContext(ctx context.Context) []string {
	return ctx.Value(CtxAuthWsSelected).([]string)
}
