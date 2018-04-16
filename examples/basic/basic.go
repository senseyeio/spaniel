package main

import (
	"fmt"
	timespan "github.com/senseyeio/spaniel"
	"time"
)

func main() {

	// Times at half-hourly intervals
	var t1 = time.Date(2018, 1, 30, 0, 0, 0, 0, time.UTC)
	var t2 = time.Date(2018, 1, 30, 0, 30, 0, 0, time.UTC)
	var t3 = time.Date(2018, 1, 30, 1, 0, 0, 0, time.UTC)
	var t4 = time.Date(2018, 1, 30, 1, 30, 0, 0, time.UTC)

	input := timespan.List{
		timespan.New(t1, t3),
		timespan.New(t2, t4),
	}

	union := input.Union()
	fmt.Println(union[0].Start(), "->", union[0].End())

	intersection := input.Intersection()
	fmt.Println(intersection[0].Start(), "->", intersection[0].End())
}
