package day02

import (
	"errors"
	"testing"
)

func mustBook(t *testing.T, b *IntervalBooker, start, end int64) {
	t.Helper()
	ok, err := b.Book(start, end)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatalf("expected Book(%d, %d) to succeed", start, end)
	}
}

func TestIntervalBookerBookOnEmptyCalendarSucceeds(t *testing.T) {
	b := NewIntervalBooker()
	mustBook(t, b, 10, 20)
	if b.Count() != 1 {
		t.Errorf("expected Count() 1, got %d", b.Count())
	}
}

func TestIntervalBookerAdjacentIntervalsDoNotOverlap(t *testing.T) {
	// Half-open semantics: [0,10) and [10,20) share only the boundary point.
	b := NewIntervalBooker()
	mustBook(t, b, 0, 10)
	mustBook(t, b, 10, 20)
	if b.Count() != 2 {
		t.Errorf("expected Count() 2, got %d", b.Count())
	}
}

func TestIntervalBookerPartialOverlapIsRejectedBothDirections(t *testing.T) {
	b := NewIntervalBooker()
	mustBook(t, b, 10, 20)
	if ok, err := b.Book(15, 25); err != nil || ok { // overlaps the tail
		t.Errorf("expected (false, nil), got (%v, %v)", ok, err)
	}
	if ok, err := b.Book(5, 15); err != nil || ok { // overlaps the head
		t.Errorf("expected (false, nil), got (%v, %v)", ok, err)
	}
	if b.Count() != 1 {
		t.Errorf("expected Count() 1, got %d", b.Count())
	}
}

func TestIntervalBookerContainmentIsRejectedBothDirections(t *testing.T) {
	b := NewIntervalBooker()
	mustBook(t, b, 10, 20)
	if ok, err := b.Book(12, 18); err != nil || ok { // inside the existing booking
		t.Errorf("expected (false, nil), got (%v, %v)", ok, err)
	}
	if ok, err := b.Book(5, 25); err != nil || ok { // swallows the existing booking
		t.Errorf("expected (false, nil), got (%v, %v)", ok, err)
	}
}

func TestIntervalBookerIdenticalIntervalIsRejected(t *testing.T) {
	b := NewIntervalBooker()
	mustBook(t, b, 10, 20)
	if ok, err := b.Book(10, 20); err != nil || ok {
		t.Errorf("expected (false, nil), got (%v, %v)", ok, err)
	}
}

func TestIntervalBookerRejectedBookChangesNothing(t *testing.T) {
	b := NewIntervalBooker()
	mustBook(t, b, 10, 20)
	if ok, err := b.Book(15, 25); err != nil || ok {
		t.Errorf("expected (false, nil), got (%v, %v)", ok, err)
	}
	if b.Count() != 1 {
		t.Errorf("expected Count() 1, got %d", b.Count())
	}
	mustBook(t, b, 20, 30) // the slot the rejected book didn't take
}

func TestIntervalBookerEmptyOrReversedIntervalErrors(t *testing.T) {
	b := NewIntervalBooker()
	if _, err := b.Book(10, 10); !errors.Is(err, ErrInvalidInterval) {
		t.Errorf("expected ErrInvalidInterval, got %v", err)
	}
	if _, err := b.Book(20, 10); !errors.Is(err, ErrInvalidInterval) {
		t.Errorf("expected ErrInvalidInterval, got %v", err)
	}
	if b.Count() != 0 {
		t.Errorf("expected Count() 0, got %d", b.Count())
	}
}

func TestIntervalBookerCancelExactMatchFreesTheSlot(t *testing.T) {
	b := NewIntervalBooker()
	mustBook(t, b, 10, 20)
	if !b.Cancel(10, 20) {
		t.Error("expected Cancel to return true for an exact match")
	}
	if b.Count() != 0 {
		t.Errorf("expected Count() 0, got %d", b.Count())
	}
	mustBook(t, b, 15, 25) // the freed range is bookable again
}

func TestIntervalBookerCancelRequiresExactMatch(t *testing.T) {
	b := NewIntervalBooker()
	mustBook(t, b, 10, 20)
	if b.Cancel(10, 15) { // sub-range: no
		t.Error("expected Cancel to return false for a sub-range")
	}
	if b.Cancel(5, 25) { // super-range: no
		t.Error("expected Cancel to return false for a super-range")
	}
	if b.Count() != 1 {
		t.Errorf("expected Count() 1, got %d", b.Count())
	}
}

func TestIntervalBookerCancelUnknownIntervalReturnsFalse(t *testing.T) {
	b := NewIntervalBooker()
	if b.Cancel(10, 20) {
		t.Error("expected Cancel to return false for an unknown interval")
	}
}

func TestIntervalBookerFirstFreeOnEmptyCalendarIsFrom(t *testing.T) {
	b := NewIntervalBooker()
	if s, err := b.FirstFree(7, 5); err != nil || s != 7 {
		t.Errorf("expected (7, nil), got (%d, %v)", s, err)
	}
}

func TestIntervalBookerFirstFreeSkipsBusyRangesAndTooSmallGaps(t *testing.T) {
	b := NewIntervalBooker()
	b.Book(10, 20)
	b.Book(20, 30)
	// [5,15) collides; the only room left starts after the last booking.
	if s, err := b.FirstFree(5, 10); err != nil || s != 30 {
		t.Errorf("expected (30, nil), got (%d, %v)", s, err)
	}
}

func TestIntervalBookerFirstFreeFindsGapOfExactSize(t *testing.T) {
	b := NewIntervalBooker()
	b.Book(0, 10)
	b.Book(15, 25)
	if s, err := b.FirstFree(0, 5); err != nil || s != 10 { // the gap [10,15) fits exactly
		t.Errorf("expected (10, nil), got (%d, %v)", s, err)
	}
}

func TestIntervalBookerFirstFreeWhenFromFallsInsideABooking(t *testing.T) {
	b := NewIntervalBooker()
	b.Book(10, 20)
	if s, err := b.FirstFree(12, 5); err != nil || s != 20 {
		t.Errorf("expected (20, nil), got (%d, %v)", s, err)
	}
}

func TestIntervalBookerFirstFreeIgnoresBookingsEntirelyBeforeFrom(t *testing.T) {
	b := NewIntervalBooker()
	b.Book(0, 10)
	if s, err := b.FirstFree(50, 5); err != nil || s != 50 {
		t.Errorf("expected (50, nil), got (%d, %v)", s, err)
	}
}

func TestIntervalBookerFirstFreeResultIsActuallyBookable(t *testing.T) {
	b := NewIntervalBooker()
	b.Book(3, 9)
	b.Book(12, 40)
	b.Book(41, 50)
	s, err := b.FirstFree(0, 4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok, err := b.Book(s, s+4); err != nil || !ok {
		t.Errorf("expected Book(%d, %d) to succeed, got (%v, %v)", s, s+4, ok, err)
	}
}

func TestIntervalBookerNonPositiveDurationErrors(t *testing.T) {
	b := NewIntervalBooker()
	if _, err := b.FirstFree(0, 0); !errors.Is(err, ErrInvalidDuration) {
		t.Errorf("expected ErrInvalidDuration, got %v", err)
	}
	if _, err := b.FirstFree(0, -3); !errors.Is(err, ErrInvalidDuration) {
		t.Errorf("expected ErrInvalidDuration, got %v", err)
	}
}
