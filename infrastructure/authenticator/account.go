package authenticator

// Account is followed the document of https://www.iana.org/assignments/jwt/jwt.xhtml#claims
type Account struct {
	Sub      string            `json:"sub" yaml:"sub"`
	Name     string            `json:"name" yaml:"name"`
	Metadata map[string]string `json:"metadata" yaml:"metadata"`
}
