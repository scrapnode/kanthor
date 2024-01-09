package identifier

import (
	"fmt"
	"time"

	"github.com/segmentio/ksuid"
)

var SafeUnixDiff = time.Second * 10

// New return a Sortable Unique IDentifier
// IMPORTANT: ksuid is replied on Unix timestamp so the datetime factor is only correct at second level
// IMPORTANT: that means the order of the ids with same unix timestamps are not guaranteed
// currently we are using KSUID (https://github.com/segmentio/ksuid)
// other standard we can consider is UUID v7 (https://www.ietf.org/archive/id/draft-peabody-dispatch-new-uuid-format-01.html#name-uuidv7-layout-and-bit-order)
func New(ns string) string {
	return fmt.Sprintf("%s_%s", ns, ksuid.New().String())
}

func Id(ns, id string) string {
	return fmt.Sprintf("%s_%s", ns, id)
}
