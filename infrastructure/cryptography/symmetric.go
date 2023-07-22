package cryptography

func NewSymmetric(conf *SymmetricConfig) (Symmetric, error) {
	return NewAES(conf)
}

type Symmetric interface {
	Encrypt(plaintext []byte) ([]byte, error)
	Decrypt(ciphertext []byte) ([]byte, error)

	StringEncrypt(plaintext string) (string, error)
	StringDecrypt(hextext string) (string, error)
}
