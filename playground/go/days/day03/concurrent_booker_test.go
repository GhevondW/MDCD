package day03

import (
	"cmp"
	"errors"
	"fmt"
	"math/rand"
	"slices"
	"sync/atomic"
	"testing"
)

// noOverlaps reports whether no two bookings in a sorted list overlap.
func noOverlaps(sorted []Interval) bool {
	for i := 1; i < len(sorted); i++ {
		if sorted[i].Start < sorted[i-1].End {
			return false
		}
	}
	return true
}

func byStart(a, b Interval) int { return cmp.Compare(a.Start, b.Start) }

// diffIntervals describes how got differs from want ("" if it doesn't).
func diffIntervals(got, want []Interval) string {
	for i := 0; i < len(got) && i < len(want); i++ {
		if got[i] != want[i] {
			return fmt.Sprintf("%d bookings, want %d; first difference at index %d: got %v, want %v",
				len(got), len(want), i, got[i], want[i])
		}
	}
	if len(got) != len(want) {
		return fmt.Sprintf("%d bookings, want %d", len(got), len(want))
	}
	return ""
}

// expectBook checks that Book(start, end) returns (want, nil).
func expectBook(t *testing.T, b *ConcurrentBooker, start, end int64, want bool) {
	t.Helper()
	if ok, err := b.Book(start, end); err != nil || ok != want {
		t.Errorf("Book(%d, %d) = (%v, %v), want (%v, nil)", start, end, ok, err, want)
	}
}

// expectBookFirstFree checks that BookFirstFree(from, duration) returns
// (want, nil).
func expectBookFirstFree(t *testing.T, b *ConcurrentBooker, from, duration, want int64) {
	t.Helper()
	if s, err := b.BookFirstFree(from, duration); err != nil || s != want {
		t.Errorf("BookFirstFree(%d, %d) = (%d, %v), want (%d, nil)", from, duration, s, err, want)
	}
}

func expectBookings(t *testing.T, b *ConcurrentBooker, want []Interval) {
	t.Helper()
	if got := b.Bookings(); !slices.Equal(got, want) {
		t.Errorf("Bookings() = %v, want %v", got, want)
	}
}

func expectCount(t *testing.T, b *ConcurrentBooker, want int) {
	t.Helper()
	if got := b.Count(); got != want {
		t.Errorf("Count() = %d, want %d", got, want)
	}
}

// ---- one goroutine: the rules ----

func TestConcurrentBookerBookRejectsOverlapAndAcceptsAdjacent(t *testing.T) {
	watchdog(t)
	b := NewConcurrentBooker()
	expectBook(t, b, 10, 20, true)
	expectBook(t, b, 20, 30, true) // adjacent: [10,20) and [20,30) don't overlap
	expectBook(t, b, 15, 25, false)
	expectBook(t, b, 0, 11, false)
	expectCount(t, b, 2)
}

func TestConcurrentBookerEmptyOrReversedIntervalErrors(t *testing.T) {
	watchdog(t)
	b := NewConcurrentBooker()
	if _, err := b.Book(5, 5); !errors.Is(err, ErrInvalidInterval) {
		t.Errorf("Book(5, 5): got error %v, want ErrInvalidInterval", err)
	}
	if _, err := b.Book(6, 5); !errors.Is(err, ErrInvalidInterval) {
		t.Errorf("Book(6, 5): got error %v, want ErrInvalidInterval", err)
	}
	expectCount(t, b, 0)
}

func TestConcurrentBookerCancelNeedsAnExactMatch(t *testing.T) {
	watchdog(t)
	b := NewConcurrentBooker()
	if ok, err := b.Book(10, 20); err != nil || !ok {
		t.Fatalf("Book(10, 20) = (%v, %v), want (true, nil)", ok, err)
	}
	if b.Cancel(10, 15) {
		t.Error("Cancel(10, 15) = true, want false: not an exact match")
	}
	if !b.Cancel(10, 20) {
		t.Error("Cancel(10, 20) = false, want true")
	}
	if b.Cancel(10, 20) {
		t.Error("Cancel(10, 20) a second time = true, want false")
	}
	expectCount(t, b, 0)
}

