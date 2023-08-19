package metrics

import (
	prometheuscore "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"net/http"
	"sync"
)

func NewPrometheus(conf *Config, logger logging.Logger) (Metrics, error) {
	logger = logger.With("monitoring.metrics", "prometheus")
	registry := prometheuscore.NewRegistry()
	m := &prometheus{
		conf:       conf,
		logger:     logger,
		registry:   registry,
		counters:   map[string]prometheuscore.Counter{},
		histograms: map[string]prometheuscore.Histogram{},
	}

	return m, nil
}

type prometheus struct {
	conf     *Config
	logger   logging.Logger
	registry *prometheuscore.Registry
	cmu      sync.RWMutex

	counters map[string]prometheuscore.Counter

	hmu        sync.RWMutex
	histograms map[string]prometheuscore.Histogram
}

func (metrics *prometheus) Count(name string, value int64) {
	metrics.cmu.Lock()
	defer metrics.cmu.Unlock()

	if counter, ok := metrics.counters[name]; ok {
		counter.Add(float64(value))
		return
	}

	counter := prometheuscore.NewCounter(prometheuscore.CounterOpts{
		Name:        name,
		ConstLabels: metrics.conf.Prometheus.Labels,
	})
	if err := metrics.registry.Register(counter); err != nil {
		metrics.logger.Errorw(err.Error(), "counter_name", name)
		return
	}
	counter.Add(float64(value))
	metrics.counters[name] = counter
}

func (metrics *prometheus) Observe(name string, value float64) {
	metrics.hmu.Lock()
	defer metrics.hmu.Unlock()

	if histogram, ok := metrics.histograms[name]; ok {
		histogram.Observe(value)
		return
	}

	histogram := prometheuscore.NewHistogram(prometheuscore.HistogramOpts{
		Name:        name,
		ConstLabels: metrics.conf.Prometheus.Labels,
		Buckets:     []float64{.1, .25, .5, 1, 2.5, 5, 10},
	})
	if err := metrics.registry.Register(histogram); err != nil {
		metrics.logger.Errorw(err.Error(), "histogram_name", name)
		return
	}
	histogram.Observe(value)
	metrics.histograms[name] = histogram
}

func (metrics *prometheus) Handler() http.Handler {
	return promhttp.HandlerFor(metrics.registry, promhttp.HandlerOpts{Registry: metrics.registry})
}
