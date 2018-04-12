package spaniel

import (
	"sort"
	"time"
)

<<<<<<< HEAD
// IntervalType represents whether the start or end of an interval is Closed or Open.
type IntervalType int

const (
	// Open means that the interval does not include a value
	Open IntervalType = iota
	// Closed means that the interval does include a value
=======
type IntervalType int

const (
	Open IntervalType = iota
>>>>>>> master
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
<<<<<<< HEAD

// MergeHandlerFunc is used by UnionWithHandler to allow for custom functionality when two spans are merged.
// It is passed the two timespans, and the start and end times of the new span.
type MergeHandlerFunc func(mergeInto, mergeFrom T, start, end time.Time, startType, endType IntervalType) T

// IntersectionHandlerFunc is used by IntersectionWithHandler to allow for custom functionality when two spans
// intersect. It is passed the two timespans that intersect, and the start and end times at which they overlap.
type IntersectionHandlerFunc func(intersectingEvent1, intersectingEvent2 T, start, end time.Time, startType, endType IntervalType) T

func getLoosestIntervalType(x, y IntervalType) IntervalType {
	if x > y {
		return x
	}
	return y
}
=======

// ByEnd sorts a list of timespans by their end time
type ByEnd List

func (s ByEnd) Len() int           { return len(s) }
func (s ByEnd) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ByEnd) Less(i, j int) bool { return s[i].End().Before(s[j].End()) }

// MergeHandlerFunc is used by UnionWithHandler to allow for custom functionality when two spans are merged.
// It is passed the two timespans, and the span which would be the result of the merge.
type MergeHandlerFunc func(mergeInto, mergeFrom, mergeSpan T) T

// IntersectionHandlerFunc is used by IntersectionWithHandler to allow for custom functionality when two spans
// intersect. It is passed the two timespans that intersect, span that represents the intersection.
type IntersectionHandlerFunc func(intersectingEvent1, intersectingEvent2, intersectionSpan T) T
>>>>>>> master

func getTightestIntervalType(x, y IntervalType) IntervalType {
	if x < y {
		return x
	}
	return y
}

func getMinStart(a, b T) (time.Time, IntervalType) {
	minStart := b.Start()
	minType := b.StartType()

	if a.Start().Before(b.Start()) {
		minStart = a.Start()
		minType = a.StartType()
	}

	return minStart, minType
}

func getMaxStart(a, b T) (time.Time, IntervalType) {
	maxStart := b.Start()
	maxType := b.StartType()

	if a.Start().After(b.Start()) {
		maxStart = a.Start()
		maxType = a.StartType()
	}
	return maxStart, maxType
}

func getMinEnd(a, b T) (time.Time, IntervalType) {

	minEnd := b.End()
	minType := b.EndType()

	if a.End().Before(b.End()) {
		minEnd = a.End()
		minType = a.EndType()
	}
	return minEnd, minType
}

func getMaxEnd(a, b T) (time.Time, IntervalType) {
	maxEnd := b.End()
	maxType := b.EndType()

	if a.End().After(b.End()) {
		maxEnd = a.End()
		maxType = a.EndType()
	}
	return maxEnd, maxType
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

<<<<<<< HEAD
// IsInstant returns true if the interval is deemed instantaneous
func IsInstant(a T) bool {
	return a.Start().Equal(a.End())
=======
func addUnionTypes(a, b, result T) T {
	// If [a,b,c] and (a,b,c) start type becomes [, end type becomes ] - the more permissive option
	var startType, endType IntervalType

	pair := List{a, b}

	if a.Start().Equal(b.Start()) {
		if a.StartType() == b.StartType() {
			// If the start types match, it's easy
			startType = a.StartType()
		} else {
			// If they differ, one must be closed - so choose that.
			startType = Closed
		}
	} else {
		// Find which starts first and choose its type
		sort.Stable(ByStart(pair))
		startType = pair[0].StartType()
	}

	if a.End().Equal(b.End()) {
		if a.EndType() == b.EndType() {
			// If the start types match, it's easy
			endType = a.EndType()
		} else {
			// If they differ, one must be closed - so choose that.
			endType = Closed
		}
	} else {
		// Find which pair ends last and choose its type
		sort.Stable(ByEnd(pair))
		endType = pair[1].EndType()
	}

	return NewEmptyWithTypes(result.Start(), result.End(), startType, endType)
}

func addIntersectionTypes(a, b, result T) T {
	// If [a,b,c,d] and (b,c], we expect to get (b,c] back - the more restrictive option
	// if [a,b,c,d] and [b,c], we expect to get [b,c] back

	var startType, endType IntervalType

	pair := List{a, b}

	if a.Start().Equal(b.Start()) {
		if a.StartType() == b.StartType() {
			startType = a.StartType()
		} else {
			// If they differ, one must be open - so choose that.
			startType = Open
		}
	} else {
		// Find which starts last and choose its type
		sort.Stable(ByStart(pair))
		startType = pair[1].StartType()
	}

	if a.End().Equal(b.End()) {
		if a.EndType() == b.EndType() {
			endType = a.EndType()
		} else {
			// If they differ, one must be open - so choose that.
			endType = Open
		}
	} else {
		// Find which pair ends first and choose its type
		sort.Stable(ByEnd(pair))
		endType = pair[0].EndType()
	}

	return NewEmptyWithTypes(result.Start(), result.End(), startType, endType)
>>>>>>> master
}

// Returns true if two timespans are side by side
func contiguous(a, b T) bool {
	// [1,2,3,4] [4,5,6,7] - not contiguous
	// [1,2,3,4) [4,5,6,7] - contiguous
	// [1,2,3,4] (4,5,6,7] - contiguous
	// [1,2,3,4) (4,5,6,7] - not contiguous
	// [1,2,3] [5,6,7] - not contiguous

<<<<<<< HEAD
=======
	// Both instants, they can't be contiguous
	if IsInstant(a) && IsInstant(b) {
		return false
	}

>>>>>>> master
	if b.Start().Before(a.Start()) {
		a, b = b, a
	}

<<<<<<< HEAD
=======
	// If one is an instant; if they match the start or end of the other, and the other type is open
	// they are contiguous
	if IsInstant(a) {
		if a.End().Equal(b.Start()) {
			return b.StartType() == Open
		} else if a.End().Equal(b.End()) {
			return b.EndType() == Open
		}
	} else if IsInstant(b) {
		if b.End().Equal(a.Start()) {
			return a.StartType() == Open
		} else if b.End().Equal(a.End()) {
			return a.EndType() == Open
		}
	}

>>>>>>> master
	aEndType := a.EndType()
	bStartType := b.StartType()

	if IsInstant(a) {
		aEndType = Closed
<<<<<<< HEAD
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
=======
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
>>>>>>> master
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
<<<<<<< HEAD
			maxTime, endType := getMaxEnd(a, b)
			minTime, startType := getMinStart(a, b)

			if a.Start().Equal(b.Start()) {
				startType = getLoosestIntervalType(a.StartType(), b.StartType())
			}
			if a.End().Equal(b.End()) {
				endType = getLoosestIntervalType(a.EndType(), b.EndType())
			}

			result[len(result)-1] = mergeHandlerFunc(a, b, minTime, maxTime, startType, endType)
=======
			span := addUnionTypes(a, b, NewEmpty(a.Start(), getMaxTime(a.End(), b.End())))
			result[len(result)-1] = mergeHandlerFunc(a, b, span)
>>>>>>> master
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
<<<<<<< HEAD
	return ts.UnionWithHandler(func(mergeInto, mergeFrom T, start, end time.Time, startType, endType IntervalType) T {
		return NewEmpty(start, end, startType, endType)
=======
	return ts.UnionWithHandler(func(mergeInto, mergeFrom, overlapSpan T) T {
		return overlapSpan
>>>>>>> master
	})
}

// [1,2,3,4,5]  [4,5,6]    = [4,5]
// [1,2,3,4,5) [4,5,6] = [4,5)

// [10:00 - 13:00)     [12:00 - 14:00]   [12:00-13:00)

// [10:00 - 13:00)
//

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
<<<<<<< HEAD
				start, startType := getMaxStart(b, a)
				end, endType := getMinEnd(b, a)
				if a.Start().Equal(b.Start()) {
					startType = getTightestIntervalType(a.StartType(), b.StartType())
				}
				if a.End().Equal(b.End()) {
					endType = getTightestIntervalType(a.EndType(), b.EndType())
				}
				intersection := intersectHandlerFunc(a, b, start, end, startType, endType)
=======
				start := getMaxTime(b.Start(), a.Start())
				end := getMinTime(b.End(), a.End())
				span := addIntersectionTypes(a, b, NewEmpty(start, end))
				intersection := intersectHandlerFunc(a, b, span)
>>>>>>> master
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
<<<<<<< HEAD
	return ts.IntersectionWithHandler(func(intersectingEvent1, intersectingEvent2 T, start, end time.Time, startType, endType IntervalType) T {
		return NewEmpty(start, end, startType, endType)
=======
	return ts.IntersectionWithHandler(func(intersectingEvent1, intersectingEvent2, intersectionSpan T) T {
		return intersectionSpan
>>>>>>> master
	})
}
