package project

import (
	"fmt"
	"os"
	"strings"
)

func RegionCode() string {
	code := os.Getenv("KANTHOR_REGION_CODE")
	if code != "" {
		return code
	}
	return "southeast"
}

func Namespace() string {
	ns := os.Getenv("KANTHOR_NAMESPACE")
	if ns != "" {
		return ns
	}
	return "kanthor"
}

func Tier() string {
	tier := os.Getenv("KANTHOR_TIER")
	if tier != "" {
		return tier
	}
	return "default"
}

func DefaultWorkspaceName() string {
	name := os.Getenv("KANTHOR_DEFAULT_WORKSPACE_NAME")
	if name != "" {
		return name
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

func Topic(segments ...string) string {
	return strings.Join(segments, ".")
}

func Subject(topic string) string {
	return fmt.Sprintf("%s.%s.%s", Namespace(), Tier(), topic)
}

func IsTopic(subject, topic string) bool {
	return strings.HasPrefix(subject, Subject(topic))
}
