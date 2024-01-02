package utils

import (
	"encoding/json"
	"strings"

	"github.com/google/uuid"
)

func Key(values ...string) string {
	return strings.Join(values, "/")
}

func Stringify(value interface{}) string {
	bytes, _ := json.Marshal(value)
	return string(bytes)
}

func StringifyIndent(value interface{}) string {
	bytes, _ := json.MarshalIndent(value, "", "  ")
	return string(bytes)
}

func RandomString(n int) string {
	var str string
	count := n / 32
	for i := 0; i <= count; i++ {
		str += strings.ReplaceAll(uuid.NewString(), "-", "")
	}

	return str[:n]
}
