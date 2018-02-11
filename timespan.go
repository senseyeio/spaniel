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
	StartType() IntervalType
	EndType() IntervalType
}

// List represents a list of timespans, on which other functions operate.
type List []T

// ByStart sorts a list of timespans by their start time
type ByStart List

func (s ByStart) Len() int           { return len(s) }
func (s ByStart) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ByStart) Less(i, j int) bool { return s[i].Start().Before(s[j].Start()) }

// MergeHandlerFunc is used by UnionWithHandler to allow for custom functionality when two spans are merged.
// It is passed the two timespans, and the span which would be the result of the merge.
type MergeHandlerFunc func(mergeInto, mergeFrom, mergeSpan T) T

// IntersectionHandlerFunc is used by IntersectionWithHandler to allow for custom functionality when two spans
// intersect. It is passed the two timespans that intersect, span that represents the intersection.
type IntersectionHandlerFunc func(intersectingEvent1, intersectingEvent2, intersectionSpan T) T

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

// Returns true if two timespans are side by side
func contiguous(a, b T) bool {
	// [1,2,3,4] [4,5,6,7] - not contiguous
	// [1,2,3,4) [4,5,6,7] - contiguous
	// [1,2,3,4] (4,5,6,7] - contiguous
	// [1,2,3,4) (4,5,6,7] - not contiguous
	// [1,2,3] [5,6,7] - not contiguous

	if b.Start().Before(a.Start()) {
		a, b = b, a
	}

	aEndType := a.EndType()
	bStartType := b.StartType()

	if IsInstant(a) {
		aEndType = Closed
	}
	if IsInstant(b) {
		bStartType = Closed
	}

	// To be contiguous the ranges have to overlap on the first/last time
	if !(a.End().Equal(b.Start())) {
		return false
	}

	if aEndType == bStartType {
		return false
	}
	return true
}

func IsInstant(t T) bool {
	return t.Start().Equal(t.End())
}

// Returns true if two timespans overlap
func overlap(a, b T) bool {
	// [1,2,3,4] [4,5,6,7] - intersects
	// [1,2,3,4) [4,5,6,7] - doesn't intersect
	// [1,2,3,4] (4,5,6,7] - doesn't intersect
	// [1,2,3,4) (4,5,6,7] - doesn't intersect

	aStartType := a.StartType()
	aEndType := a.EndType()
	bStartType := b.StartType()
	bEndType := b.EndType()

	if IsInstant(a) {
		aStartType = Closed
		aEndType = Closed
	}
	if IsInstant(b) {
		bStartType = Closed
		bEndType = Closed
	}

	// Given [a_s,a_e] and [b_s,b_e]
	// If a_s > b_e || a_e < b_s, overlap == false

	c_1 := false // is a_s after b_e
	if a.Start().After(b.End()) {
		c_1 = true
	} else if a.Start().Equal(b.End()) {
		c_1 = (aStartType == Open || bEndType == Open)
	}

	c_2 := false // is a_e before b_s
	if a.End().Before(b.Start()) {
		c_2 = true
	} else if a.End().Equal(b.Start()) {
		c_2 = (aEndType == Open || bStartType == Open)
	}

	if c_1 || c_2 {
		return false
	}

	return true
}

// UnionWithHandler returns a list of TimeSpans representing the union of all of the time spans.
// For example, given a list [A,B] where A and B overlap, a list [C] would be returned, with the timespan C spanning
// both A and B. The provided handler is passed the source and destination timespans, and the currently merged empty timespan.
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
		if overlap(a, b) || contiguous(a, b) {
			span := NewEmpty(a.Start(), getMaxTime(a.End(), b.End()))
			result[len(result)-1] = mergeHandlerFunc(a, b, span)
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
	return ts.UnionWithHandler(func(mergeInto, mergeFrom, overlapSpan T) T {
		return overlapSpan
	})
}

// IntersectionWithHandler returns a list of TimeSpans representing the overlaps between the contained time spans.
// For example, given a list [A,B] where A and B overlap, a list [C] would be returned, with the timespan C covering
// the intersection of the A and B. The provided handler function is notified of the two timespans that have been found
// to overlap, and the span representing the overlap.
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
			if overlap(a, b) {
				start := getMaxTime(b.Start(), a.Start())
				end := getMinTime(b.End(), a.End())
				span := NewEmpty(start, end)
				intersection := intersectHandlerFunc(a, b, span)
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
	return ts.IntersectionWithHandler(func(intersectingEvent1, intersectingEvent2, intersectionSpan T) T {
		return intersectionSpan
	})
}
