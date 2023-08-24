package cryptography

func NewKDF(conf *KDFConfig) (KDF, error) {
	return NewBcrypt(conf)
}

// KDF is shorted of key derivation function
type KDF interface {
	Hash(value []byte) ([]byte, error)
	Compare(hash, value []byte) error

	StringHash(value string) (string, error)
	StringCompare(hash, value string) error
}
