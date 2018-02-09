package spaniel_test

import (
	"testing"
	"time"
	timespan "github.com/senseyeio/spaniel"
)
//var now = time.Date(2018, 1, 30, 0, 0, 0, 0, time.UTC)

func TestTypedUnion(t *testing.T) {

	t.Run("Two overlapping inclusive ranges should result in an inclusive range", func(t *testing.T) {
		a := timespan.NewEmpty(now, now.Add(2*time.Hour), timespan.Closed, timespan.Closed)
		b := timespan.NewEmpty(now.Add(2*time.Hour), now.Add(2*time.Hour), timespan.Closed, timespan.Closed)
		expected := timespan.List{timespan.NewEmpty(a.Start(), b.End(), timespan.Closed, timespan.Closed)}
		events := timespan.List{a, b}
		after := events.Union()
		expectEqual(t, after, expected)
	})


	t.Run("Two ranges half-closed at outer ends should keep that nature", func(t *testing.T) {
		a := timespan.NewEmpty(now, now.Add(2*time.Hour), timespan.Open, timespan.Closed)
		b := timespan.NewEmpty(now.Add(2*time.Hour), now.Add(2*time.Hour), timespan.Closed, timespan.Open)
		expected := timespan.List{timespan.NewEmpty(a.Start(), b.End(), timespan.Open, timespan.Open)}
		events := timespan.List{a, b}
		after := events.Union()
		expectEqual(t, after, expected)
	})


	t.Run("Two ranges half-closed at inner ends should keep that nature", func(t *testing.T) {
		a := timespan.NewEmpty(now, now.Add(2*time.Hour), timespan.Closed, timespan.Open)
		b := timespan.NewEmpty(now.Add(2*time.Hour), now.Add(2*time.Hour), timespan.Open, timespan.Closed)
		expected := timespan.List{timespan.NewEmpty(a.Start(), b.End(), timespan.Closed, timespan.Closed)}
		events := timespan.List{a, b}
		after := events.Union()
		expectEqual(t, after, expected)
	})


	t.Run("Two duplicate ranges should keep the more inclusive type", func(t *testing.T) {
		a := timespan.NewEmpty(now.Add(1*time.Hour), now.Add(2*time.Hour), timespan.Closed, timespan.Closed)
		b := timespan.NewEmpty(now.Add(1*time.Hour), now.Add(2*time.Hour), timespan.Open, timespan.Open)
		expected := timespan.List{timespan.NewEmpty(now.Add(1*time.Hour), now.Add(2*time.Hour), timespan.Closed, timespan.Closed)}
		events := timespan.List{a, b}
		after := events.Union()
		expectEqual(t, after, expected)
	})
}