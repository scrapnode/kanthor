package cryptography_test

import (
	"github.com/scrapnode/kanthor/infrastructure/cryptography"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBcrypt_Success(t *testing.T) {
	kdf, err := cryptography.NewKDF(&cryptography.KDFConfig{})
	assert.Nil(t, err)

	password := "alonglongpasswordthatnobodyknow"

	hashed, err := kdf.StringHash(password)
	assert.Nil(t, err)

	err = kdf.StringCompare(hashed, password)
	assert.Nil(t, err)
}
