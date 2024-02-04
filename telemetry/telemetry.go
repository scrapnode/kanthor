package telemetry

import (
	"context"

	"github.com/scrapnode/kanthor/project"
	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type ctxkey string

var (
	CtxService ctxkey = "telemetry.ctx.service.name"
	CtxTracer  ctxkey = "telemetry.ctx.tracer"
)

// Start will help configure the telemetry in local machine
// On UAT/PROD we will use environment variable to configure it
func Start(ctx context.Context) error {
	uptrace.ConfigureOpentelemetry(
		// copy your project DSN here or use UPTRACE_DSN env var
		//uptrace.WithDSN("https://token@api.uptrace.dev/project_id"),

		uptrace.WithServiceName(ctx.Value(CtxService).(string)),
		uptrace.WithServiceVersion(project.GetVersion()),
		uptrace.WithDeploymentEnvironment(project.Env()),
	)
	return nil
}

func Stop(ctx context.Context) error {
	return uptrace.Shutdown(ctx)
}

func Tracer(name string) trace.Tracer {
	return otel.Tracer(name)
}
