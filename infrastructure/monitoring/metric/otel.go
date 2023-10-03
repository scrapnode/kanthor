package metric

import (
	"context"
	"sync"
	"time"

	"github.com/scrapnode/kanthor/infrastructure/logging"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func NewOtel(conf *Config, logger logging.Logger) (Metrics, error) {
	logger = logger.With("monitoring.metrics", "otel")
	return &otel{conf: conf, logger: logger}, nil
}

type otel struct {
	conf     *Config
	logger   logging.Logger
	provider *sdkmetric.MeterProvider

	mu sync.Mutex
}

// we don't need to take care about readiness & liveness for monitoring
// if they are failed, we will receive alerts from external system
func (metric *otel) Readiness() error {
	return nil
}

func (metric *otel) Liveness() error {
	return nil
}

func (metric *otel) Connect(ctx context.Context) error {
	metric.mu.Lock()
	defer metric.mu.Unlock()

	if metric.provider != nil {
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
	kv := []attribute.KeyValue{semconv.ServiceNameKey.String(metric.conf.Otel.Service)}
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

	metric.logger.Info("connected")
	metric.provider = provider
	return nil
}

func (metric *otel) Disconnect(ctx context.Context) error {
	metric.mu.Lock()
	defer metric.mu.Unlock()

	metric.logger.Info("disconnected")

	if metric.provider != nil {
		if err := metric.provider.Shutdown(ctx); err != nil {
			metric.logger.Error(err)
		}
	}

	return nil
}

func (metric *otel) Count(ctx context.Context, name string, value int64) {
	meter := metric.provider.Meter("github.com/scrapnode/kanthor")
	if counter, err := meter.Int64Counter(name); err == nil {
		counter.Add(context.Background(), value)
	}
}

func (metric *otel) Observe(ctx context.Context, name string, value float64) {
	meter := metric.provider.Meter("github.com/scrapnode/kanthor")
	if histogram, err := meter.Float64Histogram(name); err == nil {
		histogram.Record(context.Background(), value)
	}
}
