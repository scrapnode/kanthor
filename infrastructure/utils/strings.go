package utils

import (
	"github.com/segmentio/ksuid"
	"strings"
)

func ID(ns string) string {
	return strings.Join([]string{ns, ksuid.New().String()}, "_")
}
