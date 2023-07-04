package metric

type Meter interface {
	Counter(name string, value int64, labels ...Label)
	Histogram(name string, value float64, labels ...Label)
}

type Label func(labels map[string]string)

func UseLabel(name, value string) Label {
	return func(labels map[string]string) {
		labels[name] = value
	}
}
