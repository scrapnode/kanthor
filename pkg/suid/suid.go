package suid

import (
	"fmt"
	"time"

	"github.com/segmentio/ksuid"
)

var SafeUnixDiff = time.Second * 7

// New return a Sortable Unique IDentifier
// IMPORTANT: ksuid is replied on Unix timestamp so the datetime factor is only correct at second level
// IMPORTANT: that means the order of the ids with same unix timestamps are not guaranteed
// currently we are using KSUID (https://github.com/segmentio/ksuid)
// other standard UUID v7 (https://www.ietf.org/archive/id/draft-peabody-dispatch-new-uuid-format-01.html#name-uuidv7-layout-and-bit-order)
func New(ns string) string {
	return fmt.Sprintf("%s_%s", ns, ksuid.New().String())
}

// AfterTime uses SafeUnixDiff as factor to make sure we can get an id that is always less than the given time
func BeforeTime(t time.Time) string {
	id, err := ksuid.NewRandomWithTime(t.Add(-SafeUnixDiff))
	if err != nil {
		panic(fmt.Sprintf("Couldn't generate KSUID, inconceivable! error: %v", err))
	}
	return id.String()
}

// AfterTime uses SafeUnixDiff as factor to make sure we can get an id that is always greater than the given time
func AfterTime(t time.Time) string {
	id, err := ksuid.NewRandomWithTime(t.Add(SafeUnixDiff))
	if err != nil {
		panic(fmt.Sprintf("Couldn't generate KSUID, inconceivable! error: %v", err))
	}
	return id.String()
}
