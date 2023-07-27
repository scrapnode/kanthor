package cryptography_test

import (
	"github.com/scrapnode/kanthor/infrastructure/cryptography"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAES_Success(t *testing.T) {
	key := "N1PCdw3M2B1TfJhoaY2mL736p2vCUc47"
	symmetric, err := cryptography.NewAES(&cryptography.SymmetricConfig{Key: key})
	assert.Nil(t, err)

	value := "i love kanthor so much"

	hextext, err := symmetric.StringEncrypt(value)
	assert.Nil(t, err)

	plaintext, err := symmetric.StringDecrypt(hextext)
	assert.Nil(t, err)

	assert.Equal(t, value, plaintext)
}
