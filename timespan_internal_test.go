package spaniel

import (
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
	startType IntervalType
	endType   IntervalType
}

func TestIntersectionTypes(t *testing.T) {

	for _, tt := range []struct {
		name              string
		a, b              T
		intersectionTypes IntervalTypePair
	}{
		{
			//  [--a--)
			//    (b]
			name:              "overlap c/o o/c",
			a:                 NewEmptyWithTypes(t1, t4, Closed, Open),
			b:                 NewEmptyWithTypes(t2, t3, Open, Closed),
			intersectionTypes: IntervalTypePair{Open, Closed},
		},
		{
			//  [--a--)
			//    [b)
			name:              "overlap c/o c/o",
			a:                 NewEmptyWithTypes(t1, t4, Closed, Open),
			b:                 NewEmptyWithTypes(t2, t3, Closed, Open),
			intersectionTypes: IntervalTypePair{Closed, Open},
		},
		{
			//  (--a--]
			//  [b]
			name:              "overlap o/c start instant",
			a:                 NewEmptyWithTypes(t1, t4, Open, Closed),
			b:                 NewInstant(t1),
			intersectionTypes: IntervalTypePair{Open, Closed},
		},
		{
			//  [---a---]
			//    (-b-)
			name:              "overlap c/c o/o",
			a:                 NewEmptyWithTypes(t1, t4, Closed, Closed),
			b:                 NewEmptyWithTypes(t2, t3, Open, Open),
			intersectionTypes: IntervalTypePair{Open, Open},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			empty := NewEmpty(t1, t2)
			intersectionEvent := addIntersectionTypes(tt.a, tt.b, empty)

			if intersectionEvent.StartType() != tt.intersectionTypes.startType {
				t.Errorf("intersection start")
			}

			if intersectionEvent.EndType() != tt.intersectionTypes.endType {
				t.Errorf("intersection end")
			}
		})
	}
}

func TestMergeTypes(t *testing.T) {

	for _, tt := range []struct {
		name       string
		a, b       T
		mergeTypes IntervalTypePair
	}{
		{
			//  (--a--]
			//        (--b--]
			name:       "contiguous o/c o/c",
			a:          NewEmptyWithTypes(t1, t2, Open, Closed),
			b:          NewEmptyWithTypes(t3, t4, Open, Closed),
			mergeTypes: IntervalTypePair{Open, Closed},
		},
		{
			//  [--a--]
			//        (--b--)
			name:       "contiguous c/c o/o",
			a:          NewEmptyWithTypes(t1, t2, Closed, Closed),
			b:          NewEmptyWithTypes(t3, t4, Open, Open),
			mergeTypes: IntervalTypePair{Closed, Open},
		},
		{
			//  [--a--)
			//        [--b--)
			name:       "contiguous c/o c/o",
			a:          NewEmptyWithTypes(t1, t2, Closed, Open),
			b:          NewEmptyWithTypes(t3, t4, Closed, Open),
			mergeTypes: IntervalTypePair{Closed, Open},
		},
		{
			//  [--a--)
			//    (b]
			name:       "overlap c/o o/c",
			a:          NewEmptyWithTypes(t1, t4, Closed, Open),
			b:          NewEmptyWithTypes(t2, t3, Open, Closed),
			mergeTypes: IntervalTypePair{Closed, Open},
		},
		{
			//  [--a--)
			//    [b)
			name:       "overlap c/o c/o",
			a:          NewEmptyWithTypes(t1, t4, Closed, Open),
			b:          NewEmptyWithTypes(t2, t3, Closed, Open),
			mergeTypes: IntervalTypePair{Closed, Open},
		},
		{
			//  (--a--]
			//  [b]
			name:       "overlap o/c start instant",
			a:          NewEmptyWithTypes(t1, t4, Open, Closed),
			b:          NewInstant(t1),
			mergeTypes: IntervalTypePair{Closed, Closed},
		},
		{
			//  [--a--)
			//       [b]
			name:       "overlap c/o end instant",
			a:          NewEmptyWithTypes(t1, t4, Closed, Open),
			b:          NewInstant(t4),
			mergeTypes: IntervalTypePair{Closed, Closed},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			empty := NewEmpty(t1, t2)
			mergedEvent := addUnionTypes(tt.a, tt.b, empty)

			if mergedEvent.StartType() != tt.mergeTypes.startType {
				t.Errorf("merge start")
			}

			if mergedEvent.EndType() != tt.mergeTypes.endType {
				t.Errorf("merge end")
			}
		})
	}
}

