package main

import (
	"fmt"
	"github.com/senseyeio/spaniel"
	"time"
)

func main() {

	var now = time.Date(2018, 1, 30, 0, 0, 0, 0, time.UTC)

	input := spaniel.List{
		spaniel.NewEmpty(now, now.Add(1*time.Hour), spaniel.Closed, spaniel.Open),
		spaniel.NewEmpty(now.Add(30*time.Minute), now.Add(90*time.Minute), spaniel.Closed, spaniel.Open),
	}

	union := input.Union()
	fmt.Println(union[0].Start(), union[0].End())

	intersection := input.Intersection()
	fmt.Println(intersection[0].Start(), intersection[0].End())
}
