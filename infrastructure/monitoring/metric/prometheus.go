package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/scrapnode/kanthor/infrastructure/logging"
	"github.com/scrapnode/kanthor/infrastructure/patterns"
)

func NewPrometheus(conf *Config) Meter {
	return &prometheusio{
		conf:       conf,
		counters:   &prometheusc{conf: conf, entries: map[string]prometheus.Counter{}},
		histograms: &prometheush{conf: conf, entries: map[string]prometheus.Histogram{}},
	}
}

func NewPrometheusExporter(conf *Config, logger logging.Logger) patterns.Runnable {
	return NewHttpServer(conf, logger, promhttp.Handler())
}

type prometheusio struct {
	conf       *Config
	counters   *prometheusc
	histograms *prometheush
}

func (metric *prometheusio) Count(name string, value int64, withLabels ...WithLabel) {
	counter := metric.counters.Get(name, withLabels)
	counter.Add(float64(value))
}

func (metric *prometheusio) Histogram(name string, value float64, withLabels ...WithLabel) {
	histogram := metric.histograms.Get(name, withLabels)
	histogram.Observe(value)
}