func TestOpenedClosedSpans(t *testing.T) {
	for _, tt := range []struct {
		name               string
		a, b               T
		expectedOverlap    bool
		expectedContiguous bool
	}{
		{
			//  (--a--]
			//          (--b--]
			name:               "no overlap",
			a:                  NewEmptyWithTypes(t1, t2, Open, Closed),
			b:                  NewEmptyWithTypes(t3, t4, Open, Closed),
			expectedOverlap:    false,
			expectedContiguous: false,
		}, {
			//  (--a--]
			//        (--b--]
			name:               "contiguous",
			a:                  NewEmptyWithTypes(t1, t2, Open, Closed),
			b:                  NewEmptyWithTypes(t2, t3, Open, Closed),
			expectedOverlap:    false,
			expectedContiguous: true,
		}, {
			//  (---a----]
			//        (---b----]
			name:               "small intersection",
			a:                  NewEmptyWithTypes(t1, t3, Open, Closed),
			b:                  NewEmptyWithTypes(t2, t4, Open, Closed),
			expectedOverlap:    true,
			expectedContiguous: false,
		}, {
			//  (------a------]
			//      (--b--]
			name:               "one inside the other",
			a:                  NewEmptyWithTypes(t1, t4, Open, Closed),
			b:                  NewEmptyWithTypes(t2, t3, Open, Closed),
			expectedOverlap:    true,
			expectedContiguous: false,
		}, {
			//  (--a--]
			//  (--b--]
			name:               "same",
			a:                  NewEmptyWithTypes(t1, t2, Open, Closed),
			b:                  NewEmptyWithTypes(t1, t2, Open, Closed),
			expectedOverlap:    true,
			expectedContiguous: false,
		}, {
			//  (--a--]
			//           [b]
			name:               "span vs instant, no overlap",
			a:                  NewEmptyWithTypes(t1, t3, Open, Closed),
			b:                  NewInstant(t4),
			expectedOverlap:    false,
			expectedContiguous: false,
		}, {
			//    (--a--]
			//    [b]
			name:               "span vs instant, overlap on the start border",
			a:                  NewEmptyWithTypes(t1, t3, Open, Closed),
			b:                  NewInstant(t1),
			expectedOverlap:    false,
			expectedContiguous: true,
		}, {
			//    (--a--]
			//      [b]
			name:               "span vs instant, overlap in the middle",
			a:                  NewEmptyWithTypes(t1, t3, Open, Closed),
			b:                  NewInstant(t2),
			expectedOverlap:    true,
			expectedContiguous: false,
		}, {
			//    (--a--]
			//        [b]
			name:               "span vs instant, overlap at the end",
			a:                  NewEmptyWithTypes(t1, t3, Open, Closed),
			b:                  NewInstant(t3),
			expectedOverlap:    true,
			expectedContiguous: false,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {

			obtained := overlap(tt.a, tt.b)
			if obtained != tt.expectedOverlap {
				t.Errorf("in order, overlap")
			}
			obtained = overlap(tt.b, tt.a)
			if obtained != tt.expectedOverlap {
				t.Errorf("reversed, overlap")
			}
			obtained = contiguous(tt.a, tt.b)
			if obtained != tt.expectedContiguous {
				t.Errorf("in order, contiguous")
			}
			obtained = contiguous(tt.b, tt.a)
			if obtained != tt.expectedContiguous {
				t.Errorf("reversed, contiguous")
			}
		})
	}
}