func TestConcurrentBookerBookingsAreSortedByStart(t *testing.T) {
	watchdog(t)
	b := NewConcurrentBooker()
	expectBook(t, b, 50, 60, true)
	expectBook(t, b, 0, 10, true)
	expectBook(t, b, 20, 30, true)
	expectBookings(t, b, []Interval{{0, 10}, {20, 30}, {50, 60}})
}

func TestConcurrentBookerBookFirstFreeOnEmptyCalendarStartsAtFrom(t *testing.T) {
	watchdog(t)
	b := NewConcurrentBooker()
	expectBookFirstFree(t, b, 7, 5, 7)
	expectBookings(t, b, []Interval{{7, 12}})
}

func TestConcurrentBookerBookFirstFreeSkipsBusyRangesAndSmallGaps(t *testing.T) {
	watchdog(t)
	b := NewConcurrentBooker()
	expectBook(t, b, 0, 10, true)
	expectBook(t, b, 12, 20, true)
	expectBookFirstFree(t, b, 0, 5, 20) // the gap [10,12) is too small
	expectBookFirstFree(t, b, 0, 2, 10) // ...but fits a booking of 2
	expectCount(t, b, 4)
}

func TestConcurrentBookerBookFirstFreeWhenFromIsInsideABooking(t *testing.T) {
	watchdog(t)
	b := NewConcurrentBooker()
	expectBook(t, b, 10, 20, true)
	expectBookFirstFree(t, b, 15, 5, 20)
}

func TestConcurrentBookerBookFirstFreeRejectsNonPositiveDuration(t *testing.T) {
	watchdog(t)
	b := NewConcurrentBooker()
	if _, err := b.BookFirstFree(0, 0); !errors.Is(err, ErrInvalidDuration) {
		t.Errorf("BookFirstFree(0, 0): got error %v, want ErrInvalidDuration", err)
	}
	if _, err := b.BookFirstFree(0, -1); !errors.Is(err, ErrInvalidDuration) {
		t.Errorf("BookFirstFree(0, -1): got error %v, want ErrInvalidDuration", err)
	}
	expectCount(t, b, 0)
}

func TestConcurrentBookerBookAllBooksEveryInterval(t *testing.T) {
	watchdog(t)
	b := NewConcurrentBooker()
	group := []Interval{{30, 40}, {0, 10}, {10, 20}}
	if ok, err := b.BookAll(group); err != nil || !ok {
		t.Errorf("BookAll(%v) = (%v, %v), want (true, nil)", group, ok, err)
	}
	expectBookings(t, b, []Interval{{0, 10}, {10, 20}, {30, 40}})
}

func TestConcurrentBookerBookAllIsAllOrNothing(t *testing.T) {
	watchdog(t)
	b := NewConcurrentBooker()
	if ok, err := b.Book(10, 20); err != nil || !ok {
		t.Fatalf("Book(10, 20) = (%v, %v), want (true, nil)", ok, err)
	}
	group := []Interval{{0, 5}, {15, 25}, {30, 40}}
	if ok, err := b.BookAll(group); err != nil || ok {
		t.Errorf("BookAll(%v) = (%v, %v), want (false, nil): [15,25) overlaps [10,20)", group, ok, err)
	}
	if got := b.Count(); got != 1 {
		t.Errorf("Count() = %d, want 1 -- a failed BookAll must book nothing", got)
	}
	expectBook(t, b, 0, 5, true)
	expectBook(t, b, 30, 40, true)
}

func TestConcurrentBookerBookAllRejectsOverlapInsideTheList(t *testing.T) {
	watchdog(t)
	b := NewConcurrentBooker()
	group := []Interval{{0, 10}, {5, 15}}
	if ok, err := b.BookAll(group); err != nil || ok {
		t.Errorf("BookAll(%v) = (%v, %v), want (false, nil)", group, ok, err)
	}
	expectCount(t, b, 0)
}

func TestConcurrentBookerBookAllOfEmptyListSucceeds(t *testing.T) {
	watchdog(t)
	b := NewConcurrentBooker()
	if ok, err := b.BookAll([]Interval{}); err != nil || !ok {
		t.Errorf("BookAll of an empty list = (%v, %v), want (true, nil)", ok, err)
	}
	expectCount(t, b, 0)
}

