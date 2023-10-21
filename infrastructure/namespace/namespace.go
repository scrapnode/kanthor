package namespace

import (
	"fmt"
	"os"
)

func Namespace() string {
	ns := os.Getenv("KANTHOR_NAMESPACE")
	if ns != "" {
		return ns
	}
	return "kanthor"
}

func Tier() string {
	ns := os.Getenv("KANTHOR_TIER")
	if ns != "" {
		return ns
	}
	return "default"
}

func Name(name string) string {
	return fmt.Sprintf("%s_%s_%s", Namespace(), Tier(), name)
}

func NameWithoutTier(name string) string {
	return fmt.Sprintf("%s_%s", Namespace(), name)
}

func Key(key string) string {
	return fmt.Sprintf("%s/%s/%s", Namespace(), Tier(), key)
}

func Subject(topic string) string {
	return fmt.Sprintf("%s.%s.%s", Namespace(), Tier(), topic)
}
