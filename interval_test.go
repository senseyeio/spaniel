package spaniel_test

import (
	"reflect"
	"fmt"
	"sort"
	"testing"
	"time"

	timespan "github.com/senseyeio/spaniel"
)

type Event struct {
	start     time.Time
	end       time.Time
	startType timespan.EndPointType
	endType   timespan.EndPointType
}

func NewEvent(start time.Time, end time.Time) *Event {
	return &Event{start, end, timespan.Closed, timespan.Open}
}

func (e *Event) Start() time.Time {
	return e.start
}
func (e *Event) End() time.Time {
	return e.end
}

func (e *Event) StartType() timespan.EndPointType {
	return e.startType
}

func (e *Event) EndType() timespan.EndPointType {
	return e.endType
}

func (e *Event) SetStartType(startType timespan.EndPointType) {
	e.startType = startType
}

func (e *Event) SetEndType(endType timespan.EndPointType) {
	e.endType = endType
}

func (e *Event) String() string {
	return e.start.Format("15:04:05") + "-" + e.end.Format("15:04:05")
}

type PropertyEvent struct {
	Event
	Properties []string
}

func NewPropertyEvent(start time.Time, end time.Time, properties []string) *PropertyEvent {
	return &PropertyEvent{Event{start, end, timespan.Closed, timespan.Open}, properties}
}

var now = time.Date(2018, 1, 30, 0, 0, 0, 0, time.UTC)

func expectEqual(t *testing.T, x interface{}, y interface{}) {

	if !reflect.DeepEqual(x, y) {
		t.Fatalf("Expected %v to equal %v", x, y)
	}
}

