package day03

import (
	"errors"
	"runtime"
	"slices"
	"sync/atomic"
	"testing"
	"time"
)

// The single-goroutine contract is the same as Day 2 -- a quick check
// that it still holds.

func TestBoundedStackPushPopIsLifo(t *testing.T) {
	watchdog(t)
	s := NewBoundedStack(3)
	must(t, s.Push(1))
	must(t, s.Push(2))
	must(t, s.Push(3))
	if !s.Full() {
		t.Error("expected Full() true")
	}
	if v, err := s.Top(); err != nil || v != 3 {
		t.Errorf("Top() = (%d, %v), want (3, nil)", v, err)
	}
	if v, err := s.Pop(); err != nil || v != 3 {
		t.Errorf("Pop() = (%d, %v), want (3, nil)", v, err)
	}
	if v, err := s.Pop(); err != nil || v != 2 {
		t.Errorf("Pop() = (%d, %v), want (2, nil)", v, err)
	}
	if got := s.Size(); got != 1 {
		t.Errorf("Size() = %d, want 1", got)
	}
}

func TestBoundedStackFullAndEmptyStillError(t *testing.T) {
	watchdog(t)
	s := NewBoundedStack(1)
	if !s.Empty() {
		t.Error("expected Empty() true")
	}
	if _, err := s.Pop(); !errors.Is(err, ErrEmpty) {
		t.Errorf("Pop() on an empty stack: got error %v, want ErrEmpty", err)
	}
	if _, err := s.Top(); !errors.Is(err, ErrEmpty) {
		t.Errorf("Top() on an empty stack: got error %v, want ErrEmpty", err)
	}
	must(t, s.Push(7))
	if err := s.Push(8); !errors.Is(err, ErrFull) {
		t.Errorf("Push(8) on a full stack: got error %v, want ErrFull", err)
	}
	if got := s.Size(); got != 1 {
		t.Errorf("Size() = %d, want 1", got)
	}
	if v, err := s.Top(); err != nil || v != 7 {
		t.Errorf("Top() = (%d, %v), want (7, nil)", v, err)
	}
}

func TestBoundedStackConcurrentPushesStopExactlyAtCapacity(t *testing.T) {
	watchdog(t)
	// 8 goroutines push 2,000 values into a stack that holds 1,000.
	// Exactly 1,000 pushes may succeed, and the stack must then hold
	// exactly the values whose push succeeded. Pushing past the capacity
	// can only happen at the moment the stack fills up -- so one round
	// can miss it, and the test runs 20.
	const rounds = 20
	const goroutines = 8
	const perGoroutine = 250
	for round := 0; round < rounds; round++ {
		s := NewBoundedStack(1000)
		pushed := make([][]int, goroutines)
		runTogether(t, goroutines, func(g int) {
			for i := 0; i < perGoroutine; i++ {
				v := g*1000 + i
				if err := s.Push(v); err == nil {
					pushed[g] = append(pushed[g], v)
				}
			}
		})

		var expected []int
		for _, mine := range pushed {
			expected = append(expected, mine...)
		}
		if len(expected) != 1000 {
			t.Fatalf("round %d: %d pushes succeeded, want exactly 1000 (the capacity)", round, len(expected))
		}
		if got := s.Size(); got != 1000 {
			t.Fatalf("round %d: Size() = %d, want 1000", round, got)
		}
		if !s.Full() {
			t.Fatalf("round %d: expected Full() true", round)
		}

		var held []int
		for i := 0; i < 1000 && !s.Empty(); i++ {
			v, err := s.Pop()
			if err != nil {
				t.Fatalf("round %d: Pop() on a non-empty stack: unexpected error: %v", round, err)
			}
			held = append(held, v)
		}
		if !s.Empty() {
			t.Fatalf("round %d: expected Empty() true after popping 1000 values", round)
		}
		slices.Sort(expected)
		slices.Sort(held)
		if !slices.Equal(held, expected) {
			t.Fatalf("round %d: the stack lost a value or made one up: it held %d values, "+
				"and they are not the %d values whose Push succeeded", round, len(held), len(expected))
		}
	}
}

