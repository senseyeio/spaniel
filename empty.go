package spaniel

import (
	"time"
)

// Empty represents a simple span of time, with no additional properties. It should be constructed with NewEmpty.
type Empty struct {
	start     time.Time
	end       time.Time
	startType IntervalType
	endType   IntervalType
}

// Start returns the start time of a span
func (ets *Empty) Start() time.Time { return ets.start }

// End returns the end time of a span
func (ets *Empty) End() time.Time { return ets.end }

// StartType returns the start type of a span.
func (ets *Empty) StartType() IntervalType { return ets.startType }

// EndType returns the end type of a span.
func (ets *Empty) EndType() IntervalType { return ets.endType }

// NewEmpty creates a span with just a start and end time, and is used when no handlers are provided to Union or Intersection.
func NewEmpty(start, end time.Time) *Empty {
	return &Empty{start, end, Closed, Closed}
}

// NewEmptyWithTypes creates a span with a start and end time, and accompanying types (Closed or Open). This is used when no handlers are provided to Union or Intersection.
func NewEmptyWithTypes(start, end time.Time, startType, endType IntervalType) *Empty {
	return &Empty{start, end, startType, endType}
}
