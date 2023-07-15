package authenticator

type Account struct {
	Sub string `json:"sub" yaml:"sub"`
	Iss string `json:"iss" yaml:"iss"`
	Aud string `json:"aud" yaml:"aud"`

	Name    string `json:"name" yaml:"name"`
	Email   string `json:"email" yaml:"email"`
	Picture string `json:"picture" yaml:"picture"`
	Address string `json:"address" yaml:"address"`
	Phone   string `json:"phone" yaml:"phone"`
}
