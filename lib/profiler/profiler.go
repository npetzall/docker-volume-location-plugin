package profiler

import (
	"fmt"
	"time"
)

var p = false

//SetProfiling controls if profiling information is written to stdout, defaults to false.
func SetProfiling(value bool) {
	p = value
}

//Timed is a function to print profiling information to stdout.
//defer Timer(time.Now(), [Name of the thing you'r timing])
func Timed(start time.Time, name string) {
	if p {
		elapsed := time.Since(start)
		fmt.Printf("%s took %s\n", name, elapsed)
	}
}
