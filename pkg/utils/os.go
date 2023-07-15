package utils

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

func AbsPathify(in string) string {
	if strings.HasPrefix(in, "$HOME") {
		in = UserHomeDir() + in[5:]
	}

	if strings.HasPrefix(in, "$") {
		end := strings.Index(in, string(os.PathSeparator))
		in = os.Getenv(in[1:end]) + in[end:]
	}

	if filepath.IsAbs(in) {
		return filepath.Clean(in)
	}

	p, err := filepath.Abs(in)
	if err == nil {
		return filepath.Clean(p)
	}

	return ""
}

func Liveness(name string) error {
	fp := path.Join(os.TempDir(), fmt.Sprintf("%s.liveness", name))
	return os.WriteFile(fp, []byte("live"), os.ModePerm)
}

func Readiness(name string) error {
	fp := path.Join(os.TempDir(), fmt.Sprintf("%s.readiness", name))
	return os.WriteFile(fp, []byte("live"), os.ModePerm)
}
