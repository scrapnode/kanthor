package metric

import "strings"

func genLabels(withLabels []WithLabel) map[string]string {
	kv := map[string]string{}
	if len(withLabels) > 0 {
		for _, withLabel := range withLabels {
			withLabel(kv)
		}
	}
	return kv
}

func genKey(name string, labels map[string]string) string {
	segments := []string{name}
	for k, v := range labels {
		segments = append(segments, k, v)
	}
	return strings.Join(segments, "/")
}
