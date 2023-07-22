package crypto_test

import (
	"encoding/hex"
	"fmt"
	"github.com/scrapnode/kanthor/infrastructure/crypto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAES_Success(t *testing.T) {
	key := "N1PCdw3M2B1TfJhoaY2mL736p2vCUc47"
	aes := crypto.NewAES(key)

	value := "i love kanthor so much"

	ciphertext, err := aes.Encrypt([]byte(value))
	assert.Nil(t, err)

	s := hex.EncodeToString(ciphertext)
	fmt.Sprintln(s)

	plaintext, err := aes.Decrypt(ciphertext)
	assert.Nil(t, err)

	assert.Equal(t, []byte(value), plaintext)
}
