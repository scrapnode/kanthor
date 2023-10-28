package project

import (
	_ "embed"
	"os"
	"strings"
)

func IsDev() bool {
	return strings.EqualFold(os.Getenv("KANTHOR_ENV"), "development")
}
