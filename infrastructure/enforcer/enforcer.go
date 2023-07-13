package enforcer

func New(conf *Config) Enforcer {
	return NewCasbin(conf)
}

type Enforcer interface {
	Enforce(sub, obj, act string) (bool, error)
}
