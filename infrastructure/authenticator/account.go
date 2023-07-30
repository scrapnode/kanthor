package authenticator

import "github.com/scrapnode/kanthor/pkg/utils"

type Account struct {
	Sub string `json:"sub" yaml:"sub" validate:"required"`
	Iss string `json:"iss" yaml:"iss"`
	Aud string `json:"aud" yaml:"aud"`

	Name        string `json:"name" yaml:"name"`
	Picture     string `json:"picture" yaml:"picture"`
	Email       string `json:"email" yaml:"email"`
	PhoneNumber string `json:"phone" yaml:"phone_number"`
}

func (acc *Account) Modifier() string {
	return utils.Stringify(map[string]string{"id": acc.Sub, "name": acc.Name})
}
