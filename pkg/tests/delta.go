package tests

import (
	"time"
)

// unixNanoDelta is used to compare UnixNano dates that should be close enough
// to each other to count as single event.
const unixNanoDelta = 50_000_000.

// TimesAlmostEqual returns true if both times are almost equal.
func TimesAlmostEqual(at, bt time.Time) bool {
	af, bf := float64(at.UnixNano()), float64(bt.UnixNano())
	dt := af - bf
	if dt < -unixNanoDelta || dt > unixNanoDelta {
		return false
	}
	return true
}
