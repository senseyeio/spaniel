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

// NewInstant creates a span with just a single time.
func NewInstant(time time.Time) *Instant {
	return &Instant{time}
}
