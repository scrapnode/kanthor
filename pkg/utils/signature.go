package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func SignatureSign(key, data string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(data))

	return hex.EncodeToString(mac.Sum(nil))
}
