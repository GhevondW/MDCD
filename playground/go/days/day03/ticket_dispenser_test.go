package day03

import (
	"slices"
	"sync/atomic"
	"testing"
)

func TestTicketDispenserFirstTicketIsOne(t *testing.T) {
	watchdog(t)
	d := NewTicketDispenser()
	if got := d.Next(); got != 1 {
		t.Errorf("first Next() = %d, want 1", got)
	}
}

func TestTicketDispenserTicketsCountUpByOne(t *testing.T) {
	watchdog(t)
	d := NewTicketDispenser()
	for want := int64(1); want <= 3; want++ {
		if got := d.Next(); got != want {
			t.Errorf("Next() = %d, want %d", got, want)
		}
	}
}

func TestTicketDispenserIssuedCountsHandedOutTickets(t *testing.T) {
	watchdog(t)
	d := NewTicketDispenser()
	if got := d.Issued(); got != 0 {
		t.Errorf("Issued() before any Next() = %d, want 0", got)
	}
	d.Next()
	d.Next()
	if got := d.Issued(); got != 2 {
		t.Errorf("Issued() after two Next() = %d, want 2", got)
	}
}

func TestTicketDispenserConcurrentTicketsAreUniqueAndGapless(t *testing.T) {
	watchdog(t)
	const goroutines = 8
	const perGoroutine = 20000
	d := NewTicketDispenser()
	got := make([][]int64, goroutines)
	runTogether(t, goroutines, func(g int) {
		for i := 0; i < perGoroutine; i++ {
			got[g] = append(got[g], d.Next())
		}
	})

	var all []int64
	for _, mine := range got {
		all = append(all, mine...)
	}
	slices.Sort(all)
	for i, v := range all {
		if v != int64(i)+1 {
			t.Fatalf("sorted, the tickets should be 1, 2, 3, ... but tickets[%d] = %d, want %d "+
				"-- a ticket was handed out twice, or a number was skipped", i, v, i+1)
		}
	}
	if got := d.Issued(); got != goroutines*perGoroutine {
		t.Errorf("Issued() = %d, want %d", got, goroutines*perGoroutine)
	}
}

func TestTicketDispenserEachGoroutineSeesItsTicketsIncrease(t *testing.T) {
	watchdog(t)
	const goroutines = 8
	const perGoroutine = 20000
	d := NewTicketDispenser()
	got := make([][]int64, goroutines)
	runTogether(t, goroutines, func(g int) {
		for i := 0; i < perGoroutine; i++ {
			got[g] = append(got[g], d.Next())
		}
	})

	for g, mine := range got {
		for i := 1; i < len(mine); i++ {
			if mine[i-1] >= mine[i] {
				t.Fatalf("goroutine %d got ticket %d, then ticket %d "+
					"-- a later ticket had a smaller (or the same) number", g, mine[i-1], mine[i])
			}
		}
	}
}

func TestTicketDispenserIssuedNeverGoesBackwardsWhileTicketsAreTaken(t *testing.T) {
	watchdog(t)
	const takers = 4
	const perGoroutine = 20000
	d := NewTicketDispenser()
	var wentBackwards atomic.Bool
	runTogether(t, takers+1, func(g int) {
		if g == takers { // the watcher
			var last int64
			for i := 0; i < 20000; i++ {
				now := d.Issued()
				if now < last {
					wentBackwards.Store(true)
				}
				last = now
			}
			return
		}
		for i := 0; i < perGoroutine; i++ {
			d.Next()
		}
	})

	if wentBackwards.Load() {
		t.Error("Issued() went backwards while tickets were being taken")
	}
	if got := d.Issued(); got != takers*perGoroutine {
		t.Errorf("Issued() = %d, want %d", got, takers*perGoroutine)
	}
}
