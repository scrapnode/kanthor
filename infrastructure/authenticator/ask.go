package authenticator

import "context"

var EngineAsk = "ask"

func NewAsk(conf *Ask) Verifier {
	return &ask{conf: conf}
}

type ask struct {
	conf *Ask
}

func (verifier *ask) Verify(ctx context.Context, request *Request) (*Account, error) {
	user, pass, err := ParseBasicCredentials(request.Credentials)
	if err != nil {
		return nil, err
	}

	accessOK := user == verifier.conf.AccessKey
	secretOk := pass == verifier.conf.SecretKey
	if !accessOK || !secretOk {
		return nil, ErrInvalidCredentials
	}

	return &Account{Sub: user, Name: user, Metadata: make(map[string]string, 0)}, nil
}
