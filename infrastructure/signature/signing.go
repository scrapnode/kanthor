package signature

func New() Signature {
	return NewHMAC()
}

type Signature interface {
	Sign(msg, key string) string
	Verify(msg, key, hash string) (bool, error)
}
