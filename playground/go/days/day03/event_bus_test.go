package day03

import (
	"slices"
	"sync/atomic"
	"testing"
	"time"
)

// publishWithin calls bus.Publish(event) and fails the test if it has
// not returned after callTimeLimit -- a Publish that holds a lock while
// callbacks run deadlocks as soon as a callback calls back into the bus.
func publishWithin(t *testing.T, bus *EventBus, event string) int {
	t.Helper()
	var n int
	if !within(callTimeLimit, func() { n = bus.Publish(event) }) {
		t.Fatalf("Publish(%q) did not return within %v -- deadlocked? "+
			"Is a lock held while callbacks run?", event, callTimeLimit)
	}
	return n
}

func TestEventBusPublishCallsSubscribersInOrder(t *testing.T) {
	watchdog(t)
	bus := NewEventBus()
	var log []string
	bus.Subscribe(func(e string) { log = append(log, "A:"+e) })
	bus.Subscribe(func(e string) { log = append(log, "B:"+e) })
	bus.Subscribe(func(e string) { log = append(log, "C:"+e) })
	if n := bus.Publish("x"); n != 3 {
		t.Errorf("Publish(\"x\") = %d, want 3", n)
	}
	if want := []string{"A:x", "B:x", "C:x"}; !slices.Equal(log, want) {
		t.Errorf("callbacks ran as %q, want %q", log, want)
	}
}

func TestEventBusPublishWithNoSubscribersCallsNobody(t *testing.T) {
	watchdog(t)
	bus := NewEventBus()
	if n := bus.Publish("x"); n != 0 {
		t.Errorf("Publish(\"x\") = %d, want 0", n)
	}
	if n := bus.SubscriberCount(); n != 0 {
		t.Errorf("SubscriberCount() = %d, want 0", n)
	}
}

func TestEventBusIdsArePositiveAndUnique(t *testing.T) {
	watchdog(t)
	bus := NewEventBus()
	ids := map[int64]bool{}
	for i := 0; i < 100; i++ {
		id := bus.Subscribe(func(string) {})
		if id <= 0 {
			t.Fatalf("Subscribe returned id %d, want > 0", id)
		}
		ids[id] = true
	}
	if len(ids) != 100 {
		t.Errorf("100 Subscribe calls returned %d different ids, want 100", len(ids))
	}
	if n := bus.SubscriberCount(); n != 100 {
		t.Errorf("SubscriberCount() = %d, want 100", n)
	}
}

func TestEventBusUnsubscribeStopsDelivery(t *testing.T) {
	watchdog(t)
	bus := NewEventBus()
	var log []string
	bus.Subscribe(func(string) { log = append(log, "A") })
	b := bus.Subscribe(func(string) { log = append(log, "B") })
	bus.Subscribe(func(string) { log = append(log, "C") })

	if !bus.Unsubscribe(b) {
		t.Error("Unsubscribe(b) = false, want true")
	}
	if bus.Unsubscribe(b) {
		t.Error("Unsubscribe(b) a second time = true, want false: it was already removed")
	}
	if bus.Unsubscribe(12345) {
		t.Error("Unsubscribe(12345) = true, want false: that id never existed")
	}
	if n := bus.SubscriberCount(); n != 2 {
		t.Errorf("SubscriberCount() = %d, want 2", n)
	}
	if n := bus.Publish("x"); n != 2 {
		t.Errorf("Publish(\"x\") = %d, want 2", n)
	}
	if want := []string{"A", "C"}; !slices.Equal(log, want) {
		t.Errorf("callbacks ran as %q, want %q", log, want)
	}
}

