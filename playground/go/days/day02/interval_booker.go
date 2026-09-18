package day02

// Day 2 -- Challenge: Interval Booker.
//
// A booking calendar. A booking takes the time interval [start, end):
// start is included, end is not. So [0,10) and [10,20) do NOT overlap --
// one ends exactly where the other starts.
//
// The one rule that must always hold: no two recorded bookings overlap.
//
//	Book(start, end)          if the interval overlaps no existing
//	                          booking, record it and return true.
//	                          Otherwise return false and change nothing.
//	                          start must be < end; if not, return
//	                          ErrInvalidInterval.
//	Cancel(start, end)        remove a booking. Only an EXACT match
//	                          counts: the same start and the same end.
//	                          Returns whether a booking was removed.
//	FirstFree(from, duration) find the earliest time a booking of length
//	                          duration could start: the smallest
//	                          start >= from such that
//	                          [start, start + duration) overlaps no
//	                          existing booking. duration must be > 0; if
//	                          not, return ErrInvalidDuration.
//	Count()                   how many bookings are recorded.

import "errors"

var ErrInvalidInterval = errors.New("invalid interval")
var ErrInvalidDuration = errors.New("invalid duration")

type IntervalBooker struct {
	// TODO: choose your own representation.
}

func NewIntervalBooker() *IntervalBooker {
	// TODO
	return &IntervalBooker{}
}

func (b *IntervalBooker) Book(start, end int64) (bool, error) {
	// TODO
	return false, nil
}

func (b *IntervalBooker) Cancel(start, end int64) bool {
	// TODO
	return false
}

func (b *IntervalBooker) FirstFree(from, duration int64) (int64, error) {
	// TODO
	return 0, nil
}

func (b *IntervalBooker) Count() int {
	// TODO
	return 0
}
