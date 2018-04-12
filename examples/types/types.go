package main

import (
	"fmt"
	"github.com/senseyeio/spaniel"
	"time"
)

func main() {

	var now = time.Date(2018, 1, 30, 0, 0, 0, 0, time.UTC)

	input := spaniel.List{
		spaniel.NewEmpty(now, now.Add(1*time.Hour), spaniel.Open, spaniel.Open),
		spaniel.NewEmpty(now.Add(1*time.Hour), now.Add(2*time.Hour), spaniel.Open, spaniel.Open),
	}

	union := input.Union()
	fmt.Println("As both timespans are Open, they are not contiguous - so won't be merged:")
	fmt.Println(union)

	intersection := input.Intersection()
	fmt.Println("And there will be no intersections:")
	fmt.Println(intersection)

	input = spaniel.List{
		spaniel.NewEmpty(now, now.Add(1*time.Hour), spaniel.Closed, spaniel.Closed),
		spaniel.NewEmpty(now.Add(1*time.Hour), now.Add(2*time.Hour), spaniel.Closed, spaniel.Closed),
	}
	union = input.Union()
	fmt.Println("If they are Closed, they will overlap:")
	fmt.Println(union)

	intersection = input.Intersection()
	fmt.Println("And there will be an instantaneous intersection:")
	fmt.Println(intersection)

	input = spaniel.List{
		spaniel.NewEmpty(now, now.Add(1*time.Hour), spaniel.Closed, spaniel.Open),
		spaniel.NewEmpty(now.Add(1*time.Hour), now.Add(2*time.Hour), spaniel.Closed, spaniel.Open),
	}
	union = input.Union()
	fmt.Println("If they are both [) they can be merged as they are contiguous:")
	fmt.Println(union)

	intersection = input.Intersection()
	fmt.Println("But there will be no intersection:")
	fmt.Println(intersection)
}