func TestBoundedStackConcurrentPopsReturnEachValueOnce(t *testing.T) {
	watchdog(t)
	const values = 20000
	const goroutines = 8
	s := NewBoundedStack(values)
	for v := 0; v < values; v++ {
		must(t, s.Push(v))
	}

	// Pop until the stack says it is empty. (Stop early if more values
	// come out than ever went in -- the stack is making them up.)
	popped := make([][]int, goroutines)
	var poppedCount atomic.Int32
	runTogether(t, goroutines, func(g int) {
		for poppedCount.Load() <= values {
			v, err := s.Pop()
			if err != nil {
				return
			}
			popped[g] = append(popped[g], v)
			poppedCount.Add(1)
		}
	})

	var all []int
	for _, mine := range popped {
		all = append(all, mine...)
	}
	slices.Sort(all)
	if len(all) != values {
		t.Fatalf("%d values were popped, want %d -- a value was popped twice, or lost", len(all), values)
	}
	for v := 0; v < values; v++ {
		if all[v] != v {
			t.Fatalf("value %d is missing or was popped twice (sorted pops: position %d holds %d)", v, v, all[v])
		}
	}
	if !s.Empty() {
		t.Error("expected Empty() true")
	}
}

func TestBoundedStackConcurrentPushAndPopConserveValues(t *testing.T) {
	watchdog(t)
	// A small stack, 4 pushers and 4 poppers. Pushers retry while it is
	// full, poppers retry while it is empty. Every value pushed must come
	// out exactly once.
	const pushers = 4
	const perPusher = 5000
	const total = pushers * perPusher
	s := NewBoundedStack(16)
	var poppedCount atomic.Int32
	popped := make([][]int, pushers)
	deadline := time.Now().Add(20 * time.Second)
	pastDeadline := func() bool { return time.Now().After(deadline) }

	runTogether(t, 2*pushers, func(g int) {
		if g < pushers {
			for i := 0; i < perPusher && !pastDeadline(); i++ {
				for s.Push(g*perPusher+i) != nil {
					if pastDeadline() {
						return
					}
					runtime.Gosched()
				}
			}
			return
		}
		mine := &popped[g-pushers]
		for poppedCount.Load() < total && !pastDeadline() {
			v, err := s.Pop()
			if err != nil {
				runtime.Gosched()
				continue
			}
			*mine = append(*mine, v)
			poppedCount.Add(1)
		}
	})

	if pastDeadline() {
		t.Fatal("gave up after 20 s -- values were lost")
	}
	var all []int
	for _, mine := range popped {
		all = append(all, mine...)
	}
	slices.Sort(all)
	if len(all) != total {
		t.Fatalf("%d values were popped, want %d", len(all), total)
	}
	for v := 0; v < total; v++ {
		if all[v] != v {
			t.Fatalf("value %d is missing or was popped twice (sorted pops: position %d holds %d)", v, v, all[v])
		}
	}
	if !s.Empty() {
		t.Error("expected Empty() true")
	}
}

// ---- TryPush / TryPop: check and act as one step ----

func TestBoundedStackTryPushAndTryPopFollowTheRules(t *testing.T) {
	watchdog(t)
	s := NewBoundedStack(2)
	if v, ok := s.TryPop(); ok {
		t.Fatalf("TryPop() on an empty stack = (%d, true), want (0, false)", v)
	}
	if !s.TryPush(1) || !s.TryPush(2) {
		t.Fatal("TryPush on a stack with room returned false")
	}
	if s.TryPush(3) {
		t.Fatal("TryPush on a full stack returned true")
	}
	if s.Size() != 2 {
		t.Fatalf("Size() = %d, want 2", s.Size())
	}
	for _, want := range []int{2, 1} {
		if v, ok := s.TryPop(); !ok || v != want {
			t.Fatalf("TryPop() = (%d, %v), want (%d, true)", v, ok, want)
		}
	}
	if _, ok := s.TryPop(); ok {
		t.Fatal("TryPop() on an emptied stack returned true")
	}
	if !s.Empty() {
		t.Error("expected Empty() true")
	}
}

