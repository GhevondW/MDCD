package day02

// Day 2 -- Challenge: Interval Booker.
//
// A booking calendar over HALF-OPEN intervals [start, end): the invariant
// is that no two recorded bookings ever overlap. [0,10) and [10,20) do
// not overlap.
//
//	Book(start, end)          record the booking and return true iff it
//	                          overlaps no existing booking; otherwise
//	                          return false and change nothing.
//	                          Contract: start < end, else return
//	                          ErrInvalidInterval.
//	Cancel(start, end)        remove a booking by EXACT match only.
//	FirstFree(from, duration) the smallest start >= from such that
//	                          [start, start + duration) overlaps no
//	                          existing booking.
//	                          Contract: duration > 0, else return
//	                          ErrInvalidDuration.
//	Count()                   number of recorded bookings.

import "errors"

var ErrInvalidInterval = errors.New("invalid interval")
var ErrInvalidDuration = errors.New("invalid duration")

type IntervalBooker struct {
	// TODO: choose your own representation.
}

func NewIntervalBooker() *IntervalBooker {
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