func TestConcurrentBookerBookAllWithInvalidIntervalErrorsAndBooksNothing(t *testing.T) {
	watchdog(t)
	b := NewConcurrentBooker()
	if _, err := b.BookAll([]Interval{{0, 5}, {7, 7}}); !errors.Is(err, ErrInvalidInterval) {
		t.Errorf("BookAll([{0 5} {7 7}]): got error %v, want ErrInvalidInterval", err)
	}
	expectCount(t, b, 0)
}

// ---- many goroutines ----

func TestConcurrentBookerConcurrentBooksOfOverlappingSlotsHaveOneWinner(t *testing.T) {
	watchdog(t)
	// Goroutine g books [g, g+16). Every pair of these overlaps, so in
	// each round exactly one goroutine may win.
	const goroutines = 16
	for round := 0; round < 100; round++ {
		b := NewConcurrentBooker()
		var wins atomic.Int32
		runTogether(t, goroutines, func(g int) {
			if ok, _ := b.Book(int64(g), int64(g)+16); ok {
				wins.Add(1)
			}
		})
		if n := wins.Load(); n != 1 {
			t.Fatalf("round %d: %d goroutines booked an overlapping slot, want exactly 1", round, n)
		}
		if n := b.Count(); n != 1 {
			t.Fatalf("round %d: Count() = %d, want 1", round, n)
		}
	}
}

func TestConcurrentBookerConcurrentBookFirstFreeTilesTheCalendar(t *testing.T) {
	watchdog(t)
	// 8 goroutines each take the first free 10-unit slot 200 times.
	// Nobody cancels, so the slots must fill [0, 16000) with no gaps and
	// no overlaps: 0, 10, 20, ...
	const goroutines = 8
	const perGoroutine = 200
	b := NewConcurrentBooker()
	starts := make([][]int64, goroutines)
	var bad firstFailure
	runTogether(t, goroutines, func(g int) {
		for i := 0; i < perGoroutine; i++ {
			s, err := b.BookFirstFree(0, 10)
			if err != nil {
				bad.record("BookFirstFree(0, 10): unexpected error: %v", err)
				return
			}
			starts[g] = append(starts[g], s)
		}
	})

	bad.report(t)
	var all []int64
	for _, mine := range starts {
		all = append(all, mine...)
	}
	slices.Sort(all)
	if len(all) != goroutines*perGoroutine {
		t.Fatalf("%d slots were booked, want %d", len(all), goroutines*perGoroutine)
	}
	for i, s := range all {
		if s != int64(i)*10 {
			t.Fatalf("sorted, the starts should be 0, 10, 20, ... but starts[%d] = %d, want %d "+
				"-- two goroutines got the same slot, or a slot was skipped", i, s, i*10)
		}
	}
	booked := b.Bookings()
	if len(booked) != len(all) {
		t.Errorf("Bookings() has %d bookings, want %d", len(booked), len(all))
	}
	if !noOverlaps(booked) {
		t.Error("Bookings() contains overlapping bookings")
	}
}

func TestConcurrentBookerConcurrentBookAllHasOneWinnerAndNoPartialBookings(t *testing.T) {
	watchdog(t)
	// Every goroutine's group contains the shared slot [50,60), plus two
	// slots of its own. Exactly one group can win each round -- and the
	// losers must leave none of their own slots behind.
	const goroutines = 8
	for round := 0; round < 100; round++ {
		b := NewConcurrentBooker()
		won := make([]bool, goroutines)
		errs := make([]error, goroutines)
		runTogether(t, goroutines, func(g int) {
			mine := int64(100 + 20*g)
			won[g], errs[g] = b.BookAll([]Interval{{mine, mine + 10}, {50, 60}, {1000 + mine, 1010 + mine}})
		})

		winners, w := 0, -1
		for g := range won {
			if errs[g] != nil {
				t.Fatalf("round %d: BookAll: unexpected error: %v", round, errs[g])
			}
			if won[g] {
				winners++
				w = g
			}
		}
		if winners != 1 {
			t.Fatalf("round %d: %d BookAll calls succeeded, want exactly 1", round, winners)
		}
		mine := int64(100 + 20*w)
		want := []Interval{{50, 60}, {mine, mine + 10}, {1000 + mine, 1010 + mine}}
		if got := b.Bookings(); !slices.Equal(got, want) {
			t.Fatalf("round %d: Bookings() = %v, want %v -- a losing BookAll left a partial booking", round, got, want)
		}
	}
}

