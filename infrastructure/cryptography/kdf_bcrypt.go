package cryptography

import (
	bcryptcore "golang.org/x/crypto/bcrypt"
)

func NewBcrypt(conf *KDFConfig) (KDF, error) {
	return &bcrypt{}, nil
}

type bcrypt struct {
}

func (kdf *bcrypt) Hash(value []byte) ([]byte, error) {
	return bcryptcore.GenerateFromPassword(value, bcryptcore.DefaultCost)
}

func (kdf *bcrypt) Compare(hashed, value []byte) error {
	return bcryptcore.CompareHashAndPassword(hashed, value)
}

func (kdf *bcrypt) StringHash(value string) (string, error) {
	hashed, err := kdf.Hash([]byte(value))
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func (kdf *bcrypt) StringCompare(hashed, value string) error {
	return kdf.Compare([]byte(hashed), []byte(value))
}
