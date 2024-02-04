package project

import (
	_ "embed"
	"os"
	"strings"
)

var version string

func SetVersion(v string) {
	version = v
}

func GetVersion() string {
	return version
}

func IsDev() bool {
	return strings.EqualFold(Env(), "development")
}

func Env() string {
	if env := os.Getenv("KANTHOR_ENV"); env != "" {
		return env
	}
	return "production"
}
