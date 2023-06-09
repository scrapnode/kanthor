package timer

import "time"

func New() Timer {
	return &timer{}
}

type Timer interface {
	Now() time.Time
}

type timer struct {
}

// Now return current UTC time
func (t *timer) Now() time.Time {
	return time.Now().UTC()
}
