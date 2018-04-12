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

func TestOverlap(t *testing.T) {
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
			a:                  NewEmptyTyped(t1, t2),
			b:                  NewEmptyTyped(t3, t4),
			expectedOverlap:    false,
			expectedContiguous: false,
		}, {
			//  [--a--]
			//        [--b--]
			name:               "contiguous",
			a:                  NewEmptyTyped(t1, t2),
			b:                  NewEmptyTyped(t2, t3),
			expectedOverlap:    false,
			expectedContiguous: true,
		}, {
			//  [---a----]
			//        [---b----]
			name:               "small intersection",
			a:                  NewEmptyTyped(t1, t3),
			b:                  NewEmptyTyped(t2, t4),
			expectedOverlap:    true,
			expectedContiguous: false,
		}, {
			//  [------a------]
			//      [--b--]
			name:               "one inside the other",
			a:                  NewEmptyTyped(t1, t4),
			b:                  NewEmptyTyped(t2, t3),
			expectedOverlap:    true,
			expectedContiguous: false,
		}, {
			//  [--a--]
			//  [--b--]
			name:               "same",
			a:                  NewEmptyTyped(t1, t2),
			b:                  NewEmptyTyped(t1, t2),
			expectedOverlap:    true,
			expectedContiguous: false,
		}, {
			//  [--a--]
			//           [b]
			name:               "span vs instant, no overlap",
			a:                  NewEmptyTyped(t1, t3),
			b:                  NewEmptyTyped(t4, t4),
			expectedOverlap:    false,
			expectedContiguous: false,
		}, {
			//    [--a--]
			//    [b]
			name:               "span vs instant, overlap on the start border",
			a:                  NewEmptyTyped(t1, t3),
			b:                  NewEmptyTyped(t1, t1),
			expectedOverlap:    true,
			expectedContiguous: false,
		}, {
			//    [--a--]
			//      [b]
			name:               "span vs instant, overlap in the middle",
			a:                  NewEmptyTyped(t1, t3),
			b:                  NewEmptyTyped(t2, t2),
			expectedOverlap:    true,
			expectedContiguous: false,
		}, {
			//    [--a--]
			//        [b]
			name:               "span vs instant, overlap at the end",
			a:                  NewEmptyTyped(t1, t3),
			b:                  NewEmptyTyped(t3, t3),
			expectedOverlap:    false,
			expectedContiguous: true,
		}, {
			//    [a]
			//         [b]
			name:               "both instants, no overlap",
			a:                  NewEmptyTyped(t1, t1),
			b:                  NewEmptyTyped(t2, t2),
			expectedOverlap:    false,
			expectedContiguous: false,
		}, {
			//    [a]
			//    [b]
			name:               "both instants, overlap",
			a:                  NewEmptyTyped(t1, t1),
			b:                  NewEmptyTyped(t1, t1),
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