func TestBoundedStackConcurrentTryPopTakesEachValueOnce(t *testing.T) {
	watchdog(t)
	// 8 goroutines empty the stack with TryPop -- the one-call version of
	// "if !s.Empty() { v, _ := s.Top(); s.Pop() }".
	const values = 20000
	const goroutines = 8
	s := NewBoundedStack(values)
	for v := 0; v < values; v++ {
		must(t, s.Push(v))
	}

	popped := make([][]int, goroutines)
	var poppedCount atomic.Int32
	runTogether(t, goroutines, func(g int) {
		for poppedCount.Load() <= values { // more than values: made up
			v, ok := s.TryPop()
			if !ok {
				return
			}
			popped[g] = append(popped[g], v)
			poppedCount.Add(1)
		}
	})

	var all []int
	for _, mine := range popped {
		all = append(all, mine...)
	}
	slices.Sort(all)
	if len(all) != values {
		t.Fatalf("%d values were popped, want %d -- a value was popped twice, or lost", len(all), values)
	}
	for v := 0; v < values; v++ {
		if all[v] != v {
			t.Fatalf("value %d is missing or was popped twice (sorted pops: position %d holds %d)", v, v, all[v])
		}
	}
	if !s.Empty() {
		t.Error("expected Empty() true")
	}
}

func TestBoundedStackConcurrentTryPushStopsExactlyAtCapacity(t *testing.T) {
	watchdog(t)
	// Like TestBoundedStackConcurrentPushesStopExactlyAtCapacity, with
	// TryPush.
	const rounds = 20
	const goroutines = 8
	const perGoroutine = 250
	for round := 0; round < rounds; round++ {
		s := NewBoundedStack(1000)
		pushed := make([][]int, goroutines)
		runTogether(t, goroutines, func(g int) {
			for i := 0; i < perGoroutine; i++ {
				v := g*1000 + i
				if s.TryPush(v) {
					pushed[g] = append(pushed[g], v)
				}
			}
		})

		var expected []int
		for _, mine := range pushed {
			expected = append(expected, mine...)
		}
		if len(expected) != 1000 {
			t.Fatalf("round %d: %d TryPush calls returned true, want exactly 1000 (the capacity)", round, len(expected))
		}
		if got := s.Size(); got != 1000 {
			t.Fatalf("round %d: Size() = %d, want 1000", round, got)
		}
		var held []int
		for i := 0; i < 1000; i++ {
			v, ok := s.TryPop()
			if !ok {
				break
			}
			held = append(held, v)
		}
		slices.Sort(expected)
		slices.Sort(held)
		if !slices.Equal(held, expected) {
			t.Fatalf("round %d: the stack lost a value or made one up", round)
		}
	}
}

func TestBoundedStackConcurrentTryPushAndTryPopConserveValues(t *testing.T) {
	watchdog(t)
	// A small stack, 4 pushers and 4 poppers, all using the Try calls.
	const pushers = 4
	const perPusher = 5000
	const total = pushers * perPusher
	s := NewBoundedStack(16)
	var poppedCount atomic.Int32
	popped := make([][]int, pushers)
	deadline := time.Now().Add(10 * time.Second)
	pastDeadline := func() bool { return time.Now().After(deadline) }

	runTogether(t, 2*pushers, func(g int) {
		if g < pushers {
			for i := 0; i < perPusher && !pastDeadline(); i++ {
				for !s.TryPush(g*perPusher + i) {
					if pastDeadline() {
						return
					}
					runtime.Gosched()
				}
			}
			return
		}
		mine := &popped[g-pushers]
		for poppedCount.Load() < total && !pastDeadline() {
			v, ok := s.TryPop()
			if !ok {
				runtime.Gosched()
				continue
			}
			*mine = append(*mine, v)
			poppedCount.Add(1)
		}
	})

	if pastDeadline() {
		t.Fatal("gave up after 10 s -- values were lost")
	}
	var all []int
	for _, mine := range popped {
		all = append(all, mine...)
	}
	slices.Sort(all)
	if len(all) != total {
		t.Fatalf("%d values were popped, want %d", len(all), total)
	}
	for v := 0; v < total; v++ {
		if all[v] != v {
			t.Fatalf("value %d is missing or was popped twice (sorted pops: position %d holds %d)", v, v, all[v])
		}
	}
	if !s.Empty() {
		t.Error("expected Empty() true")
	}
}
