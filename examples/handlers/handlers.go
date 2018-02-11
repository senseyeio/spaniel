package handlers

import (
	"fmt"
	"github.com/senseyeio/spaniel"
	"sort"
	"time"
)

type PropertyEvent struct {
	start      time.Time
	end        time.Time
	Properties []string
}

func (e *PropertyEvent) Start() time.Time {
	return e.start
}
func (e *PropertyEvent) End() time.Time {
	return e.end
}

func (e *PropertyEvent) StartType() spaniel.IntervalType {
	return spaniel.Closed
}
func (e *PropertyEvent) EndType() spaniel.IntervalType {
	return spaniel.Closed
}

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

	var now = time.Date(2018, 1, 30, 0, 0, 0, 0, time.UTC)

	input := spaniel.List{
		NewPropertyEvent(now, now.Add(1*time.Hour), []string{"1"}),
		NewPropertyEvent(now.Add(30*time.Minute), now.Add(90*time.Minute), []string{"2"}),
	}

	var mergeHandlerFunc spaniel.MergeHandlerFunc = func(mergeInto, mergeFrom, mergeSpan spaniel.T) spaniel.T {
		a := mergeInto.(*PropertyEvent)
		b := mergeFrom.(*PropertyEvent)
		// Return your object that implements timespan.T
		return NewPropertyEvent(mergeSpan.Start(), mergeSpan.End(), mergeProperties(a.Properties, b.Properties))
	}

	var intersectionHandlerFunc spaniel.IntersectionHandlerFunc = func(intersectingEvent1, intersectingEvent2, intersectionSpan spaniel.T) spaniel.T {
		a := intersectingEvent1.(*PropertyEvent)
		b := intersectingEvent2.(*PropertyEvent)
		// Return your object that implements timespan.T
		return NewPropertyEvent(intersectionSpan.Start(), intersectionSpan.End(), mergeProperties(a.Properties, b.Properties))
	}

	union := input.UnionWithHandler(mergeHandlerFunc)
	fmt.Println(union[0].Start(), union[0].End(), union[0].(*PropertyEvent).Properties)

	intersection := input.IntersectionWithHandler(intersectionHandlerFunc)
	fmt.Println(intersection[0].Start(), intersection[0].End(), intersection[0].(*PropertyEvent).Properties)
}
