package utils

import (
	"encoding/json"
	"github.com/google/uuid"
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

func RandomString(n int) string {
	var str string
	count := n / 32
	for i := 0; i <= count; i++ {
		str += strings.ReplaceAll(uuid.New().String(), "-", "")
	}

	return str[:n]
}
