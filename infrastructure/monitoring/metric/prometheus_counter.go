package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"sync"
)

type prometheusc struct {
	conf    *Config
	entries map[string]prometheus.Counter
	mu      sync.RWMutex
}

func (counter *prometheusc) Get(name string, withLabels []WithLabel) prometheus.Counter {
	counter.mu.Lock()
	defer counter.mu.Unlock()

	labels := genLabels(withLabels)
	key := genKey(name, labels)

	if entry, ok := counter.entries[key]; ok {
		return entry
	}

	counter.entries[key] = promauto.NewCounter(prometheus.CounterOpts{
		Namespace:   counter.conf.Namespace,
		Name:        name,
		ConstLabels: labels,
	})
	return counter.entries[key]
}
