# Spaniel
*Time span handling for Go*

[![godoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/senseyeio/spaniel) [![license](https://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/senseyeio/spaniel/master/LICENSE) [![build](https://img.shields.io/travis/senseyeio/spaniel.svg?style=flat)](https://travis-ci.org/senseyeio/spaniel)

**Spaniel** contains functionality for timespan handling, specifically for merging overlapping timespans and finding the intersections between multiple timespans. It lets you specify the type of interval you want to use (open, closed), and provide handlers for when you want to add more functionality when merging/intersecting.

## Install

This package is "go-gettable", just do:

    go get github.com/senseyeio/spaniel

## Basics

Spaniel operates on lists of timespans, where a timespan is represented as the interval between a start and end time. It has a built-in minimal timespan representation for convenience, or you can use your own type, so long as it implements the timespan.T interface.

To import spaniel and create a new list of timespans:

	package main

	import (
  		timespan "github.com/senseyeio/spaniel"
		"time"
		"fmt"
	)

	func main() {
		// Times at half-hourly intervals
		var t1 = time.Date(2018, 1, 30, 0, 0, 0, 0, time.UTC)
		var t2 = time.Date(2018, 1, 30, 0, 30, 0, 0, time.UTC)
		var t3 = time.Date(2018, 1, 30, 1, 0, 0, 0, time.UTC)
		var t4 = time.Date(2018, 1, 30, 1, 30, 0, 0, time.UTC)

		input := timespan.Spans{
			timespan.New(t1, t3),
			timespan.New(t2, t4),
		}
		fmt.Println(input)
	}
    
You can then use the Union function to merge the timestamps:

	union := input.Union()
	fmt.Println(union[0].Start(), "->", union[0].End()) // 2018-01-30 00:00:00 +0000 UTC -> 2018-01-30 01:30:00 +0000 UTC

Or the Intersection function to find the overlaps:

	intersection := input.Intersection()
	fmt.Println(intersection[0].Start(), "->", intersection[0].End()) // 2018-01-30 00:30:00 +0000 UTC -> 2018-01-30 01:00:00 +0000 UTC
 
## Types
 
`timespan.New` sets the span to be [`[)`](https://en.wikipedia.org/wiki/Interval_(mathematics)#Notations_for_intervals) by default - i.e. including the left-most point, excluding the right-most. In other words, `[1,2,3)` and `[3,4,5)` do not overlap, but are contiguous. Instants are `[]` by default (they contain a single time).

If you would like to override these types, you can use NewWithTypes:

    openSpan := timespan.NewWithTypes(t1, t3, timespan.Open, timespan.Open)
 
You can see a more involved example of types in ``examples/types/types.go``
 
## Handlers
 
If you need to use a more complex object, you can call UnionWithHandler and IntersectionWithHandler. There is an example of this in ``examples/handlers/handlers.go``.


## More Examples

All of the above examples are available in the ``examples`` folder.
