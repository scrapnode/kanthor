package validator

import (
	"fmt"
	"net/url"
	"strings"
)

func StringRequired(prop, value string) Fn {
	return func() error {
		if strings.Trim(value, " ") == "" {
			return fmt.Errorf("%s is required", prop)
		}
		return nil
	}
}

func StringStartsWith(prop, value, prefix string) Fn {
	return func() error {
		if err := StringRequired(prop, value)(); err != nil {
			return err
		}
		if !strings.HasPrefix(strings.Trim(value, " "), prefix) {
			return fmt.Errorf("%s (value:%s) is must be started with %s", prop, prefix, value)
		}
		return nil
	}
}

func StringUri(prop, value string) Fn {
	return func() error {
		if err := StringRequired(prop, value)(); err != nil {
			return err
		}
		if _, err := url.ParseRequestURI(value); err != nil {
			return fmt.Errorf("%s (error:%s) is not a valid uri ", prop, err.Error())
		}
		return nil
	}
}

func StringLen(prop, value string, min, max int) Fn {
	return func() error {
		if len(value) < min {
			return fmt.Errorf("%s (len:%d) length must be greater than or equal %d", prop, len(value), min)
		}
		if len(value) > max {
			return fmt.Errorf("%s (len:%d) length must be mess than or equal %d", prop, len(value), max)
		}
		return nil
	}
}

func StringOneOf(prop, value string, oneOf []string) Fn {
	m := map[string]bool{}
	for _, o := range oneOf {
		m[o] = true
	}

	return func() error {
		if err := StringRequired(prop, value)(); err != nil {
			return err
		}

		for _, o := range oneOf {
			if b, has := m[o]; has && b {
				return nil
			}
		}

		return fmt.Errorf("%s (value:%s) must be one of %q", prop, value, oneOf)
	}
}

func StringHostPort(prop, value string) Fn {
	return func() error {
		if !IsHostPort(value) {
			return fmt.Errorf("%s (value:%s) is not a valid host:port string", prop, value)
		}

		return nil
	}
}

func StringAlphanumeric(prop, value string) Fn {
	return func() error {
		if !IsAlphanumeric(value) {
			return fmt.Errorf("%s (value:%s) is not matched %s", prop, value, Alphanumeric)
		}

		return nil
	}
}
