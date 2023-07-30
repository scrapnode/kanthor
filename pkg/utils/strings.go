package utils

import (
	"crypto/md5"
	"encoding/hex"
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

func MD5(values ...string) string {
	hash := md5.Sum([]byte(strings.Join(values, "/")))
	return hex.EncodeToString(hash[:])
}

func RandomString(n int) string {
	var str string
	count := n / 32
	for i := 0; i <= count; i++ {
		str += strings.ReplaceAll(uuid.New().String(), "-", "")
	}

	return str[:n]
}
