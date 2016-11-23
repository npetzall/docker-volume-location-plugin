package profiler

import (
	"fmt"
	"time"
)

var p = false

func SetProfiling(value bool) {
	p = value
}

func Timed(start time.Time, name string) {
	if p {
		elapsed := time.Since(start)
		fmt.Printf("%s took %s\n", name, elapsed)
	}
}
