package validator

import (
	"net"
	"strconv"
	"strings"
)

// IsAlphanumeric checks if the string contains only letters and numbers. Empty string is invalid.
func IsAlphanumeric(str string) bool {
	if len(str) == 0 {
		return false
	}
	return rxAlphanumeric.MatchString(str)
}

// IsIP checks if a string is either IP version 4 or 6. Alias for `net.ParseIP`
func IsIP(str string) bool {
	return net.ParseIP(str) != nil
}

// IsDNSName will validate the given string as a DNS name
func IsDNSName(str string) bool {
	if str == "" || len(strings.Replace(str, ".", "", -1)) > 255 {
		// constraints already violated
		return false
	}
	return !IsIP(str) && rxDNSName.MatchString(str)
}

// IsPort checks if a string represents a valid port
func IsPort(str string) bool {
	if i, err := strconv.Atoi(str); err == nil && i > 0 && i < 65536 {
		return true
	}
	return false
}

// IsHostPort validates the given string for usage with the various Dial() functions
func IsHostPort(str string) bool {
	if h, p, err := net.SplitHostPort(str); err == nil && h != "" && p != "" && (IsDNSName(h) || IsIP(h)) && IsPort(p) {
		return true
	}

	return false
}
