package metric

func NewNoop(conf *Config) (Meter, error) {
	return &noop{}, nil
}

type noop struct {
}

func (metric *noop) Count(name string, value int64, withLabels ...WithLabel) {}

func (metric *noop) Histogram(name string, value float64, withLabels ...WithLabel) {}
