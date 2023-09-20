package validator

import "regexp"

var (
	DNSName   string = `^([a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62}){1}(\.[a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62})*[\._]?$`
	rxDNSName        = regexp.MustCompile(DNSName)
)

var (
	Alphanumeric   string = "^[a-zA-Z0-9]+$"
	rxAlphanumeric        = regexp.MustCompile(Alphanumeric)
)
