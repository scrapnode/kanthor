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

func StringStartsWithIfNotEmpty(prop, value, prefix string) Fn {
	v := strings.Trim(value, " ")

	return func() error {
		if len(v) == 0 {
			return nil
		}

		if !strings.HasPrefix(v, prefix) {
			return fmt.Errorf("%s (value:%s) must be started with %s", prop, value, prefix)
		}
		return nil
	}
}

func StringStartsWith(prop, value, prefix string) Fn {
	return func() error {
		if err := StringRequired(prop, value)(); err != nil {
			return err
		}
		if err := StringStartsWithIfNotEmpty(prop, value, prefix)(); err != nil {
			return err
		}
		return nil
	}
}

func StringStartsWithOneOf(prop, value string, prefixes []string) Fn {
	return func() error {
		for i := range prefixes {
			if err := StringStartsWithIfNotEmpty(prop, value, prefixes[i])(); err == nil {
				return nil
			}
		}
		return fmt.Errorf("%s (value:%s) prefix must be started with one of %s", prop, value, prefixes)
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

func StringLenIfNotEmpty(prop, value string, min, max int) Fn {
	return func() error {
		if len(value) == 0 {
			return nil
		}
		return StringLen(prop, value, min, max)()
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

		if _, has := m[value]; !has {
			return fmt.Errorf("%s (value:%s) must be one of %q", prop, value, oneOf)
		}

		return nil
	}
}

func StringIn(prop string, values []string, in []string) Fn {
	m := map[string]bool{}
	for _, o := range in {
		m[o] = true
	}

	return func() error {
		for _, value := range values {
			if err := StringRequired(prop, value)(); err != nil {
				return err
			}

			if _, has := m[value]; has {
				return nil
			}

			if _, has := m[value]; !has {
				return fmt.Errorf("%s contains unrecognized %q (expected:%v)", prop, value, in)
			}
		}

		return nil
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
