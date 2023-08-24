package cryptography

import (
	bcryptcore "golang.org/x/crypto/bcrypt"
)

func NewBcrypt(conf *KDFConfig) (KDF, error) {
	return &bcrypt{conf: conf}, nil
}

type bcrypt struct {
	conf *KDFConfig
}

func (kdf *bcrypt) Hash(value []byte) ([]byte, error) {
	return bcryptcore.GenerateFromPassword(value, bcryptcore.DefaultCost)
}

func (kdf *bcrypt) Compare(hash, value []byte) error {
	return bcryptcore.CompareHashAndPassword(hash, value)
}

func (kdf *bcrypt) StringHash(value string) (string, error) {
	hash, err := kdf.Hash([]byte(value))
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (kdf *bcrypt) StringCompare(hash, value string) error {
	return kdf.Compare([]byte(hash), []byte(value))
}
