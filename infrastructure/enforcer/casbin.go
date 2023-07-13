package enforcer

import (
	gocasbin "github.com/casbin/casbin/v2"
	"net/url"
)

func NewCasbin(conf *Config) Enforcer {
	modelUrl, err := url.Parse(conf.Casbin.ModelSource)
	if err != nil {
		panic(err)
	}

	policyUrl, err := url.Parse(conf.Casbin.PolicySource)
	if err != nil {
		panic(err)
	}

	client, err := gocasbin.NewEnforcer(modelUrl.RawPath, policyUrl.RawPath)
	if err != nil {
		panic(err)
	}

	return &casbin{client: client}
}

type casbin struct {
	client *gocasbin.Enforcer
}

func (enforcer *casbin) Enforce(sub, obj, act string) (bool, error) {
	return enforcer.client.Enforce(sub, obj, act)
}
