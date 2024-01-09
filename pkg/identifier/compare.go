package identifier

import (
	"fmt"
	"time"

	"github.com/segmentio/ksuid"
)

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
	id, err := ksuid.NewRandomWithTime(t.Add(+SafeUnixDiff))
	if err != nil {
		panic(fmt.Sprintf("Couldn't generate KSUID, inconceivable! error: %v", err))
	}
	return id.String()
}
