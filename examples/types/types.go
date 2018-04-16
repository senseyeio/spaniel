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
	var t5 = time.Date(2018, 1, 30, 2, 0, 0, 0, time.UTC)

	input := timespan.List{
		timespan.NewWithTypes(t1, t3, timespan.Open, timespan.Open),
		timespan.NewWithTypes(t2, t4, timespan.Open, timespan.Open),
	}

	union := input.Union()
	fmt.Println("As both timespans are Open, they are not contiguous - so won't be merged:")
	fmt.Println(union)

	intersection := input.Intersection()
	fmt.Println("And there will be no intersections:")
	fmt.Println(intersection)

	input = timespan.List{
		timespan.NewWithTypes(t1, t3, timespan.Closed, timespan.Closed),
		timespan.NewWithTypes(t3, t5, timespan.Closed, timespan.Closed),
	}
	union = input.Union()
	fmt.Println("If they are Closed, they will overlap:")
	fmt.Println(union)

	intersection = input.Intersection()
	fmt.Println("And there will be an instantaneous intersection:")
	fmt.Println(intersection)

	input = timespan.List{
		timespan.NewWithTypes(t1, t3, timespan.Closed, timespan.Open),
		timespan.NewWithTypes(t3, t5, timespan.Closed, timespan.Open),
	}
	union = input.Union()
	fmt.Println("If they are both [) they can be merged as they are contiguous:")
	fmt.Println(union)

	intersection = input.Intersection()
	fmt.Println("But there will be no intersection:")
	fmt.Println(intersection)
}
