package metric

import "fmt"

func New(conf *Config) (Meter, error) {
	if conf.Engine == EngineNoop {
		return NewNoop(conf)
	}

	return nil, fmt.Errorf("monitoring.metric: unknown engine")
}

type Meter interface {
	Count(name string, value int64, withLabels ...WithLabel)
	Histogram(name string, value float64, withLabels ...WithLabel)
}

type WithLabel func(labels map[string]string)

func Label(name, value string) WithLabel {
	return func(labels map[string]string) {
		labels[name] = value
	}
}
