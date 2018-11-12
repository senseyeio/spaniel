package spaniel

import (
	"time"
)

// TimeSpan represents a simple span of time, with no additional properties. It should be constructed with NewEmpty.
type TimeSpan struct {
	start     time.Time
	end       time.Time
	startType EndPointType
	endType   EndPointType
}

// Start returns the start time of a span
func (ts TimeSpan) Start() time.Time { return ts.start }

// End returns the end time of a span
func (ts TimeSpan) End() time.Time { return ts.end }

// StartType returns the type of the start of the interval (Open in this case)
func (ts TimeSpan) StartType() EndPointType { return ts.startType }

// EndType returns the type of the end of the interval (Closed in this case)
func (ts TimeSpan) EndType() EndPointType { return ts.endType }

// String returns a string representation of a timespan
func (ts TimeSpan) String() string {
	s := ""
	if ts.StartType() == Closed {
		s += "["
	} else {
		s += "("
	}

	s += ts.Start().String()
	if ts.Start() != ts.End() {
		s += ","
		s += ts.End().String()
	}

	if ts.EndType() == Closed {
		s += "]"
	} else {
		s += ")"
	}
	return s
}

// MarshalJSON implements json.Marshal
func (ts TimeSpan) MarshalJSON() ([]byte, error) {
	o := struct {
		Start         time.Time `json:"start"`
		End           time.Time `json:"end"`
		StartIncluded bool      `json:"start_included"`
		EndIncluded   bool      `json:"end_included"`
	}{
		Start: ts.start,
		End:   ts.end,
	}

	o.StartIncluded = endPointInclusionMarshal(ts.startType)
	o.EndIncluded = endPointInclusionMarshal(ts.endType)

	return json.Marshal(o)
}

// UnmarshalJSON implements json.Unmarshal
func (ts *TimeSpan) UnmarshalJSON(b []byte) (err error) {
	var i struct {
		Start         time.Time `json:"start"`
		End           time.Time `json:"end"`
		StartIncluded bool      `json:"start_included"`
		EndIncluded   bool      `json:"end_included"`
	}

	err = json.Unmarshal(b, &i)
	if err != nil {
		return err
	}

	ts.start = i.Start
	ts.end = i.End
	ts.startType = endPointIncluseionUnmarhsal(i.StartIncluded)
	ts.endType = endPointIncluseionUnmarhsal(i.EndIncluded)

	return
}

func endPointInclusionMarshal(e EndPointType) (included bool) {
	switch e {
	case Open:
		included = false
	case Closed:
		included = true
	}

	return included
}

func endPointIncluseionUnmarhsal(b bool) (e EndPointType) {
	switch b {
	case true:
		e = Closed
	case false:
		e = Open
	}

	return e
}

// NewWithTypes creates a span with just a start and end time, and associated types, and is used when no handlers are provided to Union or Intersection.
func NewWithTypes(start, end time.Time, startType, endType EndPointType) *TimeSpan {
	return &TimeSpan{start, end, startType, endType}
}

// NewInstant creates a span with just a single time.
func NewInstant(time time.Time) *TimeSpan {
	return New(time, time)
}

// New creates a span with a start and end time, with the types set to [] for instants and [) for spans.
func New(start time.Time, end time.Time) *TimeSpan {
	if start.Equal(end) {
		// An instantaneous event has to be Closed (i.e. inclusive)
		return NewWithTypes(start, end, Closed, Closed)
	}
	return NewWithTypes(start, end, Closed, Open)
}
