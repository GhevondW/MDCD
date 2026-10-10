package day04

import (
	"sort"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// A goroutine that calls a blocking method, with its result delivered
// through a signal and a channel so the test can look at it later.
type popCall struct {
	fired *signal
	v     int
	ok    bool
}

func startPop(q *BoundedBlockingQueue[int]) *popCall {
	c := &popCall{fired: newSignal()}
	go func() {
		c.v, c.ok = q.Pop()
		c.fired.fire()
	}()
	return c
}

type pushCall struct {
	fired *signal
	ok    bool
}

func startPush(q *BoundedBlockingQueue[int], v int) *pushCall {
	c := &pushCall{fired: newSignal()}
	go func() {
		c.ok = q.Push(v)
		c.fired.fire()
	}()
	return c
}

// ---- single goroutine ----------------------------------------------------

func TestQueueFifoOrder(t *testing.T) {
	watchdog(t)
	q := NewBoundedBlockingQueue[int](3)
	if q.Cap() != 3 || q.Len() != 0 || q.Closed() {
		t.Fatalf("fresh queue: Cap=%d Len=%d Closed=%v, want 3 0 false", q.Cap(), q.Len(), q.Closed())
	}
	for _, v := range []int{1, 2, 3} {
		if !q.Push(v) {
			t.Fatalf("Push(%d) on a queue with room returned false", v)
		}
	}
	if q.Len() != 3 {
		t.Errorf("Len() = %d, want 3", q.Len())
	}
	for _, want := range []int{1, 2, 3} {
		if v, ok := q.Pop(); !ok || v != want {
			t.Errorf("Pop() = (%d, %v), want (%d, true)", v, ok, want)
		}
	}
	if q.Len() != 0 {
		t.Errorf("Len() = %d, want 0", q.Len())
	}
}

func TestQueueRejectsNonPositiveCapacity(t *testing.T) {
	watchdog(t)
	for _, c := range []int{0, -1} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("NewBoundedBlockingQueue(%d) did not panic", c)
				}
			}()
			NewBoundedBlockingQueue[int](c)
		}()
	}
}

func TestQueueWorksForAnyElementType(t *testing.T) {
	watchdog(t)
	q := NewBoundedBlockingQueue[string](2)
	q.Push("a")
	q.Push("b")
	if v, ok := q.Pop(); !ok || v != "a" {
		t.Errorf(`Pop() = (%q, %v), want ("a", true)`, v, ok)
	}
	p := NewBoundedBlockingQueue[*int](1)
	x := 5
	p.Push(&x)
	if v, ok := p.Pop(); !ok || v != &x {
		t.Errorf("Pop() did not return the pushed pointer")
	}
}

func TestQueueTryPushAndTryPopNeverWait(t *testing.T) {
	watchdog(t)
	q := NewBoundedBlockingQueue[int](2)
	if !within(callTimeLimit, func() {
		if _, ok := q.TryPop(); ok {
			t.Error("TryPop() on an empty queue returned true")
		}
		if !q.TryPush(1) || !q.TryPush(2) {
			t.Error("TryPush on a queue with room returned false")
		}
		if q.TryPush(3) {
			t.Error("TryPush on a full queue returned true")
		}
	}) {
		t.Fatal("TryPush / TryPop waited")
	}
	if q.Len() != 2 {
		t.Errorf("Len() = %d, want 2 (the rejected TryPush must change nothing)", q.Len())
	}
	if v, ok := q.TryPop(); !ok || v != 1 {
		t.Errorf("TryPop() = (%d, %v), want (1, true)", v, ok)
	}
}

// ---- waiting ---------------------------------------------------------------

