package signature

import (
	corehmac "crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func NewHMAC() Signature {
	return &hmac{}
}

type hmac struct {
}

func (signature *hmac) Sign(msg, key string) string {
	mac := corehmac.New(sha256.New, []byte(key))
	mac.Write([]byte(msg))

	return hex.EncodeToString(mac.Sum(nil))
}

func (signature *hmac) Verify(msg, key, hash string) (bool, error) {
	sig, err := hex.DecodeString(hash)
	if err != nil {
		return false, err
	}

	mac := corehmac.New(sha256.New, []byte(key))
	mac.Write([]byte(msg))

	return corehmac.Equal(sig, mac.Sum(nil)), nil
}
