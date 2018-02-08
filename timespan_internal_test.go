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
		name                    string
		a, b                    T
		expectedNonContiguous   bool
		expectedAllowContiguous bool
	}{
		{
			//  [--a--]
			//          [--b--]
			name: "no overlap",
			a:    NewEmpty(t1, t2),
			b:    NewEmpty(t3, t4),
			expectedNonContiguous:   false,
			expectedAllowContiguous: false,
		}, {
			//  [--a--]
			//        [--b--]
			name: "contiguous",
			a:    NewEmpty(t1, t2),
			b:    NewEmpty(t2, t3),
			expectedNonContiguous:   false,
			expectedAllowContiguous: true,
		}, {
			//  [---a----]
			//        [---b----]
			name: "small intersection",
			a:    NewEmpty(t1, t3),
			b:    NewEmpty(t2, t4),
			expectedNonContiguous:   true,
			expectedAllowContiguous: true,
		}, {
			//  [------a------]
			//      [--b--]
			name: "one inside the other",
			a:    NewEmpty(t1, t4),
			b:    NewEmpty(t2, t3),
			expectedNonContiguous:   true,
			expectedAllowContiguous: true,
		}, {
			//  [--a--]
			//  [--b--]
			name: "same",
			a:    NewEmpty(t1, t2),
			b:    NewEmpty(t1, t2),
			expectedNonContiguous:   true,
			expectedAllowContiguous: true,
		}, {
			//  [--a--]
			//           [b]
			name: "span vs instant, no overlap",
			a:    NewEmpty(t1, t3),
			b:    NewEmpty(t4, t4),
			expectedNonContiguous:   false,
			expectedAllowContiguous: false,
		}, {
			//    [--a--]
			//    [b]
			name: "span vs instant, overlap on the start border",
			a:    NewEmpty(t1, t3),
			b:    NewEmpty(t1, t1),
			expectedNonContiguous:   false,
			expectedAllowContiguous: true,
		}, {
			//    [--a--]
			//      [b]
			name: "span vs instant, overlap in the middle",
			a:    NewEmpty(t1, t3),
			b:    NewEmpty(t2, t2),
			expectedNonContiguous:   false,
			expectedAllowContiguous: true,
		}, {
			//    [--a--]
			//        [b]
			name: "span vs instant, overlap at the end",
			a:    NewEmpty(t1, t3),
			b:    NewEmpty(t3, t3),
			expectedNonContiguous:   false,
			expectedAllowContiguous: true,
		}, {
			//    [a]
			//         [b]
			name: "both instants, no overlap",
			a:    NewEmpty(t1, t1),
			b:    NewEmpty(t2, t2),
			expectedNonContiguous:   false,
			expectedAllowContiguous: false,
		}, {
			//    [a]
			//    [b]
			name: "both instants, overlap",
			a:    NewEmpty(t1, t1),
			b:    NewEmpty(t1, t1),
			expectedNonContiguous:   true,
			expectedAllowContiguous: true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			obtained := overlap(tt.a, tt.b, false)
			if obtained != tt.expectedNonContiguous {
				t.Errorf("in order, allowContiguous=false")
			}
			obtained = overlap(tt.b, tt.a, false)
			if obtained != tt.expectedNonContiguous {
				t.Errorf("reversed, allowContiguous=false")
			}
			obtained = overlap(tt.a, tt.b, true)
			if obtained != tt.expectedAllowContiguous {
				t.Errorf("in order, allowContiguous=true")
			}
			obtained = overlap(tt.b, tt.a, true)
			if obtained != tt.expectedAllowContiguous {
				t.Errorf("reversed, allowContiguous=true")
			}
		})
	}
}