func TestClosedOpenedSpans(t *testing.T) {
	for _, tt := range []struct {
		name               string
		a, b               T
		expectedOverlap    bool
		expectedContiguous bool
	}{
		{
			//  [--a--)
			//          [--b--)
			name:               "no overlap",
			a:                  NewEmptyWithTypes(t1, t2, Closed, Open),
			b:                  NewEmptyWithTypes(t3, t4, Closed, Open),
			expectedOverlap:    false,
			expectedContiguous: false,
		}, {
			//  [--a--)
			//        [--b--)
			name:               "contiguous",
			a:                  NewEmptyWithTypes(t1, t2, Closed, Open),
			b:                  NewEmptyWithTypes(t2, t3, Closed, Open),
			expectedOverlap:    false,
			expectedContiguous: true,
		}, {
			//  [---a----)
			//        [---b----)
			name:               "small intersection",
			a:                  NewEmptyWithTypes(t1, t3, Closed, Open),
			b:                  NewEmptyWithTypes(t2, t4, Closed, Open),
			expectedOverlap:    true,
			expectedContiguous: false,
		}, {
			//  [------a------)
			//      [--b--)
			name:               "one inside the other",
			a:                  NewEmptyWithTypes(t1, t4, Closed, Open),
			b:                  NewEmptyWithTypes(t2, t3, Closed, Open),
			expectedOverlap:    true,
			expectedContiguous: false,
		}, {
			//  [--a--)
			//  [--b--)
			name:               "same",
			a:                  NewEmptyWithTypes(t1, t2, Closed, Open),
			b:                  NewEmptyWithTypes(t1, t2, Closed, Open),
			expectedOverlap:    true,
			expectedContiguous: false,
		}, {
			//  [--a--)
			//           [b]
			name:               "span vs instant, no overlap",
			a:                  NewEmptyWithTypes(t1, t3, Closed, Open),
			b:                  NewInstant(t4),
			expectedOverlap:    false,
			expectedContiguous: false,
		}, {
			//    [--a--)
			//    [b]
			name:               "span vs instant, overlap on the start border",
			a:                  NewEmptyWithTypes(t1, t3, Closed, Open),
			b:                  NewInstant(t1),
			expectedOverlap:    true,
			expectedContiguous: false,
		}, {
			//    [--a--)
			//      [b]
			name:               "span vs instant, overlap in the middle",
			a:                  NewEmptyWithTypes(t1, t3, Closed, Open),
			b:                  NewInstant(t2),
			expectedOverlap:    true,
			expectedContiguous: false,
		}, {
			//    [--a--)
			//        [b]
			name:               "span vs instant, overlap at the end",
			a:                  NewEmptyWithTypes(t1, t3, Closed, Open),
			b:                  NewInstant(t3),
			expectedOverlap:    false,
			expectedContiguous: true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {

			obtained := overlap(tt.a, tt.b)
			if obtained != tt.expectedOverlap {
				t.Errorf("in order, overlap")
			}
			obtained = overlap(tt.b, tt.a)
			if obtained != tt.expectedOverlap {
				t.Errorf("reversed, overlap")
			}
			obtained = contiguous(tt.a, tt.b)
			if obtained != tt.expectedContiguous {
				t.Errorf("in order, contiguous")
			}
			obtained = contiguous(tt.b, tt.a)
			if obtained != tt.expectedContiguous {
				t.Errorf("reversed, contiguous")
			}
		})
	}
}

func TestOpenedSpans(t *testing.T) {

	for _, tt := range []struct {
		name               string
		a, b               T
		expectedOverlap    bool
		expectedContiguous bool
	}{
		{
			//  (--a--)
			//          (--b--)
			name:               "no overlap",
			a:                  NewEmptyWithTypes(t1, t2, Open, Open),
			b:                  NewEmptyWithTypes(t3, t4, Open, Open),
			expectedOverlap:    false,
			expectedContiguous: false,
		}, {
			//  (--a--)
			//        (--b--)
			name:               "contiguous",
			a:                  NewEmptyWithTypes(t1, t2, Open, Open),
			b:                  NewEmptyWithTypes(t2, t3, Open, Open),
			expectedOverlap:    false,
			expectedContiguous: false,
		}, {
			//  (---a----)
			//        (---b----)
			name:               "small intersection",
			a:                  NewEmptyWithTypes(t1, t3, Open, Open),
			b:                  NewEmptyWithTypes(t2, t4, Open, Open),
			expectedOverlap:    true,
			expectedContiguous: false,
		}, {
			//  (------a------)
			//      (--b--)
			name:               "one inside the other",
			a:                  NewEmptyWithTypes(t1, t4, Open, Open),
			b:                  NewEmptyWithTypes(t2, t3, Open, Open),
			expectedOverlap:    true,
			expectedContiguous: false,
		}, {
			//  (--a--)
			//  (--b--)
			name:               "same",
			a:                  NewEmptyWithTypes(t1, t2, Open, Open),
			b:                  NewEmptyWithTypes(t1, t2, Open, Open),
			expectedOverlap:    true,
			expectedContiguous: false,
		}, {
			//  (--a--)
			//           [b]
			name:               "span vs instant, no overlap",
			a:                  NewEmptyWithTypes(t1, t3, Open, Open),
			b:                  NewEmptyWithTypes(t4, t4, Open, Open),
			expectedOverlap:    false,
			expectedContiguous: false,
		}, {
			//    (--a--)
			//    [b]
			name:               "span vs instant, overlap on the start border",
			a:                  NewEmptyWithTypes(t1, t3, Open, Open),
			b:                  NewInstant(t1),
			expectedOverlap:    false,
			expectedContiguous: true,
		}, {
			//    (--a--)
			//      [b]
			name:               "span vs instant, overlap in the middle",
			a:                  NewEmptyWithTypes(t1, t3, Open, Open),
			b:                  NewInstant(t2),
			expectedOverlap:    true,
			expectedContiguous: false,
		}, {
			//    (--a--)
			//        [b]
			name:               "span vs instant, overlap at the end",
			a:                  NewEmptyWithTypes(t1, t3, Open, Open),
			b:                  NewInstant(t3),
			expectedOverlap:    false,
			expectedContiguous: true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {

			obtained := overlap(tt.a, tt.b)
			if obtained != tt.expectedOverlap {
				t.Errorf("in order, overlap")
			}
			obtained = overlap(tt.b, tt.a)
			if obtained != tt.expectedOverlap {
				t.Errorf("reversed, overlap")
			}
			obtained = contiguous(tt.a, tt.b)
			if obtained != tt.expectedContiguous {
				t.Errorf("in order, contiguous")
			}
			obtained = contiguous(tt.b, tt.a)
			if obtained != tt.expectedContiguous {
				t.Errorf("reversed, contiguous")
			}
		})
	}
}

