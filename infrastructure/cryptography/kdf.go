package cryptography

func NewKDF(conf *KDFConfig) (KDF, error) {
	return NewBcrypt(conf)
}

// KDF is shorted of key derivation function
type KDF interface {
	Hash(value []byte) ([]byte, error)
	Compare(hashed, value []byte) error

	StringHash(value string) (string, error)
	StringCompare(hashed, value string) error
}