func TestHandlers(t *testing.T) {

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

	var mergePropertiesFunc = func(mergeInto timespan.Span, mergeFrom timespan.Span, mergeSpan timespan.Span) timespan.Span {
		// The union will contain the properties from both merged events
		a, ok := mergeInto.(*PropertyEvent)
		if !ok {
			t.Fatalf("Couldn't cast mergeInto timespan into a PropertyEvent")
		}

		b, ok := mergeFrom.(*PropertyEvent)
		if !ok {
			t.Fatalf("Couldn't cast mergeFrom timespan into a PropertyEvent")
		}
		return NewPropertyEvent(mergeSpan.Start(), mergeSpan.End(), mergeProperties(a.Properties, b.Properties))
	}

	var intersectPropertiesFunc = func(intersectingEvent1 timespan.Span, intersectingEvent2 timespan.Span, intersectSpan timespan.Span) timespan.Span {
		// The intersection will contain the properties from both intersecting events
		a, ok := intersectingEvent1.(*PropertyEvent)
		if !ok {
			t.Fatalf("Couldn't cast intersectingEvent1 timespan into a PropertyEvent")
		}
		b, ok := intersectingEvent2.(*PropertyEvent)
		if !ok {
			t.Fatalf("Couldn't cast intersectingEvent2 timespan into a PropertyEvent")
		}
		return NewPropertyEvent(intersectSpan.Start(), intersectSpan.End(), mergeProperties(a.Properties, b.Properties))

	}

	t.Run("Should allow a combination of union and intersection with handlers", func(t *testing.T) {
		// Input timestamps are merged
		inputTimeSpans := timespan.Spans{
			NewPropertyEvent(now, now.Add(time.Hour), []string{}),                            // 00:00 - 01:00
			NewPropertyEvent(now.Add(30*time.Minute), now.Add(30*time.Minute), []string{}),   // An instantaneous event at 00:30
			NewPropertyEvent(now.Add(90*time.Minute), now.Add(2*time.Hour), []string{}),      // 01:30 - 02:00
			NewPropertyEvent(now.Add(100*time.Minute), now.Add(130*time.Minute), []string{}), // 01:40 - 02:10
		}

		mergedInputTimeSpans := inputTimeSpans.UnionWithHandler(mergePropertiesFunc)
		expectEqual(t, mergedInputTimeSpans, timespan.Spans{
			NewPropertyEvent(now, now.Add(time.Hour), []string{}),                           // 00:00 - 01:00
			NewPropertyEvent(now.Add(90*time.Minute), now.Add(130*time.Minute), []string{}), // 01:30 - 02:10
		})

		// Get some property timespans
		propSpans1 := timespan.Spans{
			NewPropertyEvent(now.Add(10*time.Minute), now.Add(25*time.Minute), []string{"prop1"}),  // 00:10 - 00:25
			NewPropertyEvent(now.Add(50*time.Minute), now.Add(100*time.Minute), []string{"prop1"}), // 00:50 - 01:40
			NewPropertyEvent(now.Add(55*time.Minute), now.Add(60*time.Minute), []string{"prop1"}),  // 00:55 - 01:35
		}.UnionWithHandler(mergePropertiesFunc)

		expectEqual(t, propSpans1, timespan.Spans{
			NewPropertyEvent(now.Add(10*time.Minute), now.Add(25*time.Minute), []string{"prop1"}),  // 00:10 - 00:25
			NewPropertyEvent(now.Add(50*time.Minute), now.Add(100*time.Minute), []string{"prop1"}), // 00:50 - 01:40
		})

		// Intersect the property spans with the input timespans
		intersectionPropSpans1 := append(propSpans1, mergedInputTimeSpans...).IntersectionWithHandler(intersectPropertiesFunc)

		expectEqual(t, intersectionPropSpans1, timespan.Spans{
			NewPropertyEvent(now.Add(10*time.Minute), now.Add(25*time.Minute), []string{"prop1"}),  // 00:10 - 00:25
			NewPropertyEvent(now.Add(50*time.Minute), now.Add(60*time.Minute), []string{"prop1"}),  // 00:50 - 01:00
			NewPropertyEvent(now.Add(90*time.Minute), now.Add(100*time.Minute), []string{"prop1"}), // 01:30 - 01:40
		})

		propSpans2 := timespan.Spans{
			NewPropertyEvent(now.Add(35*time.Minute), now.Add(110*time.Minute), []string{"prop2"}), // 00:35 - 01:50
			NewPropertyEvent(now.Add(2*time.Hour), now.Add(150*time.Minute), []string{"prop2"}),    // 02:00 - 02:30
		}.UnionWithHandler(mergePropertiesFunc)

		expectEqual(t, propSpans2, timespan.Spans{
			NewPropertyEvent(now.Add(35*time.Minute), now.Add(110*time.Minute), []string{"prop2"}), // 00:35 - 01:50
			NewPropertyEvent(now.Add(2*time.Hour), now.Add(150*time.Minute), []string{"prop2"}),    // 02:00 - 02:30
		})

		// Intersect the property spans with the input timespans
		intersectionPropSpans2 := append(propSpans2, mergedInputTimeSpans...).IntersectionWithHandler(intersectPropertiesFunc)

		expectEqual(t, intersectionPropSpans2, timespan.Spans{
			NewPropertyEvent(now.Add(35*time.Minute), now.Add(60*time.Minute), []string{"prop2"}),   // 00:35 - 01:00
			NewPropertyEvent(now.Add(90*time.Minute), now.Add(110*time.Minute), []string{"prop2"}),  // 01:30 - 01:50
			NewPropertyEvent(now.Add(120*time.Minute), now.Add(130*time.Minute), []string{"prop2"}), // 02:00 - 02:10
		})

		// Merge the intersected rule spans
		outputPropSpans := append(intersectionPropSpans1, intersectionPropSpans2...).UnionWithHandler(mergePropertiesFunc)

		expectEqual(t, outputPropSpans, timespan.Spans{
			NewPropertyEvent(now.Add(10*time.Minute), now.Add(25*time.Minute), []string{"prop1"}),
			NewPropertyEvent(now.Add(35*time.Minute), now.Add(1*time.Hour), []string{"prop1", "prop2"}),
			NewPropertyEvent(now.Add(90*time.Minute), now.Add(110*time.Minute), []string{"prop1", "prop2"}),
			NewPropertyEvent(now.Add(2*time.Hour), now.Add(130*time.Minute), []string{"prop2"}),
		})
	})
}

