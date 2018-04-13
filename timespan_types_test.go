package spaniel_test

import (
	timespan "github.com/senseyeio/spaniel"
	"testing"
	"time"
)

var (
	t1 = time.Date(2018, 1, 30, 0, 0, 0, 0, time.UTC)
	t2 = t1.Add(time.Second)
	t3 = t2.Add(time.Second)
	t4 = t3.Add(time.Second)
)

type IntervalTypePair struct {
	startType timespan.EndPointType
	endType   timespan.EndPointType
}

func TestTypedUnion(t *testing.T) {

	for _, tt := range []struct {
		name       string
		a, b       timespan.T
		mergeTypes IntervalTypePair
	}{
		{
			//  (--a--]
			//        (--b--]
			name:       "contiguous o/c o/c",
			a:          timespan.NewEmptyWithTypes(t1, t2, timespan.Open, timespan.Closed),
			b:          timespan.NewEmptyWithTypes(t2, t3, timespan.Open, timespan.Closed),
			mergeTypes: IntervalTypePair{timespan.Open, timespan.Closed},
		},
		{
			//  [--a--]
			//        (--b--)
			name:       "contiguous c/c o/o",
			a:          timespan.NewEmptyWithTypes(t1, t2, timespan.Closed, timespan.Closed),
			b:          timespan.NewEmptyWithTypes(t2, t3, timespan.Open, timespan.Open),
			mergeTypes: IntervalTypePair{timespan.Closed, timespan.Open},
		},
		{
			//  [--a--)
			//        [--b--)
			name:       "contiguous c/o c/o",
			a:          timespan.NewEmptyWithTypes(t1, t2, timespan.Closed, timespan.Open),
			b:          timespan.NewEmptyWithTypes(t2, t3, timespan.Closed, timespan.Open),
			mergeTypes: IntervalTypePair{timespan.Closed, timespan.Open},
		},
		{
			//  [--a--)
			//    (b]
			name:       "overlap c/o o/c",
			a:          timespan.NewEmptyWithTypes(t1, t4, timespan.Closed, timespan.Open),
			b:          timespan.NewEmptyWithTypes(t2, t3, timespan.Open, timespan.Closed),
			mergeTypes: IntervalTypePair{timespan.Closed, timespan.Open},
		},
		{
			//  [--a--)
			//    [b)
			name:       "overlap c/o c/o",
			a:          timespan.NewEmptyWithTypes(t1, t4, timespan.Closed, timespan.Open),
			b:          timespan.NewEmptyWithTypes(t2, t3, timespan.Closed, timespan.Open),
			mergeTypes: IntervalTypePair{timespan.Closed, timespan.Open},
		},
		{
			//  (--a--]
			//  [b]
			name:       "overlap o/c start instant",
			a:          timespan.NewEmptyWithTypes(t1, t4, timespan.Open, timespan.Closed),
			b:          timespan.NewInstant(t1),
			mergeTypes: IntervalTypePair{timespan.Closed, timespan.Closed},
		},
		{
			//  [--a--)
			//       [b]
			name:       "overlap c/o end instant",
			a:          timespan.NewEmptyWithTypes(t1, t4, timespan.Closed, timespan.Open),
			b:          timespan.NewInstant(t4),
			mergeTypes: IntervalTypePair{timespan.Closed, timespan.Closed},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			input := timespan.List{
				tt.a, tt.b,
			}
			merges := input.Union()
			if len(merges) != 1 {
				t.Errorf("in order, no merges")
				return
			}

			if merges[0].StartType() != tt.mergeTypes.startType {
				t.Errorf("in order, merge start")
			}

			if merges[0].EndType() != tt.mergeTypes.endType {
				t.Errorf("in order, merge end")
			}

			input = timespan.List{
				tt.b, tt.a,
			}
			merges = input.Union()
			if len(merges) != 1 {
				t.Errorf("reversed, no merges")
				return
			}

			if merges[0].StartType() != tt.mergeTypes.startType {
				t.Errorf("reversed, merge start")
			}

			if merges[0].EndType() != tt.mergeTypes.endType {
				t.Errorf("reversed, merge end")
			}
		})
	}
}

func TestTypedIntersection(t *testing.T) {

	for _, tt := range []struct {
		name              string
		a, b              timespan.T
		intersectionTypes IntervalTypePair
	}{
		{
			//  [--a--)
			//    (b]
			name:              "overlap c/o o/c",
			a:                 timespan.NewEmptyWithTypes(t1, t4, timespan.Closed, timespan.Open),
			b:                 timespan.NewEmptyWithTypes(t2, t3, timespan.Open, timespan.Closed),
			intersectionTypes: IntervalTypePair{timespan.Open, timespan.Closed},
		},
		{
			//  [--a--)
			//    [b)
			name:              "overlap c/o c/o",
			a:                 timespan.NewEmptyWithTypes(t1, t4, timespan.Closed, timespan.Open),
			b:                 timespan.NewEmptyWithTypes(t2, t3, timespan.Closed, timespan.Open),
			intersectionTypes: IntervalTypePair{timespan.Closed, timespan.Open},
		},
		{
			//  [--a--]
			//  (--b--)
			name:              "overlap c/c o/o",
			a:                 timespan.NewEmptyWithTypes(t1, t4, timespan.Closed, timespan.Closed),
			b:                 timespan.NewEmptyWithTypes(t1, t4, timespan.Open, timespan.Open),
			intersectionTypes: IntervalTypePair{timespan.Open, timespan.Open},
		},
		{
			//  [---a---]
			//    (-b-)
			name:              "overlap c/c o/o",
			a:                 timespan.NewEmptyWithTypes(t1, t4, timespan.Closed, timespan.Closed),
			b:                 timespan.NewEmptyWithTypes(t2, t3, timespan.Open, timespan.Open),
			intersectionTypes: IntervalTypePair{timespan.Open, timespan.Open},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			input := timespan.List{
				tt.a, tt.b,
			}
			intersections := input.Intersection()
			if len(intersections) != 1 {
				return
			}

			if intersections[0].StartType() != tt.intersectionTypes.startType {
				t.Errorf("in order, intersection start")
			}

			if intersections[0].EndType() != tt.intersectionTypes.endType {
				t.Errorf("in order, intersection end")
			}

			input = timespan.List{
				tt.b, tt.a,
			}
			intersections = input.Intersection()
			if len(intersections) != 1 {
				return
			}

			if intersections[0].StartType() != tt.intersectionTypes.startType {
				t.Errorf("reversed, intersection start")
			}

			if intersections[0].EndType() != tt.intersectionTypes.endType {
				t.Errorf("reversed, intersection end")
			}
		})
	}
}
