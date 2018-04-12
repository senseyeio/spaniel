# Spaniel
*Time span handling for Go*

[![godoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/senseyeio/spaniel) [![license](https://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/senseyeio/spaniel/master/LICENSE) [![build](https://img.shields.io/travis/senseyeio/spaniel.svg?style=flat)](https://travis-ci.org/senseyeio/spaniel)

**Spaniel** contains functionality for timespan handling, specifically for merging overlapping timespans and finding the overlaps between multiple timespans.

## Install

This package is "go-gettable", just do:

    go get github.com/senseyeio/spaniel

## Examples

These examples are all available in the ``examples`` folder.

### Basics

Spaniel operates on lists of timespans, it has a built-in Empty timespan for convenience or you can use your own type, so long as it implements the timespan.T interface.

To create a new list of timespans:

	var now = time.Date(2018, 1, 30, 0, 0, 0, 0, time.UTC)

	input := timespan.List{
		timespan.NewEmptyTyped(now, now.Add(1*time.Hour)),
		timespan.NewEmptyTyped(now.Add(30*time.Minute), now.Add(90*time.Minute)),
	}

    
You can then use the Union function to merge the timestamps:

	union := input.Union()
	fmt.Println(union[0].Start(), union[0].End()) // 00:00 - 01:30

Or the Intersection function to find the overlaps:

	intersection := input.Intersection()
	fmt.Println(intersection[0].Start(), intersection[0].End()) // 00:30 - 01:00
 
 ### Types
 
 NewEmptyTyped sets the span to be [) by default - i.e. including the left-most point, excluding the right-most. In other words, [1,2,3) and [3,4,5) do not overlap, but are contiguous. Instants are [] by default (they contain a single time).

If you would like to override these types, you can use NewEmpty:

    openSpan := timespan.NewEmpty(now, now.Add(1*time.Hour)), timespan.Open, timespan.Open)
 
 ### Handlers
 
 If you need to use a more complex object, you can call UnionWithHandler and IntersectionWithHandler. There is an
 example of this in ``examples/handlers.go``.