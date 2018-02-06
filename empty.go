package spaniel

import (
	"time"
)

// Empty represents a simple span of time, with no additional properties. It should be constructed with NewEmpty.
type Empty struct {
	start time.Time
	end   time.Time
}

// Start returns the start time of a span
func (ets *Empty) Start() time.Time { return ets.start }
// End returns the end time of a span
func (ets *Empty) End() time.Time   { return ets.end }

// NewEmpty creates a span with just a start and end time, and is used when no handlers are provided to Union or Intersection.
func NewEmpty(start time.Time, end time.Time) *Empty {
	return &Empty{start, end}
}
