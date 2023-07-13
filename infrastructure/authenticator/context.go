package authenticator

import "context"

type ctxkey string

const (
	CtxAuthAccount ctxkey = "kanthor.auth.account"
)

func AccountWithContext(ctx context.Context, account *Account) context.Context {
	return context.WithValue(ctx, CtxAuthAccount, account)
}

func AccountFromContext(ctx context.Context) *Account {
	return ctx.Value(CtxAuthAccount).(*Account)
}