func TestUnion(t *testing.T) {

	t.Run("Should keep two instants separate", func(t *testing.T) {
		a := timespan.New(now, now)
		b := timespan.New(now.Add(2*time.Hour), now.Add(2*time.Hour))
		events := timespan.Spans{a, b}
		after := events.Union()
		expectEqual(t, after, events)
	})

	t.Run("Should keep two separate timespans separate", func(t *testing.T) {
		a := timespan.New(now, now.Add(time.Hour))
		b := timespan.New(now.Add(2*time.Hour), now.Add(3*time.Hour))
		events := timespan.Spans{a, b}
		after := events.Union()
		expectEqual(t, after, events)
	})

	t.Run("Should handle a single timespan by returning that timespan", func(t *testing.T) {
		a := timespan.New(now, now.Add(time.Hour))
		events := timespan.Spans{a}
		after := events.Union()
		expectEqual(t, after, events)
	})

	t.Run("Should merge two overlapping timespans", func(t *testing.T) {
		a := timespan.New(now, now.Add(time.Hour))
		b := timespan.New(now.Add(30*time.Minute), now.Add(3*time.Hour))
		expected := timespan.Spans{timespan.New(a.Start(), b.End())}
		events := timespan.Spans{a, b}
		after := events.Union()
		expectEqual(t, after, expected)
	})

	t.Run("Should merge two consecutive timespans", func(t *testing.T) {
		a := timespan.New(now, now.Add(time.Hour))
		b := timespan.New(now.Add(time.Hour), now.Add(3*time.Hour))
		expected := timespan.Spans{timespan.New(a.Start(), b.End())}
		events := timespan.Spans{a, b}
		after := events.Union()
		expectEqual(t, after, expected)
	})

	t.Run("Should merge three overlapping timespans", func(t *testing.T) {
		a := timespan.New(now, now.Add(time.Hour))
		b := timespan.New(now.Add(30*time.Minute), now.Add(3*time.Hour))
		c := timespan.New(now.Add(20*time.Minute), now.Add(35*time.Minute))
		expected := timespan.Spans{timespan.New(a.Start(), b.End())}
		events := timespan.Spans{a, b, c}
		after := events.Union()
		expectEqual(t, after, expected)
	})

	t.Run("Should merge two timespans overlapped by one timespan", func(t *testing.T) {
		a := timespan.New(now, now.Add(30*time.Minute))
		b := timespan.New(now.Add(15*time.Minute), now.Add(60*time.Minute))
		c := timespan.New(now.Add(45*time.Minute), now.Add(75*time.Minute))
		expected := timespan.Spans{
			timespan.New(a.Start(), c.End()),
		}
		events := timespan.Spans{a, b, c}
		after := events.Union()
		expectEqual(t, after, expected)
	})

	t.Run("Should merge one timespan overlapped by two timespans", func(t *testing.T) {
		a := timespan.New(now, now.Add(60*time.Minute))
		b := timespan.New(now.Add(15*time.Minute), now.Add(20*time.Minute))
		c := timespan.New(now.Add(40*time.Minute), now.Add(45*time.Minute))
		expected := timespan.Spans{
			a,
		}
		events := timespan.Spans{a, b, c}
		after := events.Union()
		expectEqual(t, after, expected)
	})

	t.Run("Should merge three identical timespans", func(t *testing.T) {
		a := timespan.New(now, now.Add(60*time.Minute))
		b := timespan.New(now, now.Add(60*time.Minute))
		c := timespan.New(now, now.Add(60*time.Minute))
		expected := timespan.Spans{
			a,
		}
		events := timespan.Spans{a, b, c}
		after := events.Union()
		expectEqual(t, after, expected)
	})

	t.Run("Should not merge two consecutive timespans if non-inclusive", func(t *testing.T) {
		a := NewEvent(now, now.Add(time.Hour))
		a.SetEndType(timespan.Open)

		b := NewEvent(now.Add(time.Hour), now.Add(3*time.Hour))
		b.SetStartType(timespan.Open)
		expected := timespan.Spans{a, b}
		events := timespan.Spans{a, b}
		after := events.Union()
		expectEqual(t, after, expected)
	})

}

