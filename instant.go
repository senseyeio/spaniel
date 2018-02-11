package spaniel

import (
	"time"
)

// Instant represents an instantaneous event with no duration or additional properties. It should be constructed with NewInstant.
type Instant struct {
	time time.Time
}

// Start returns the start time of a span
func (its *Instant) Start() time.Time { return its.time }

// End returns the end time of a span - due to it being an Instant, this is the same as the start time.
func (its *Instant) End() time.Time { return its.time }

// StartType returns the start type of a span. In this case, Closed as it is an instantaneous event.
func (ets *Instant) StartType() IntervalType { return Closed }

// EndType returns the end type of a span. In this case, Closed as it is an instantaneous event.
func (ets *Instant) EndType() IntervalType { return Closed }

// NewInstant creates a span with just a single time.
func NewInstant(time time.Time) *Instant {
	return &Instant{time}
}
