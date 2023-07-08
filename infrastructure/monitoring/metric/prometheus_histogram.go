package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"sync"
)

type prometheush struct {
	conf    *Config
	entries map[string]prometheus.Histogram
	mu      sync.RWMutex
}

func (histogram *prometheush) Get(name string, withLabels []WithLabel) prometheus.Histogram {
	histogram.mu.Lock()
	defer histogram.mu.Unlock()

	labels := genLabels(withLabels)
	key := genKey(name, labels)

	if entry, ok := histogram.entries[key]; ok {
		return entry
	}

	histogram.entries[key] = promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace:   histogram.conf.Namespace,
		Name:        name,
		ConstLabels: labels,
		Buckets:     prometheus.DefBuckets,
	})
	return histogram.entries[name]
}
