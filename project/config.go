package project

import (
	_ "embed"
	"os"
	"strings"
)

//go:embed .version
var version string

func Version() string {
	return version
}

func IsDev() bool {
	return strings.EqualFold(os.Getenv("KANTHOR_ENV"), "development")
}
