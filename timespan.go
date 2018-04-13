package spaniel

import (
	"sort"
	"time"
)

// EndPointType represents whether the start or end of an interval is Closed or Open.
type EndPointType int

const (
	// Open means that the interval does not include a value
	Open EndPointType = iota
	// Closed means that the interval does include a value
	Closed
)

type EndPoint struct {
	Element time.Time
	Type    EndPointType
}

func (e EndPoint) Before(a EndPoint) bool { return e.Element.Before(a.Element) }
func (e EndPoint) After(a EndPoint) bool  { return e.Element.After(a.Element) }
func (e EndPoint) Equal(a EndPoint) bool  { return e.Element.Equal(a.Element) }

// T represents a basic timespan, with a start and end time.
type T interface {
	Start() time.Time
	StartType() EndPointType
	End() time.Time
	EndType() EndPointType
}

// List represents a list of timespans, on which other functions operate.
type List []T

// ByStart sorts a list of timespans by their start time
type ByStart List

func (s ByStart) Len() int           { return len(s) }
func (s ByStart) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ByStart) Less(i, j int) bool { return s[i].Start().Before(s[j].Start()) }

// MergeHandlerFunc is used by UnionWithHandler to allow for custom functionality when two spans are merged.
// It is passed the two timespans, and the start and end times of the new span.
type MergeHandlerFunc func(mergeInto, mergeFrom, mergeSpan T) T

// IntersectionHandlerFunc is used by IntersectionWithHandler to allow for custom functionality when two spans
// intersect. It is passed the two timespans that intersect, and the start and end times at which they overlap.
type IntersectionHandlerFunc func(intersectingEvent1, intersectingEvent2, intersectionSpan T) T

func getLoosestIntervalType(x, y EndPointType) EndPointType {
	if x > y {
		return x
	}
	return y
}

func getTightestIntervalType(x, y EndPointType) EndPointType {
	if x < y {
		return x
	}
	return y
}

func getMin(a, b EndPoint) EndPoint {
	if a.Before(b) {
		return a
	}
	return b
}

func getMax(a, b EndPoint) EndPoint {
	if a.After(b) {
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

// IsInstant returns true if the interval is deemed instantaneous
func IsInstant(a T) bool {
	return a.Start().Equal(a.End())
}

// Returns true if two timespans are side by side
func contiguous(a, b T) bool {
	// [1,2,3,4] [4,5,6,7] - not contiguous
	// [1,2,3,4) [4,5,6,7] - contiguous
	// [1,2,3,4] (4,5,6,7] - contiguous
	// [1,2,3,4) (4,5,6,7] - not contiguous
	// [1,2,3] [5,6,7] - not contiguous
	// [1] (1,2,3] - contiguous

	// Two instants can't be contiguous
	if IsInstant(a) && IsInstant(b) {
		return false
	}

	if b.Start().Before(a.Start()) {
		a, b = b, a
	}

	aStartType := a.StartType()
	aEndType := a.EndType()
	bStartType := b.StartType()

	if IsInstant(a) {
		aEndType = Closed
		aStartType = Closed
	}
	if IsInstant(b) {
		bStartType = Closed
	}

	// If a and b start at the same time, just check that their start types are different.
	if a.Start().Equal(b.Start()) {
		return aStartType != bStartType
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

	c1 := false // is a_s after b_e
	if a.Start().After(b.End()) {
		c1 = true
	} else if a.Start().Equal(b.End()) {
		c1 = (aStartType == Open || bEndType == Open)
	}

	c2 := false // is a_e before b_s
	if a.End().Before(b.Start()) {
		c2 = true
	} else if a.End().Equal(b.Start()) {
		c2 = (aEndType == Open || bStartType == Open)
	}

	if c1 || c2 {
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

			spanStart := getMin(EndPoint{a.Start(), a.StartType()}, EndPoint{b.Start(), b.StartType()})
			spanEnd := getMax(EndPoint{a.End(), a.EndType()}, EndPoint{b.End(), b.EndType()})

			if a.Start().Equal(b.Start()) {
				spanStart.Type = getLoosestIntervalType(a.StartType(), b.StartType())
			}
			if a.End().Equal(b.End()) {
				spanEnd.Type = getLoosestIntervalType(a.EndType(), b.EndType())
			}

			span := NewEmptyWithTypes(spanStart.Element, spanEnd.Element, spanStart.Type, spanEnd.Type)
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
	return ts.UnionWithHandler(func(mergeInto, mergeFrom, mergeSpan T) T {
		return mergeSpan
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
			return b.Start().After(t.End())
		})

		for _, a := range actives {
			if overlap(a, b) {
				spanStart := getMax(EndPoint{a.Start(), a.StartType()}, EndPoint{b.Start(), b.StartType()})
				spanEnd := getMin(EndPoint{a.End(), a.EndType()}, EndPoint{b.End(), b.EndType()})

				if a.Start().Equal(b.Start()) {
					spanStart.Type = getTightestIntervalType(a.StartType(), b.StartType())
				}
				if a.End().Equal(b.End()) {
					spanEnd.Type = getTightestIntervalType(a.EndType(), b.EndType())
				}
				span := NewEmptyWithTypes(spanStart.Element, spanEnd.Element, spanStart.Type, spanEnd.Type)
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