func TestEventBusOrderIsKeptAfterUnsubscribeAndSubscribe(t *testing.T) {
	watchdog(t)
	bus := NewEventBus()
	var log []string
	bus.Subscribe(func(string) { log = append(log, "A") })
	b := bus.Subscribe(func(string) { log = append(log, "B") })
	bus.Subscribe(func(string) { log = append(log, "C") })
	bus.Unsubscribe(b)
	bus.Subscribe(func(string) { log = append(log, "D") })
	bus.Publish("x")
	if want := []string{"A", "C", "D"}; !slices.Equal(log, want) {
		t.Errorf("callbacks ran as %q, want %q", log, want)
	}
}

func TestEventBusSubscribeInsideACallbackCountsFromTheNextPublish(t *testing.T) {
	watchdog(t)
	bus := NewEventBus()
	var log []string
	added := false
	bus.Subscribe(func(string) {
		log = append(log, "A")
		if !added {
			added = true
			bus.Subscribe(func(string) { log = append(log, "B") })
		}
	})

	if n := publishWithin(t, bus, "first"); n != 1 {
		t.Errorf("Publish(\"first\") = %d, want 1 -- B joined mid-publish, so this publish must not call it", n)
	}
	if want := []string{"A"}; !slices.Equal(log, want) {
		t.Errorf("first publish: callbacks ran as %q, want %q", log, want)
	}
	log = nil
	if n := publishWithin(t, bus, "second"); n != 2 {
		t.Errorf("Publish(\"second\") = %d, want 2", n)
	}
	if want := []string{"A", "B"}; !slices.Equal(log, want) {
		t.Errorf("second publish: callbacks ran as %q, want %q", log, want)
	}
}

func TestEventBusUnsubscribeInsideACallbackCountsFromTheNextPublish(t *testing.T) {
	watchdog(t)
	bus := NewEventBus()
	var log []string
	var b int64
	removed := false
	unsubscribed := false // what Unsubscribe(b) returned inside the callback
	bus.Subscribe(func(string) {
		log = append(log, "A")
		if !removed {
			removed = true
			unsubscribed = bus.Unsubscribe(b)
		}
	})
	b = bus.Subscribe(func(string) { log = append(log, "B") })

	if n := publishWithin(t, bus, "first"); n != 2 {
		t.Errorf("Publish(\"first\") = %d, want 2 -- B left mid-publish, so this publish must still call it", n)
	}
	if !unsubscribed {
		t.Error("Unsubscribe(b) inside a callback returned false, want true")
	}
	if want := []string{"A", "B"}; !slices.Equal(log, want) {
		t.Errorf("first publish: callbacks ran as %q, want %q", log, want)
	}
	log = nil
	if n := publishWithin(t, bus, "second"); n != 1 {
		t.Errorf("Publish(\"second\") = %d, want 1", n)
	}
	if want := []string{"A"}; !slices.Equal(log, want) {
		t.Errorf("second publish: callbacks ran as %q, want %q", log, want)
	}
}

func TestEventBusCallbackCanUnsubscribeItself(t *testing.T) {
	watchdog(t)
	bus := NewEventBus()
	aCalls := 0
	var a int64
	a = bus.Subscribe(func(string) {
		aCalls++
		bus.Unsubscribe(a)
	})
	bus.Subscribe(func(string) {})

	if n := publishWithin(t, bus, "first"); n != 2 {
		t.Errorf("Publish(\"first\") = %d, want 2", n)
	}
	if n := publishWithin(t, bus, "second"); n != 1 {
		t.Errorf("Publish(\"second\") = %d, want 1", n)
	}
	if aCalls != 1 {
		t.Errorf("the self-removing callback ran %d times, want 1", aCalls)
	}
}

func TestEventBusCallbackCanPublish(t *testing.T) {
	watchdog(t)
	bus := NewEventBus()
	var bLog []string
	bus.Subscribe(func(e string) {
		if e == "outer" {
			bus.Publish("inner")
		}
	})
	bus.Subscribe(func(e string) { bLog = append(bLog, e) })

	if n := publishWithin(t, bus, "outer"); n != 2 {
		t.Errorf("Publish(\"outer\") = %d, want 2", n)
	}
	// A handles "outer" by publishing "inner", which reaches B first.
	if want := []string{"inner", "outer"}; !slices.Equal(bLog, want) {
		t.Errorf("B heard %q, want %q", bLog, want)
	}
}

