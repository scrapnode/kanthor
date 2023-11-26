package suid

import (
	"fmt"
	"strings"
	"time"

	"github.com/segmentio/ksuid"
)

var SafeUnixDiff = time.Second * 10

// New return a Sortable Unique IDentifier
// IMPORTANT: ksuid is replied on Unix timestamp so the datetime factor is only correct at second level
// IMPORTANT: that means the order of the ids with same unix timestamps are not guaranteed
// currently we are using KSUID (https://github.com/segmentio/ksuid)
// other standard UUID v7 (https://www.ietf.org/archive/id/draft-peabody-dispatch-new-uuid-format-01.html#name-uuidv7-layout-and-bit-order)
func New(ns string) string {
	return fmt.Sprintf("%s_%s", ns, ksuid.New().String())
}

func Parse(id string) (ns string, suid string) {
	if id == "" {
		return "", ""
	}
	segments := strings.Split(id, "_")
	if len(segments) == 0 {
		return "", ""
	}
	return segments[0], segments[1]
}

func Ns(id string) string {
	if id == "" {
		return ""
	}
	segments := strings.Split(id, "_")
	if len(segments) == 0 {
		return ""
	}
	return segments[0]
}

func Valid(id string) bool {
	ns, suid := Parse(id)
	if ns == "" || suid == "" {
		return false
	}
	if _, err := ksuid.Parse(suid); err != nil {
		return false
	}
	return true
}

// AfterTime uses SafeUnixDiff as factor to make sure we can get an id that is always less than the given time
func BeforeTime(t time.Time) string {
	id, err := ksuid.NewRandomWithTime(t.Add(SafeUnixDiff))
	if err != nil {
		panic(fmt.Sprintf("Couldn't generate KSUID, inconceivable! error: %v", err))
	}
	return id.String()
}

// AfterTime uses SafeUnixDiff as factor to make sure we can get an id that is always greater than the given time
func AfterTime(t time.Time) string {
	id, err := ksuid.NewRandomWithTime(t.Add(-SafeUnixDiff))
	if err != nil {
		panic(fmt.Sprintf("Couldn't generate KSUID, inconceivable! error: %v", err))
	}
	return id.String()
}