func TestConcurrentBookerAFailedBookAllIsNeverSeenHalfDone(t *testing.T) {
	watchdog(t)
	// [50,60) is taken, so BookAll([[0,10), [50,60)]) must always fail and
	// book nothing -- which leaves [0,10) free for the one goroutine that
	// books and cancels it over and over. A BookAll that books [0,10) first
	// and takes it back once [50,60) turns out to be taken is not
	// all-or-nothing: for a moment, that goroutine finds [0,10) taken.
	const groupGoroutines = 3
	const tries = 20000
	b := NewConcurrentBooker()
	if ok, err := b.Book(50, 60); !ok || err != nil {
		t.Fatalf("Book(50, 60) = (%v, %v), want (true, nil)", ok, err)
	}
	var groupWins, blocked atomic.Int32
	var bad firstFailure
	runTogether(t, groupGoroutines+1, func(g int) {
		for i := 0; i < tries; i++ {
			if g < groupGoroutines {
				ok, err := b.BookAll([]Interval{{0, 10}, {50, 60}})
				if err != nil {
					bad.record("BookAll: unexpected error: %v", err)
					return
				}
				if ok {
					groupWins.Add(1)
				}
			} else if ok, err := b.Book(0, 10); err != nil {
				bad.record("Book(0, 10): unexpected error: %v", err)
				return
			} else if ok {
				b.Cancel(0, 10)
			} else {
				blocked.Add(1)
			}
		}
	})
	bad.report(t)

	if n := groupWins.Load(); n != 0 {
		t.Errorf("BookAll succeeded %d times although [50,60) was taken", n)
	}
	if n := blocked.Load(); n != 0 {
		t.Errorf("Book(0, 10) failed %d times: a BookAll that failed held [0,10) for a moment", n)
	}
	if got, want := b.Bookings(), []Interval{{50, 60}}; !slices.Equal(got, want) {
		t.Errorf("Bookings() = %v, want %v", got, want)
	}
}

func TestConcurrentBookerConcurrentMixedOperationsKeepTheCalendarConsistent(t *testing.T) {
	watchdog(t)
	// Goroutines book, cancel and BookFirstFree at random. Afterwards the
	// calendar must hold exactly the bookings that were made and not
	// cancelled -- and none may overlap.
	const goroutines = 8
	const ops = 3000
	b := NewConcurrentBooker()
	kept := make([][]Interval, goroutines)
	var bad firstFailure
	runTogether(t, goroutines, func(g int) {
		rng := rand.New(rand.NewSource(int64(1234 + g)))
		mine := &kept[g]
		for i := 0; i < ops; i++ {
			switch rng.Intn(3) {
			case 0:
				s := rng.Int63n(5000)
				length := 1 + rng.Int63n(20)
				ok, err := b.Book(s, s+length)
				if err != nil {
					bad.record("Book(%d, %d): unexpected error: %v", s, s+length, err)
					return
				}
				if ok {
					*mine = append(*mine, Interval{s, s + length})
				}
			case 1:
				length := 1 + rng.Int63n(20)
				from := rng.Int63n(5000)
				s, err := b.BookFirstFree(from, length)
				if err != nil {
					bad.record("BookFirstFree(%d, %d): unexpected error: %v", from, length, err)
					return
				}
				*mine = append(*mine, Interval{s, s + length})
			default:
				if len(*mine) == 0 {
					continue
				}
				k := rng.Intn(len(*mine))
				iv := (*mine)[k]
				if !b.Cancel(iv.Start, iv.End) {
					bad.record("could not cancel %v, a booking this goroutine made", iv)
					return
				}
				*mine = slices.Delete(*mine, k, k+1)
			}
		}
	})

	bad.report(t)
	var expected []Interval
	for _, mine := range kept {
		expected = append(expected, mine...)
	}
	slices.SortFunc(expected, byStart)
	booked := b.Bookings()
	if !slices.IsSortedFunc(booked, byStart) {
		t.Error("Bookings() is not sorted by Start")
	}
	if !noOverlaps(booked) {
		t.Error("Bookings() contains overlapping bookings")
	}
	if d := diffIntervals(booked, expected); d != "" {
		t.Errorf("Bookings() is not exactly the bookings made and not cancelled: %s", d)
	}
	if got := b.Count(); got != len(expected) {
		t.Errorf("Count() = %d, want %d", got, len(expected))
	}
}
