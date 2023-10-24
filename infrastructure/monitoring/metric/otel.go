package metric

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func NewOtel(conf *Config, logger logging.Logger) (Metric, error) {
	logger = logger.With("monitoring.metrics", "otel")
	return &otel{conf: conf, logger: logger}, nil
}

type otel struct {
	conf     *Config
	logger   logging.Logger
	provider *sdkmetric.MeterProvider

	mu     sync.Mutex
	status int
}

// we don't need to take care about readiness & liveness for monitoring
// if they are failed, we will receive alerts from external system
func (metric *otel) Readiness() error {
	if metric.status == patterns.StatusDisconnected {
		return nil
	}
	if metric.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	return nil
}

func (metric *otel) Liveness() error {
	if metric.status == patterns.StatusDisconnected {
		return nil
	}
	if metric.status != patterns.StatusConnected {
		return ErrNotConnected
	}

	return nil
}

func (metric *otel) Connect(ctx context.Context) error {
	metric.mu.Lock()
	defer metric.mu.Unlock()

	if metric.status == patterns.StatusConnected {
		return ErrAlreadyConnected
	}

	exporter, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithEndpoint(metric.conf.Otel.Endpoint),
		otlpmetricgrpc.WithTimeout(time.Second*30),
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithRetry(otlpmetricgrpc.RetryConfig{
			Enabled:         true,
			InitialInterval: 5 * time.Second,
			MaxInterval:     30 * time.Second,
			MaxElapsedTime:  time.Minute,
		}),
	)
	if err != nil {
		return err
	}

	// labels/tags/resources that are common to all metrics.
	kv := []attribute.KeyValue{}
	if len(metric.conf.Otel.Labels) > 0 {
		for k, v := range metric.conf.Otel.Labels {
			kv = append(kv, attribute.String(k, v))
		}
	}
	rs := resource.NewWithAttributes(semconv.SchemaURL, kv...)

	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(rs),
		sdkmetric.WithReader(
			// collects and exports metric data every 30 seconds.
			sdkmetric.NewPeriodicReader(exporter, sdkmetric.WithInterval(time.Millisecond*time.Duration(metric.conf.Otel.Interval))),
		),
	)

	metric.status = patterns.StatusConnected
	metric.logger.Info("connected")
	metric.provider = provider
	return nil
}

func (metric *otel) Disconnect(ctx context.Context) error {
	metric.mu.Lock()
	defer metric.mu.Unlock()

	if metric.status != patterns.StatusConnected {
		return ErrNotConnected
	}
	metric.status = patterns.StatusDisconnected
	metric.logger.Info("disconnected")

	var returning error
	if err := metric.provider.Shutdown(ctx); err != nil {
		returning = errors.Join(returning, err)
	}

	return returning
}

func (metric *otel) Count(ctx context.Context, service, name string, value int64) {
	meter := metric.provider.Meter(metric.meter(service))
	if counter, err := meter.Int64Counter(name); err == nil {
		counter.Add(ctx, value)
	}
}

func (metric *otel) Observe(ctx context.Context, service, name string, value float64) {
	meter := metric.provider.Meter(metric.meter(service))
	if histogram, err := meter.Float64Histogram(name); err == nil {
		histogram.Record(ctx, value)
	}
}

func (metric *otel) meter(service string) string {
	return fmt.Sprintf("kanthorlabs/kanthor/%s", service)
}
