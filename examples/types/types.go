package main

import (
	"fmt"
	"github.com/senseyeio/spaniel"
	"time"
)

func main() {

	// Times at half-hourly intervals
	var t1 = time.Date(2018, 1, 30, 0, 0, 0, 0, time.UTC)
	var t2 = time.Date(2018, 1, 30, 0, 30, 0, 0, time.UTC)
	var t3 = time.Date(2018, 1, 30, 1, 0, 0, 0, time.UTC)
	var t4 = time.Date(2018, 1, 30, 1, 30, 0, 0, time.UTC)
	var t5 = time.Date(2018, 1, 30, 2, 0, 0, 0, time.UTC)

	input := spaniel.List{
		spaniel.NewWithTypes(t1, t3, spaniel.Open, spaniel.Open),
		spaniel.NewWithTypes(t2, t4, spaniel.Open, spaniel.Open),
	}

	union := input.Union()
	fmt.Println("As both timespans are Open, they are not contiguous - so won't be merged:")
	fmt.Println(union)

	intersection := input.Intersection()
	fmt.Println("And there will be no intersections:")
	fmt.Println(intersection)

	input = spaniel.List{
		spaniel.NewWithTypes(t1, t3, spaniel.Closed, spaniel.Closed),
		spaniel.NewWithTypes(t3, t5, spaniel.Closed, spaniel.Closed),
	}
	union = input.Union()
	fmt.Println("If they are Closed, they will overlap:")
	fmt.Println(union)

	intersection = input.Intersection()
	fmt.Println("And there will be an instantaneous intersection:")
	fmt.Println(intersection)

	input = spaniel.List{
		spaniel.NewWithTypes(t1, t3, spaniel.Closed, spaniel.Open),
		spaniel.NewWithTypes(t3, t5, spaniel.Closed, spaniel.Open),
	}
	union = input.Union()
	fmt.Println("If they are both [) they can be merged as they are contiguous:")
	fmt.Println(union)

	intersection = input.Intersection()
	fmt.Println("But there will be no intersection:")
	fmt.Println(intersection)
}
