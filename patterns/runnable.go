package patterns

import "context"

type Runnable interface {
	Stop(ctx context.Context) error
	Start(ctx context.Context) error
	Run(ctx context.Context) error
}

type CommandLine interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Usecase() any
}
