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
	ns := os.Getenv("KANTHOR_DEFAULT_TIER")
	if ns != "" {
		return ns
	}
	return "default"
}

func WorkspaceName() string {
	ns := os.Getenv("KANTHOR_DEFAULT_WORKSPACE_NAME")
	if ns != "" {
		return ns
	}
	return "main"
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

func SubjectInternal(topic string) string {
	return fmt.Sprintf("%s.internal.%s", Namespace(), topic)
}