func TestClosedSpans(t *testing.T) {
	for _, tt := range []struct {
		name               string
		a, b               T
		expectedOverlap    bool
		expectedContiguous bool
	}{
		{
			//  [--a--]
			//          [--b--]
			name:               "no overlap",
			a:                  NewEmpty(t1, t2),
			b:                  NewEmpty(t3, t4),
			expectedOverlap:    false,
			expectedContiguous: false,
		}, {
			//  [--a--]
			//        [--b--]
			name:               "contiguous",
			a:                  NewEmpty(t1, t2),
			b:                  NewEmpty(t2, t3),
			expectedOverlap:    true,
			expectedContiguous: false,
		}, {
			//  [---a----]
			//        [---b----]
			name:               "small intersection",
			a:                  NewEmpty(t1, t3),
			b:                  NewEmpty(t2, t4),
			expectedOverlap:    true,
			expectedContiguous: false,
		}, {
			//  [------a------]
			//      [--b--]
			name:               "one inside the other",
			a:                  NewEmpty(t1, t4),
			b:                  NewEmpty(t2, t3),
			expectedOverlap:    true,
			expectedContiguous: false,
		}, {
			//  [--a--]
			//  [--b--]
			name:               "same",
			a:                  NewEmpty(t1, t2),
			b:                  NewEmpty(t1, t2),
			expectedOverlap:    true,
			expectedContiguous: false,
		}, {
			//  [--a--]
			//           [b]
			name:               "span vs instant, no overlap",
			a:                  NewEmpty(t1, t3),
			b:                  NewInstant(t4),
			expectedOverlap:    false,
			expectedContiguous: false,
		}, {
			//    [--a--]
			//    [b]
			name:               "span vs instant, overlap on the start border",
			a:                  NewEmpty(t1, t3),
			b:                  NewInstant(t1),
			expectedOverlap:    true,
			expectedContiguous: false,
		}, {
			//    [--a--]
			//      [b]
			name:               "span vs instant, overlap in the middle",
			a:                  NewEmpty(t1, t3),
			b:                  NewInstant(t2),
			expectedOverlap:    true,
			expectedContiguous: false,
		}, {
			//    [--a--]
			//        [b]
			name:               "span vs instant, overlap at the end",
			a:                  NewEmpty(t1, t3),
			b:                  NewInstant(t3),
			expectedOverlap:    true,
			expectedContiguous: false,
		}, {
			//    [a]
			//         [b]
			name:               "both instants, no overlap",
			a:                  NewInstant(t1),
			b:                  NewInstant(t2),
			expectedOverlap:    false,
			expectedContiguous: false,
		}, {
			//    [a]
			//    [b]
			name:               "both instants, overlap",
			a:                  NewInstant(t1),
			b:                  NewInstant(t1),
			expectedOverlap:    true,
			expectedContiguous: false,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			obtained := overlap(tt.a, tt.b)
			if obtained != tt.expectedOverlap {
				t.Errorf("in order, overlap")
			}
			obtained = overlap(tt.b, tt.a)
			if obtained != tt.expectedOverlap {
				t.Errorf("reversed, overlap")
			}
			obtained = contiguous(tt.a, tt.b)
			if obtained != tt.expectedContiguous {
				t.Errorf("in order, contiguous")
			}
			obtained = contiguous(tt.b, tt.a)
			if obtained != tt.expectedContiguous {
				t.Errorf("reversed, contiguous")
			}
		})
	}
}
