package signature

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func Sign(key, msg string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(msg))

	return hex.EncodeToString(mac.Sum(nil))
}

func Verify(key, msg, hash string) (bool, error) {
	sig, err := hex.DecodeString(hash)
	if err != nil {
		return false, err
	}

	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(msg))

	return hmac.Equal(sig, mac.Sum(nil)), nil
}
