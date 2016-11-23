package profiler

import (
	"fmt"
	"time"
)

var p = false

// Enable profiling of api executions, defaults to false
func SetProfiling(value bool) {
	p = value
}

//Used to print profiling information to stdout.
//defer Timer(time.Now(), [Name of the thing you'r timing])
func Timed(start time.Time, name string) {
	if p {
		elapsed := time.Since(start)
		fmt.Printf("%s took %s\n", name, elapsed)
	}
}
