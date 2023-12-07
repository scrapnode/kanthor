package authenticator

import "strings"

func Parse(authorization string) (string, string) {
	segments := strings.Split(authorization, " ")
	if len(segments) != 2 {
		return "", ""
	}
	return segments[0], segments[1]
}
