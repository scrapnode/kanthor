package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
)

func NewAES(key string) *AES {
	return &AES{key: []byte(key)}
}

type AES struct {
	key []byte

	gcm cipher.AEAD
}

func (c *AES) Cipher() (cipher.AEAD, error) {
	if c.gcm != nil {
		return c.gcm, nil
	}

	// generate a new aes cipher using our 32 byte long key
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, err
	}

	// gcm or Galois/Counter Mode, is a mode of operation
	// for symmetric key cryptographic block ciphers
	// - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	c.gcm = gcm
	return c.gcm, nil
}

func (c *AES) Encrypt(plaintext []byte) ([]byte, error) {
	gcm, err := c.Cipher()
	if err != nil {
		return nil, err
	}

	// create a new byte array the size of the nonce which must be passed to Seal
	nonce := make([]byte, gcm.NonceSize())
	// populates our nonce with a cryptographically secure random sequence
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// here we encrypt our text using the Seal function
	// Seal encrypts and authenticates plaintext, authenticates the
	// additional data and appends the result to dst, returning the updated
	// slice. The nonce must be NonceSize() bytes long and unique for all
	// time, for a given key.
	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func (c *AES) EncryptString(plaintext string) (string, error) {
	bytes, err := c.Encrypt([]byte(plaintext))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func (c *AES) Decrypt(ciphertext []byte) ([]byte, error) {
	gcm, err := c.Cipher()
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("invalid ciphertext")
	}

	// since we know the ciphertext is actually nonce+ciphertext
	// and len(nonce) == NonceSize(). We can separate the two.
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func (c *AES) DecryptString(hextext string) (string, error) {
	ciphertext, err := hex.DecodeString(hextext)
	if err != nil {
		return "", err
	}

	bytes, err := c.Decrypt(ciphertext)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
