package telemetry

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Spanner struct {
	Tracer   trace.Tracer
	Contexts map[string]context.Context

	mu    sync.Mutex
	spans map[string]map[string]trace.Span
}

func (spanner *Spanner) Start(name string, kv ...attribute.KeyValue) {
	for refId := range spanner.Contexts {
		spanner.StartWithRefId(name, refId, kv...)
	}
}

func (spanner *Spanner) StartWithRefId(name, refId string, kv ...attribute.KeyValue) {
	spanner.mu.Lock()
	defer spanner.mu.Unlock()

	ctx, span := spanner.Tracer.Start(spanner.Contexts[refId], name, trace.WithAttributes(kv...))
	// override with new context that contains tracing id
	spanner.Contexts[refId] = ctx

	if spanner.spans == nil {
		spanner.spans = make(map[string]map[string]trace.Span)
	}
	if _, exist := spanner.spans[name]; !exist {
		spanner.spans[name] = make(map[string]trace.Span)
	}

	spanner.spans[name][refId] = span
}

func (spanner *Spanner) End(name string, kv ...attribute.KeyValue) {
	spanner.mu.Lock()
	defer spanner.mu.Unlock()

	if spans, ok := spanner.spans[name]; ok {
		for _, span := range spans {
			if len(kv) > 0 {
				span.SetAttributes(kv...)
			}
			span.End()
		}

		delete(spanner.spans, name)
	}
}

func (spanner *Spanner) Link(dest, src string) {
	spanner.Contexts[dest] = spanner.Contexts[src]
}
