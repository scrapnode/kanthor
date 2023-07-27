package cryptography

func New(conf *Config) (Cryptography, error) {
	kdf, err := NewKDF(&conf.KDF)
	if err != nil {
		return nil, err
	}

	symmetric, err := NewSymmetric(&conf.Symmetric)
	if err != nil {
		return nil, err
	}

	return &cryptology{kdf: kdf, symmetric: symmetric}, nil
}

type Cryptography interface {
	KDF() KDF
	Symmetric() Symmetric
}

// cryptology is not a good name, but I have no idea what name should I use
type cryptology struct {
	kdf       KDF
	symmetric Symmetric
}

func (c *cryptology) KDF() KDF {
	return c.kdf
}

func (c *cryptology) Symmetric() Symmetric {
	return c.symmetric
}