func TestQueuePushWaitsWhileFull(t *testing.T) {
	watchdog(t)
	q := NewBoundedBlockingQueue[int](1)
	q.Push(1)
	c := startPush(q, 2)
	if c.fired.waitFor(150 * time.Millisecond) {
		t.Fatal("Push on a full queue returned without waiting")
	}
	if v, ok := q.Pop(); !ok || v != 1 {
		t.Fatalf("Pop() = (%d, %v), want (1, true)", v, ok)
	}
	if !c.fired.waitFor(5*time.Second) || !c.ok {
		t.Fatal("the waiting Push did not complete (or returned false) after a Pop made room")
	}
	if v, ok := q.Pop(); !ok || v != 2 {
		t.Errorf("Pop() = (%d, %v), want (2, true)", v, ok)
	}
}

func TestQueuePopWaitsWhileEmpty(t *testing.T) {
	watchdog(t)
	q := NewBoundedBlockingQueue[int](1)
	c := startPop(q)
	if c.fired.waitFor(150 * time.Millisecond) {
		t.Fatal("Pop on an empty queue returned without waiting")
	}
	q.Push(7)
	if !c.fired.waitFor(5 * time.Second) {
		t.Fatal("the waiting Pop did not wake up after a Push")
	}
	if !c.ok || c.v != 7 {
		t.Errorf("Pop() = (%d, %v), want (7, true)", c.v, c.ok)
	}
}

func TestQueueBackpressureStopsProducersAtCapacity(t *testing.T) {
	watchdog(t)
	q := NewBoundedBlockingQueue[int](4)
	var pushed atomic.Int32
	for p := 0; p < 3; p++ {
		go func() {
			for i := 0; i < 100; i++ {
				if !q.Push(i) {
					return
				}
				pushed.Add(1)
			}
		}()
	}
	time.Sleep(200 * time.Millisecond)
	if got := pushed.Load(); got != 4 {
		t.Fatalf("with nobody consuming, %d pushes completed; want exactly the capacity, 4", got)
	}
	q.Pop()
	deadline := time.Now().Add(5 * time.Second)
	for pushed.Load() != 5 && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	time.Sleep(50 * time.Millisecond)
	if got := pushed.Load(); got != 5 {
		t.Errorf("after one Pop, %d pushes completed; want 5", got)
	}
	q.Close()
}

// ---- close -----------------------------------------------------------------

func TestQueueCloseWakesWaitingConsumers(t *testing.T) {
	watchdog(t)
	q := NewBoundedBlockingQueue[int](2)
	calls := make([]*popCall, 5)
	for i := range calls {
		calls[i] = startPop(q)
	}
	time.Sleep(100 * time.Millisecond)
	q.Close()
	for i, c := range calls {
		if !c.fired.waitFor(5 * time.Second) {
			t.Fatalf("consumer %d is still waiting after Close() -- Close must wake EVERY sleeper (Broadcast)", i)
		}
		if c.ok {
			t.Errorf("consumer %d: Pop() on a closed, empty queue returned ok=true", i)
		}
	}
}

func TestQueueCloseWakesWaitingProducers(t *testing.T) {
	watchdog(t)
	q := NewBoundedBlockingQueue[int](1)
	q.Push(100)
	calls := make([]*pushCall, 5)
	for i := range calls {
		calls[i] = startPush(q, i)
	}
	time.Sleep(100 * time.Millisecond)
	q.Close()
	for i, c := range calls {
		if !c.fired.waitFor(5 * time.Second) {
			t.Fatalf("producer %d is still waiting after Close() -- Close must wake EVERY sleeper (Broadcast)", i)
		}
		if c.ok {
			t.Errorf("producer %d: Push returned true although the queue was closed", i)
		}
	}
	if q.Len() != 1 {
		t.Errorf("Len() = %d, want 1: a Push that returns false must not add its item", q.Len())
	}
	if v, ok := q.Pop(); !ok || v != 100 {
		t.Errorf("Pop() = (%d, %v), want (100, true): the item from before Close is still there", v, ok)
	}
	if _, ok := q.Pop(); ok {
		t.Error("a second Pop() returned an item that was never added")
	}
}

