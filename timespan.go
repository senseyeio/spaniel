package spaniel

import (
	"sort"
	"time"
)

type IntervalType int
const (
	Open IntervalType = iota
	Closed
)

// T represents a basic timespan, with a start and end time.
type T interface {
	Start() time.Time
	End() time.Time
	LeftType() IntervalType
	RightType() IntervalType
}

// List represents a list of timespans, on which other functions operate.
type List []T

// ByStart sorts a list of timespans by their start time
type ByStart List
func (s ByStart) Len() int { return len(s) }
func (s ByStart) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s ByStart) Less(i, j int) bool { return s[i].Start().Before(s[j].Start())}

// MergeHandlerFunc is used by UnionWithHandler to allow for custom functionality when two spans are merged.
// It is passed the two timespans, and the start and end times of the new span.
type MergeHandlerFunc func(mergeInto, mergeFrom T, start, end time.Time) T

// IntersectionHandlerFunc is used by IntersectionWithHandler to allow for custom functionality when two spans
// intersect. It is passed the two timespans that intersect, and the start and end times at which they overlap.
type IntersectionHandlerFunc func(intersectingEvent1, intersectingEvent2 T, start, end time.Time) T

func getMaxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

func getMinTime(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}

func filter(timeSpans List, filterFunc func(T) bool) List {
	filtered := List{}
	for _, timeSpan := range timeSpans {
		if !filterFunc(timeSpan) {
			filtered = append(filtered, timeSpan)
		}
	}
	return filtered
}

func IsInstant(a T) bool {
	return a.Start().Equal(a.End())
}

func timeWithin(a time.Time, b T, contiguous bool) bool {
	if contiguous || b.LeftType() == Closed {
		if a.Equal(b.Start()) {
			return true
		}
	}

	if contiguous || b.RightType() == Closed {
		if a.Equal(b.End()) {
			return true
		}
	}

	return a.After(b.Start()) && a.Before(b.End())
}



func overlap(a, b T, contiguous bool) bool {

	// [x,y] // includes x and y
	// [x,y) // excludes y
	// (x,y] // excludes x
	// (x,y) // excludes x and y
	// closed: included []
	// open: excluded ()

	// If we use an instant (i.e. start == end), it doesn't use
	// left and right types.

	if IsInstant(a) {
		return timeWithin(a.Start(), b, contiguous)
	}
	if IsInstant(b) {
		return timeWithin(b.Start(), a, contiguous)
	}

	// Given [e,g] and [f,h]
	// If e > h || g < f, overlap == false

	c_1 := false // is a_s after b_e
	if a.Start().After(b.End()) {
		// given 5,6,7 and 1,2,3,4
		c_1 = true
	}
	if a.Start().Equal(b.End()) {
		// given 5,6,7 and 1,2,3,4,5
		if contiguous || (a.LeftType() == Closed && b.RightType() == Closed) {
			// a: 5,6,7 b: 1,2,3,4,5
		} else {
			c_1 = true
		}
	}

	c_2 := false // is a_e before b_s
	if a.End().Before(b.Start()) {
		c_2 = true
	}
	if a.End().Equal(b.Start()) {
		// given 1,2,3,4,5 and 5,6,7,8
		if contiguous || (a.RightType() == Closed && b.LeftType() == Closed) {
			// a: 1, 2, 3, 4, 5
			// b: 5, 6, 7, 8
			// not before
		} else {
			c_2 = true
		}
	}

	if c_1 || c_2 {
		return false
	}

	return true
}

// UnionWithHandler returns a list of TimeSpans representing the union of all of the time spans.
// For example, given a list [A,B] where A and B overlap, a list [C] would be returned, with the timespan C spanning
// both A and B. The provided handler is passed the source and destination timespans), and the start and end time of
// the currently merged timespan.
func (ts List) UnionWithHandler(mergeHandlerFunc MergeHandlerFunc) List {

	if len(ts) < 2 {
		return ts
	}

	var sorted List
	sorted = append(sorted, ts...)
	sort.Stable(ByStart(sorted))

	result := List{sorted[0]}

	for _, b := range sorted[1:] {
		// A: current timespan in merged array; B: current timespan in sorted array
		// If B overlaps with A, it can be merged with A.
		a := result[len(result)-1]
		if overlap(a, b, true) {
			result[len(result)-1] = mergeHandlerFunc(a, b, a.Start(), getMaxTime(a.End(), b.End()))
			continue
		}
		result = append(result, b)
	}

	return result
}

// Union returns a list of TimeSpans representing the union of all of the time spans.
// For example, given a list [A,B] where A and B overlap, a list [C] would be returned, with the timespan C spanning
// both A and B.
func (ts List) Union() List {
	return ts.UnionWithHandler(func(mergeInto, mergeFrom T, start, end time.Time) T {
		return NewEmpty(start, end)
	})
}

// IntersectionWithHandler returns a list of TimeSpans representing the overlaps between the contained time spans.
// For example, given a list [A,B] where A and B overlap, a list [C] would be returned, with the timespan C covering
// the intersection of the A and B. The provided handler function is notified of the two timespans that have been found
// to overlap, and the start and end time of the overlap.
func (ts List) IntersectionWithHandler(intersectHandlerFunc IntersectionHandlerFunc) List {
	var sorted List
	sorted = append(sorted, ts...)
	sort.Stable(ByStart(sorted))

	actives := List{sorted[0]}

	intersections := List{}

	for _, b := range sorted[1:] {
		// Tidy up the active span list
		actives = filter(actives, func(t T) bool {
			// If this value is identical to one in actives, don't filter it.
			if b.Start() == t.Start() && b.End() == t.End() {
				return false
			}
			// If this value starts after the one in actives finishes, filter the active.
			return b.Start().After(t.End()) || b.Start().Equal(t.End())
		})

		for _, a := range actives {
			if overlap(a, b, false) {
				start := getMaxTime(b.Start(), a.Start())
				end := getMinTime(b.End(), a.End())
				intersection := intersectHandlerFunc(a, b, start, end)
				intersections = append(intersections, intersection)
			}
		}
		actives = append(actives, b)
	}
	return intersections
}

// Intersection returns a list of TimeSpans representing the overlaps between the contained time spans.
// For example, given a list [A,B] where A and B overlap, a list [C] would be returned,
// with the timespan C covering the intersection of the A and B.
func (ts List) Intersection() List {
	return ts.IntersectionWithHandler(func(intersectingEvent1, intersectingEvent2 T, start, end time.Time) T {
		return NewEmpty(start, end)
	})
}
