package spaniel_test

import (
	"fmt"
	"github.com/senseyeio/spaniel"
	"sort"
	"time"
)

type dur struct {
	from time.Time
	to   time.Time
}

var times []dur
var input spaniel.Spans

func init() {
	times = []dur{{from: time.Date(2018, 1, 30, 0, 0, 0, 0, time.UTC), to: time.Date(2018, 1, 30, 1, 0, 0, 0, time.UTC)},
		{from: time.Date(2018, 1, 30, 0, 30, 0, 0, time.UTC), to: time.Date(2018, 1, 30, 1, 30, 0, 0, time.UTC)},
		{from: time.Date(2018, 1, 30, 1, 31, 0, 0, time.UTC), to: time.Date(2018, 1, 30, 1, 35, 0, 0, time.UTC)},
		{from: time.Date(2018, 1, 30, 1, 33, 0, 0, time.UTC), to: time.Date(2018, 1, 30, 1, 34, 0, 0, time.UTC)},
	}
	for t := range times {
		input = append(input, spaniel.New(times[t].from, times[t].to))
	}
}

func ExampleTimeSpan() {
	start, _ := time.Parse("2006-01-02 15:04:05", "2018-01-01 00:00:00")
	timespan := spaniel.New(start, time.Unix(1514768400, 0).UTC())

	fmt.Printf("Start: %v (%v)\nEnd: %v (%v)\n", timespan.Start(), timespan.StartType(), timespan.End(), timespan.EndType())

	// Output:
	// Start: 2018-01-01 00:00:00 +0000 UTC (1)
	// End: 2018-01-01 01:00:00 +0000 UTC (0)
}

func ExampleByStart() {
	sort.Stable(spaniel.ByStart(input))

	for i := range input {
		fmt.Println(input[i].Start())
	}

	// Output:
	// 2018-01-30 00:00:00 +0000 UTC
	// 2018-01-30 00:30:00 +0000 UTC
	// 2018-01-30 01:31:00 +0000 UTC
	// 2018-01-30 01:33:00 +0000 UTC
}

func ExampleByEnd() {
	sort.Stable(spaniel.ByEnd(input))

	for i := range input {
		fmt.Println(input[i].End())
	}

	// Output:
	// 2018-01-30 01:00:00 +0000 UTC
	// 2018-01-30 01:30:00 +0000 UTC
	// 2018-01-30 01:34:00 +0000 UTC
	// 2018-01-30 01:35:00 +0000 UTC
}

func ExampleSpans_Union() {
	union := input.Union()

	for u := range union {
		fmt.Println(union[u].Start(), "->", union[u].End(), ": ", union[u].End().Sub(union[u].Start()))
	}

	// Output:
	// 2018-01-30 00:00:00 +0000 UTC -> 2018-01-30 01:30:00 +0000 UTC :  1h30m0s
	// 2018-01-30 01:31:00 +0000 UTC -> 2018-01-30 01:35:00 +0000 UTC :  4m0s
}

func ExampleSpans_Intersection() {
	intersection := input.Intersection()

	for i := range intersection {
		fmt.Println(intersection[i].Start(), "->", intersection[i].End(), ": ", intersection[i].End().Sub(intersection[i].Start()))
	}

	// Output:
	// 2018-01-30 00:30:00 +0000 UTC -> 2018-01-30 01:00:00 +0000 UTC :  30m0s
	// 2018-01-30 01:33:00 +0000 UTC -> 2018-01-30 01:34:00 +0000 UTC :  1m0s
}
