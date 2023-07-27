package cryptography

import (
	aescore "crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
)

func NewAES(conf *SymmetricConfig) (Symmetric, error) {
	// generate a new aes cipher using our 32 byte long key
	block, err := aescore.NewCipher([]byte(conf.Key))
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

	return &aes{gcm: gcm}, nil
}

type aes struct {
	gcm cipher.AEAD
}

func (symmetric *aes) Encrypt(plaintext []byte) ([]byte, error) {
	// create a new byte array the size of the nonce which must be passed to Seal
	nonce := make([]byte, symmetric.gcm.NonceSize())
	// populates our nonce with a cryptographically secure random sequence
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// here we encrypt our text using the Seal function
	// Seal encrypts and authenticates plaintext, authenticates the
	// additional data and appends the result to dst, returning the updated
	// slice. The nonce must be NonceSize() bytes long and unique for all
	// time, for a given key.
	return symmetric.gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func (symmetric *aes) Decrypt(ciphertext []byte) ([]byte, error) {
	nonceSize := symmetric.gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("invalid ciphertext")
	}

	// since we know the ciphertext is actually nonce+ciphertext
	// and len(nonce) == NonceSize(). We can separate the two.
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return symmetric.gcm.Open(nil, nonce, ciphertext, nil)
}

func (symmetric *aes) StringEncrypt(plaintext string) (string, error) {
	bytes, err := symmetric.Encrypt([]byte(plaintext))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func (symmetric *aes) StringDecrypt(hextext string) (string, error) {
	ciphertext, err := hex.DecodeString(hextext)
	if err != nil {
		return "", err
	}

	bytes, err := symmetric.Decrypt(ciphertext)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
