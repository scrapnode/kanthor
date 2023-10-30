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
	return strings.EqualFold(os.Getenv("KANTHOR_ENV"), "development")
}
