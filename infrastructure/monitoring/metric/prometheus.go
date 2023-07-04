package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
	"sync"
)

func NewPrometheus(conf *Config) Meter {
	return &prometheusio{}
}

func NewPrometheusExporter(conf *Config, logger logging.Logger) patterns.Runnable {
	return NewHttpServer(conf, logger, promhttp.Handler())
}

type prometheusio struct {
	counters   sync.Map
	histograms sync.Map
}

func (metric *prometheusio) Counter(name string, value int64, labels ...Label) {
	counter := metric.counter(name, metric.labels(labels))
	counter.Add(float64(value))
}

func (metric *prometheusio) counter(name string, labels map[string]string) prometheus.Counter {
	if value, ok := metric.counters.Load(name); ok {
		return value.(prometheus.Counter)
	}

	counter := promauto.NewCounter(prometheus.CounterOpts{
		Name:        name,
		ConstLabels: labels,
	})
	metric.counters.Store(name, counter)
	return counter
}

func (metric *prometheusio) Histogram(name string, value float64, labels ...Label) {
	histogram := metric.histogram(name, metric.labels(labels))
	histogram.Observe(value)
}

func (metric *prometheusio) histogram(name string, labels map[string]string) prometheus.Histogram {
	if value, ok := metric.histograms.Load(name); ok {
		return value.(prometheus.Histogram)
	}

	histogram := promauto.NewHistogram(prometheus.HistogramOpts{
		Name:        name,
		ConstLabels: labels,
		Buckets:     prometheus.DefBuckets,
	})
	metric.histograms.Store(name, histogram)
	return histogram
}

func (metric *prometheusio) labels(labels []Label) map[string]string {
	kv := map[string]string{}
	if len(labels) > 0 {
		for _, l := range labels {
			l(kv)
		}
	}
	return kv
}
