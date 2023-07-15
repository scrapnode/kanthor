package authorizator

func New(conf *Config) Authorizator {
	return NewCasbin(conf)
}

type Authorizator interface {
	Enforce(sub, dom, obj, act string) (bool, error)
}
