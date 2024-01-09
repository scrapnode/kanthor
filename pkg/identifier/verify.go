package identifier

import (
	"strings"

	"github.com/segmentio/ksuid"
)

func Parse(id string) (ns string, suid string) {
	if id == "" {
		return "", ""
	}
	segments := strings.Split(id, "_")
	if len(segments) == 0 {
		return "", ""
	}
	return segments[0], segments[1]
}

func Ns(id string) string {
	if id == "" {
		return ""
	}
	segments := strings.Split(id, "_")
	if len(segments) == 0 {
		return ""
	}
	return segments[0]
}

func Valid(id string) bool {
	ns, suid := Parse(id)
	if ns == "" || suid == "" {
		return false
	}
	if _, err := ksuid.Parse(suid); err != nil {
		return false
	}
	return true
}