func TestIntersection(t *testing.T) {

	t.Run("Should find overlaps for two instants", func(t *testing.T) {
		a := timespan.New(now, now)
		b := timespan.New(now, now)
		expected := timespan.Spans{timespan.New(a.Start(), a.End())}
		events := timespan.Spans{a, b}
		after := events.Intersection()
		expectEqual(t, after, expected)
	})

	t.Run("Should find no overlaps if timespans are separate", func(t *testing.T) {
		a := timespan.New(now, now.Add(time.Hour))
		b := timespan.New(now.Add(2*time.Hour), now.Add(3*time.Hour))
		events := timespan.Spans{a, b}
		after := events.Intersection()
		expectEqual(t, after, timespan.Spans{})
	})

	t.Run("Should find no intersections if a single timespan", func(t *testing.T) {
		a := timespan.New(now, now.Add(time.Hour))
		events := timespan.Spans{a}
		after := events.Intersection()
		expectEqual(t, after, timespan.Spans{})
	})

	t.Run("Should return the intersection of two overlapping timespans", func(t *testing.T) {
		a := timespan.New(now, now.Add(time.Hour))
		b := timespan.New(now.Add(30*time.Minute), now.Add(3*time.Hour))
		expected := timespan.Spans{timespan.New(b.Start(), a.End())}
		events := timespan.Spans{a, b}
		after := events.Intersection()
		expectEqual(t, after, expected)
	})

	t.Run("Should return the intersection of three overlapping timespans", func(t *testing.T) {
		a := timespan.New(now, now.Add(time.Hour))
		b := timespan.New(now.Add(30*time.Minute), now.Add(3*time.Hour))
		c := timespan.New(now.Add(20*time.Minute), now.Add(35*time.Minute))
		events := timespan.Spans{a, b, c}
		expected := timespan.Spans{
			timespan.New(now.Add(20*time.Minute), now.Add(35*time.Minute)),
			timespan.New(now.Add(30*time.Minute), now.Add(1*time.Hour)),
			timespan.New(now.Add(30*time.Minute), now.Add(35*time.Minute)),
		}
		after := events.Intersection()
		expectEqual(t, after, expected)
	})

	t.Run("Should return the intersection of three overlapping timespans", func(t *testing.T) {
		a := NewEvent(now.Add(1*time.Hour), now.Add(6*time.Hour))
		b := NewEvent(now.Add(2*time.Hour), now.Add(5*time.Hour))
		c := NewEvent(now.Add(3*time.Hour), now.Add(4*time.Hour))
		events := timespan.Spans{a, b, c}

		expected := timespan.Spans{
			timespan.New(now.Add(2*time.Hour), now.Add(5*time.Hour)),
			timespan.New(now.Add(3*time.Hour), now.Add(4*time.Hour)),
			timespan.New(now.Add(3*time.Hour), now.Add(4*time.Hour)),
		}

		after := events.Intersection()
		expectEqual(t, after, expected)
	})

	t.Run("Should return the intersections of two timespans overlapped by one timespan", func(t *testing.T) {
		a := NewEvent(now, now.Add(30*time.Minute))
		b := NewEvent(now.Add(15*time.Minute), now.Add(60*time.Minute))
		c := NewEvent(now.Add(45*time.Minute), now.Add(75*time.Minute))
		expected := timespan.Spans{
			timespan.New(b.Start(), a.End()),
			timespan.New(c.Start(), b.End()),
		}
		events := timespan.Spans{a, b, c}
		after := events.Intersection()
		expectEqual(t, after, expected)
	})

	t.Run("Should return the intersections of one timespan overlapped by two timespans", func(t *testing.T) {
		a := NewEvent(now, now.Add(60*time.Minute))
		b := NewEvent(now.Add(15*time.Minute), now.Add(20*time.Minute))
		c := NewEvent(now.Add(40*time.Minute), now.Add(45*time.Minute))
		expected := timespan.Spans{
			timespan.New(b.Start(), b.End()),
			timespan.New(c.Start(), c.End()),
		}
		events := timespan.Spans{a, b, c}
		after := events.Intersection()
		expectEqual(t, after, expected)
	})

	t.Run("Should return the intersections of three identical timespans", func(t *testing.T) {
		a := NewEvent(now, now.Add(60*time.Minute))
		b := NewEvent(now, now.Add(60*time.Minute))
		c := NewEvent(now, now.Add(60*time.Minute))
		expected := timespan.Spans{
			timespan.New(a.Start(), a.End()),
			timespan.New(b.Start(), b.End()),
			timespan.New(c.Start(), c.End()),
		}
		events := timespan.Spans{a, b, c}
		after := events.Intersection()
		expectEqual(t, after, expected)
	})

	t.Run("Should be no intersections if consecutive", func(t *testing.T) {
		a := NewEvent(now, now.Add(time.Hour))
		b := NewEvent(now.Add(time.Hour), now.Add(3*time.Hour))
		expected := timespan.Spans{}
		events := timespan.Spans{a, b}
		after := events.Intersection()
		expectEqual(t, after, expected)
	})
}

func TestIntersectionWith(t *testing.T) {
	a := timespan.New(now, now.Add(time.Hour))
	b := timespan.New(now.Add(2*time.Hour), now.Add(3*time.Hour))
	c := timespan.New(now.Add(4*time.Hour), now.Add(5*time.Hour))

	candidate := timespan.New(now.Add(45*time.Minute), now.Add(3*time.Hour))
	events := timespan.Spans{a, b, c}
	after := events.IntersectionWith(candidate)

	t.Run("Should find overlap for partially overlapping", func(t *testing.T) {
		expected := timespan.New(now.Add(45*time.Minute), now.Add(time.Hour))
		expectEqual(t, after[0], expected)
	})

	t.Run("Should find overlap for totally overlapping", func(t *testing.T) {
		expectEqual(t, after[1], b)
	})

	t.Run("Should indicate nil for non-overlapping", func(t *testing.T) {
		var expected *timespan.TimeSpan
		expectEqual(t, after[2], expected)
	})
}
