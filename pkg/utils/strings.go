package utils

import (
	"encoding/json"
	"errors"
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

func UrlScheme(rawURL string) (scheme string, err error) {
	for i := 0; i < len(rawURL); i++ {
		c := rawURL[i]
		switch {
		case 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z':
		// do nothing
		case '0' <= c && c <= '9' || c == '+' || c == '-' || c == '.':
			if i == 0 {
				return "", nil
			}
		case c == ':':
			if i == 0 {
				return "", errors.New("missing protocol scheme")
			}
			return rawURL[:i], nil
		default:
			// we have encountered an invalid character,
			// so there is no valid scheme
			return "", nil
		}
	}
	return "", nil
}