func TestQueueCloseWakesTimedWaits(t *testing.T) {
	watchdog(t)
	q := NewBoundedBlockingQueue[int](1)
	popDone, pushDone := newSignal(), newSignal()
	go func() { q.PopFor(30 * time.Second); popDone.fire() }()
	full := NewBoundedBlockingQueue[int](1)
	full.Push(1)
	go func() { full.PushFor(2, 30*time.Second); pushDone.fire() }()
	time.Sleep(100 * time.Millisecond)
	q.Close()
	full.Close()
	if !popDone.waitFor(5*time.Second) || !pushDone.waitFor(5*time.Second) {
		t.Fatal("PopFor / PushFor kept waiting after Close()")
	}
}

func TestQueueClosedQueueDrainsThenStops(t *testing.T) {
	watchdog(t)
	q := NewBoundedBlockingQueue[int](5)
	q.Push(1)
	q.Push(2)
	q.Push(3)
	q.Close()
	q.Close() // twice is fine
	if !q.Closed() {
		t.Error("Closed() = false after Close()")
	}
	if q.Push(4) || q.TryPush(4) || q.PushFor(4, 50*time.Millisecond) {
		t.Error("a push on a closed queue returned true")
	}
	for _, want := range []int{1, 2, 3} {
		if v, ok := q.Pop(); !ok || v != want {
			t.Errorf("Pop() after Close = (%d, %v), want (%d, true): consumers drain what is left", v, ok, want)
		}
	}
	if _, ok := q.Pop(); ok {
		t.Error("Pop() on a closed, drained queue returned ok=true")
	}
	if _, ok := q.TryPop(); ok {
		t.Error("TryPop() on a closed, drained queue returned ok=true")
	}
	if _, ok := q.PopFor(50 * time.Millisecond); ok {
		t.Error("PopFor on a closed, drained queue returned ok=true")
	}
}

func TestQueueCloseFromManyGoroutinesAtOnce(t *testing.T) {
	watchdog(t)
	q := NewBoundedBlockingQueue[int](2)
	runTogether(t, 8, func(g int) { q.Close() })
	if !q.Closed() {
		t.Error("Closed() = false")
	}
}

// ---- timeouts --------------------------------------------------------------

func TestQueuePopForTimesOut(t *testing.T) {
	watchdog(t)
	q := NewBoundedBlockingQueue[int](1)
	start := time.Now()
	_, ok := q.PopFor(100 * time.Millisecond)
	took := time.Since(start)
	if ok {
		t.Fatal("PopFor on an empty queue returned ok=true")
	}
	if took < 90*time.Millisecond {
		t.Errorf("PopFor(100ms) gave up after only %v", took)
	}
	if took > 3*time.Second {
		t.Errorf("PopFor(100ms) took %v", took)
	}
}

func TestQueuePushForTimesOut(t *testing.T) {
	watchdog(t)
	q := NewBoundedBlockingQueue[int](1)
	q.Push(1)
	start := time.Now()
	ok := q.PushFor(2, 100*time.Millisecond)
	took := time.Since(start)
	if ok {
		t.Fatal("PushFor on a full queue returned true")
	}
	if took < 90*time.Millisecond {
		t.Errorf("PushFor(100ms) gave up after only %v", took)
	}
	if took > 3*time.Second {
		t.Errorf("PushFor(100ms) took %v", took)
	}
	if q.Len() != 1 {
		t.Errorf("Len() = %d, want 1: a timed-out PushFor must not add its item", q.Len())
	}
}

func TestQueueTimedCallsReturnEarlyWhenTheyCan(t *testing.T) {
	watchdog(t)
	q := NewBoundedBlockingQueue[int](1)
	go func() { time.Sleep(50 * time.Millisecond); q.Push(9) }()
	start := time.Now()
	v, ok := q.PopFor(10 * time.Second)
	if !ok || v != 9 {
		t.Fatalf("PopFor = (%d, %v), want (9, true)", v, ok)
	}
	if took := time.Since(start); took > 3*time.Second {
		t.Errorf("PopFor(10s) took %v although the item arrived after 50ms", took)
	}

	go func() { time.Sleep(50 * time.Millisecond); q.Pop() }()
	q.Push(1) // full again
	start = time.Now()
	if !q.PushFor(2, 10*time.Second) {
		t.Fatal("PushFor returned false although room appeared after 50ms")
	}
	if took := time.Since(start); took > 3*time.Second {
		t.Errorf("PushFor(10s) took %v although room appeared after 50ms", took)
	}
}