func TestEventBusOtherGoroutinesCanUseTheBusWhileACallbackRuns(t *testing.T) {
	watchdog(t)
	// The callback waits (up to 2 s) for another goroutine to subscribe.
	// If Publish holds a lock while calling back, that goroutine is stuck.
	bus := NewEventBus()
	outcome := ""
	var helperDone *signal // set when the callback starts the helper
	bus.Subscribe(func(string) {
		done := newSignal()
		helperDone = done
		go func() {
			bus.Subscribe(func(string) {})
			done.fire()
		}()
		if done.waitFor(2 * time.Second) {
			outcome = "ok"
		} else {
			outcome = "blocked"
		}
	})

	publishWithin(t, bus, "x")
	if helperDone != nil && !helperDone.waitFor(callTimeLimit) {
		t.Fatalf("the other goroutine's Subscribe had not returned %v after Publish returned", callTimeLimit)
	}
	if outcome != "ok" {
		t.Errorf("outcome = %q, want \"ok\" -- another goroutine could not subscribe while a callback ran", outcome)
	}
	if n := bus.SubscriberCount(); n != 2 {
		t.Errorf("SubscriberCount() = %d, want 2", n)
	}
}

func TestEventBusConcurrentSubscribesGetUniqueIds(t *testing.T) {
	watchdog(t)
	const goroutines = 8
	const perGoroutine = 1000
	bus := NewEventBus()
	ids := make([][]int64, goroutines)
	runTogether(t, goroutines, func(g int) {
		for i := 0; i < perGoroutine; i++ {
			ids[g] = append(ids[g], bus.Subscribe(func(string) {}))
		}
	})

	all := map[int64]bool{}
	for _, mine := range ids {
		for _, id := range mine {
			all[id] = true
		}
	}
	if len(all) != goroutines*perGoroutine {
		t.Errorf("%d Subscribe calls returned %d different ids, want %d",
			goroutines*perGoroutine, len(all), goroutines*perGoroutine)
	}
	if n := bus.SubscriberCount(); n != goroutines*perGoroutine {
		t.Errorf("SubscriberCount() = %d, want %d", n, goroutines*perGoroutine)
	}
}

func TestEventBusConcurrentPublishAndChurn(t *testing.T) {
	watchdog(t)
	// 3 subscribers stay for the whole test and must hear every event
	// exactly once, while other goroutines keep subscribing and
	// unsubscribing around them.
	const publishers = 4
	const churners = 4
	const perGoroutine = 2000
	bus := NewEventBus()
	var heard [3]atomic.Int64
	for i := range heard {
		h := &heard[i]
		bus.Subscribe(func(string) { h.Add(1) })
	}
	var badCount, badUnsubscribe atomic.Bool

	runTogether(t, publishers+churners, func(g int) {
		for i := 0; i < perGoroutine; i++ {
			if g < publishers {
				if bus.Publish("tick") < 3 {
					badCount.Store(true)
				}
			} else {
				id := bus.Subscribe(func(string) {})
				if !bus.Unsubscribe(id) {
					badUnsubscribe.Store(true)
				}
			}
		}
	})

	for i := range heard {
		if n := heard[i].Load(); n != publishers*perGoroutine {
			t.Errorf("permanent subscriber %d heard %d events, want %d", i, n, publishers*perGoroutine)
		}
	}
	if badCount.Load() {
		t.Error("a publish missed a permanent subscriber")
	}
	if badUnsubscribe.Load() {
		t.Error("Unsubscribe of a fresh id returned false")
	}
	if n := bus.SubscriberCount(); n != 3 {
		t.Errorf("SubscriberCount() = %d, want 3", n)
	}
}
