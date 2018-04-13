package main

import (
	"fmt"
	"github.com/senseyeio/spaniel"
	"sort"
	"time"
)

// PropertyEvent represents an event with an associated list of property strings.
type PropertyEvent struct {
	start      time.Time
	end        time.Time
	Properties []string
}

// Start represents the start time of the property event.
func (e *PropertyEvent) Start() time.Time {
	return e.start
}

// End represents the end time of the property event.
func (e *PropertyEvent) End() time.Time {
	return e.end
}

// StartType represents the type of the start of the interval (Closed in this case).
func (e *PropertyEvent) StartType() spaniel.EndPointType {
	return spaniel.Closed
}

// EndType represents the type of the end of the interval (Open in this case).
func (e *PropertyEvent) EndType() spaniel.EndPointType {
	return spaniel.Open
}

// NewPropertyEvent creates a new PropertyEvent with start and end times and a list of properties.
func NewPropertyEvent(start time.Time, end time.Time, properties []string) *PropertyEvent {
	return &PropertyEvent{start, end, properties}
}

var mergeProperties = func(a []string, b []string) []string {
	for _, mergeFromProperty := range b {
		found := false
		for _, mergeInProperty := range a {
			if mergeInProperty == mergeFromProperty {
				found = true
			}
		}
		if !found {
			a = append(a, mergeFromProperty)
		}
	}
	sort.Strings(a)
	return a
}

func main() {


	// Times at half-hourly intervals
	var t1 = time.Date(2018, 1, 30, 0, 0, 0, 0, time.UTC)
	var t2 = time.Date(2018, 1, 30, 0, 30, 0, 0, time.UTC)
	var t3 = time.Date(2018, 1, 30, 1, 0, 0, 0, time.UTC)
	var t4 = time.Date(2018, 1, 30, 1, 30, 0, 0, time.UTC)

	input := spaniel.List{
		NewPropertyEvent(t1, t3, []string{"1"}),
		NewPropertyEvent(t2, t4, []string{"2"}),
	}

	var mergeHandlerFunc spaniel.MergeHandlerFunc = func(mergeInto, mergeFrom spaniel.T, mergeSpan spaniel.T) spaniel.T {
		a, ok := mergeInto.(*PropertyEvent)
		if !ok {
			panic(fmt.Sprintf("Expected mergeInto to be a PropertyEvent"))
		}
		b, ok := mergeFrom.(*PropertyEvent)
		if !ok {
			panic(fmt.Errorf("Expected mergeFrom to be a PropertyEvent"))
		}

		// Return your object that implements timespan.T
		return NewPropertyEvent(mergeSpan.Start(), mergeSpan.End(), mergeProperties(a.Properties, b.Properties))
	}

	var intersectionHandlerFunc spaniel.IntersectionHandlerFunc = func(intersectingEvent1, intersectingEvent2, intersectionSpan spaniel.T) spaniel.T {
		a, ok := intersectingEvent1.(*PropertyEvent)
		if !ok {
			panic(fmt.Errorf("Expected intersectingEvent1 to be a PropertyEvent"))
		}
		b, ok := intersectingEvent2.(*PropertyEvent)
		if !ok {
			panic(fmt.Errorf("Expected intersectingEvent2 to be a PropertyEvent"))
		}

		// Return your object that implements timespan.T
		return NewPropertyEvent(intersectionSpan.Start(), intersectionSpan.End(), mergeProperties(a.Properties, b.Properties))
	}

	union := input.UnionWithHandler(mergeHandlerFunc)
	fmt.Println(union[0].Start(), union[0].End(), union[0].(*PropertyEvent).Properties)

	intersection := input.IntersectionWithHandler(intersectionHandlerFunc)
	fmt.Println(intersection[0].Start(), intersection[0].End(), intersection[0].(*PropertyEvent).Properties)
}