func TestQueueTimeoutIsOneDeadlineNotRestartedOnEveryWakeUp(t *testing.T) {
	watchdog(t)
	// Another goroutine keeps taking the item out and putting it back every
	// few milliseconds, waking the waiter again and again without ever
	// leaving it room for long. A PopFor / PushFor that starts a fresh
	// timeout after each wake-up would wait for as long as this goes on.
	q := NewBoundedBlockingQueue[int](1)
	q.Push(1)
	stop := newSignal()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for !stop.isSet() {
			if v, ok := q.TryPop(); ok {
				q.TryPush(v)
			}
			time.Sleep(5 * time.Millisecond)
		}
	}()
	defer func() { stop.fire(); wg.Wait() }()

	start := time.Now()
	done := within(10*time.Second, func() { q.PushFor(2, 300*time.Millisecond) })
	took := time.Since(start)
	if !done || took > 1500*time.Millisecond {
		t.Errorf("PushFor(300ms) took %v while other goroutines kept waking it: the timeout must be one deadline for the whole call", took)
	}
}

// ---- many goroutines ---------------------------------------------------------

// stress: `producers` goroutines push perProducer distinct values each,
// `consumers` goroutines pop until the queue is closed and drained. The
// last producer to finish closes the queue.
func stress(t *testing.T, producers, consumers, perProducer, capacity int) {
	t.Helper()
	q := NewBoundedBlockingQueue[int](capacity)
	var fail firstFailure
	var producersLeft atomic.Int32
	producersLeft.Store(int32(producers))
	got := make([][]int, consumers)

	runTogether(t, producers+consumers, func(g int) {
		if g < producers {
			for i := 0; i < perProducer; i++ {
				if !q.Push(g*1_000_000 + i) {
					fail.record("producer %d: Push returned false on a queue nobody closed", g)
					break
				}
				if n := q.Len(); n > capacity {
					fail.record("Len() = %d, more than the capacity %d", n, capacity)
				}
			}
			if producersLeft.Add(-1) == 0 {
				q.Close()
			}
			return
		}
		c := g - producers
		for {
			v, ok := q.Pop()
			if !ok {
				return
			}
			got[c] = append(got[c], v)
		}
	})
	fail.report(t)

	var all []int
	for c, vs := range got {
		// FIFO: what one consumer sees from one producer is in push order.
		last := map[int]int{}
		for _, v := range vs {
			p, i := v/1_000_000, v%1_000_000
			if prev, seen := last[p]; seen && i < prev {
				t.Fatalf("consumer %d saw producer %d's item %d after item %d: not FIFO", c, p, i, prev)
			}
			last[p] = i
		}
		all = append(all, vs...)
	}
	if want := producers * perProducer; len(all) != want {
		t.Fatalf("consumers received %d items, want %d: the queue lost or invented items", len(all), want)
	}
	sort.Ints(all)
	for i := 1; i < len(all); i++ {
		if all[i] == all[i-1] {
			t.Fatalf("item %d was delivered twice", all[i])
		}
	}
	if q.Len() != 0 {
		t.Errorf("Len() = %d after everything was consumed", q.Len())
	}
}

func TestQueueOneProducerOneConsumerCapacityOne(t *testing.T) {
	watchdog(t)
	stress(t, 1, 1, 20000, 1)
}

func TestQueueManyProducersManyConsumers(t *testing.T) {
	watchdog(t)
	stress(t, 4, 4, 5000, 16)
}

func TestQueueMoreProducersThanConsumers(t *testing.T) {
	watchdog(t)
	stress(t, 8, 2, 3000, 3)
}

func TestQueueMoreConsumersThanProducers(t *testing.T) {
	watchdog(t)
	stress(t, 2, 8, 3000, 3)
}
