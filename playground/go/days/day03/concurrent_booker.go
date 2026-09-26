package day03

// Day 3 -- Challenge: Interval Booker, shared.
//
// Day 2's interval booker, now used by many goroutines at the same time
// -- many people booking one meeting room at once. Intervals are
// half-open, as on Day 2: [start, end) includes start but not end, so
// [0,10) and [10,20) do NOT overlap. The one rule that must always hold:
// no two recorded bookings overlap.
//
// On Day 2, booking "the first free slot" took two calls:
//
//	s, _ := booker.FirstFree(from, d) // check
//	booker.Book(s, s+d)               // act
//
// That is check-then-act: another goroutine can take the slot in
// between. This booker has operations that find and book as ONE step
// instead.
//
//	Book(start, end)     if [start, end) overlaps no booking, record it
//	                     and return true. Otherwise return false and
//	                     change nothing. start must be < end; if not,
//	                     return ErrInvalidInterval.
//	Cancel(start, end)   remove the booking with exactly this start and
//	                     end. Returns whether a booking was removed.
//	BookFirstFree(from, duration)
//	                     find the smallest start >= from such that
//	                     [start, start + duration) overlaps no booking,
//	                     book it, and return start -- as one step.
//	                     duration must be > 0; if not, return
//	                     ErrInvalidDuration.
//	BookAll(intervals)   book every interval in the list, or none of
//	                     them. If any one overlaps an existing booking or
//	                     another interval in the same list, return false
//	                     and change nothing. An empty list returns true.
//	                     If any interval has Start >= End, return
//	                     ErrInvalidInterval and change nothing.
//	Bookings()           every booking, sorted by Start, as a new slice
//	                     the caller may keep.
//	Count()              how many bookings are recorded.
//
// Every method may be called by many goroutines at the same time.

import "errors"

var ErrInvalidInterval = errors.New("invalid interval")
var ErrInvalidDuration = errors.New("invalid duration")

type Interval struct {
	Start int64
	End   int64
}

type ConcurrentBooker struct {
	// TODO: choose your own representation.
}

func NewConcurrentBooker() *ConcurrentBooker {
	// TODO
	return &ConcurrentBooker{}
}

func (b *ConcurrentBooker) Book(start, end int64) (bool, error) {
	// TODO
	return false, nil
}

func (b *ConcurrentBooker) Cancel(start, end int64) bool {
	// TODO
	return false
}

func (b *ConcurrentBooker) BookFirstFree(from, duration int64) (int64, error) {
	// TODO
	return 0, nil
}

func (b *ConcurrentBooker) BookAll(intervals []Interval) (bool, error) {
	// TODO
	return false, nil
}

func (b *ConcurrentBooker) Bookings() []Interval {
	// TODO
	return nil
}

func (b *ConcurrentBooker) Count() int {
	// TODO
	return 0
}
