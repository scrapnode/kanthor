package debugging

import (
	"fmt"
	"time"
)

func Timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s ---> %v\n", name, time.Since(start))
	}
}
