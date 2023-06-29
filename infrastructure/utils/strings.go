package utils

import (
	"encoding/json"
	"github.com/segmentio/ksuid"
	"strings"
)

func ID(ns string) string {
	return strings.Join([]string{ns, ksuid.New().String()}, "_")
}

func Key(values ...string) string {
	return strings.Join(values, "/")
}

func Stringify(value interface{}) string {
	bytes, _ := json.Marshal(value)
	return string(bytes)
}
