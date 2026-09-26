package day03

import (
	"errors"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// getWithin calls c.Get(key, compute) and fails the test if it has not
// returned after callTimeLimit -- for example because it waits for a
// computation that already finished, or already failed.
func getWithin[V any](t *testing.T, c *OnceCache[V], key string, compute func() (V, error)) (V, error) {
	t.Helper()
	var v V
	var err error
	if !within(callTimeLimit, func() { v, err = c.Get(key, compute) }) {
		t.Fatalf("Get(%q) did not return within %v -- is it waiting for a computation "+
			"that already finished, or already failed?", key, callTimeLimit)
	}
	return v, err
}

func TestOnceCacheMissComputesAndReturnsTheValue(t *testing.T) {
	watchdog(t)
	c := NewOnceCache[int]()
	calls := 0
	v, err := getWithin(t, c, "a", func() (int, error) { calls++; return 42, nil })
	if err != nil || v != 42 {
		t.Errorf("Get(\"a\") = (%d, %v), want (42, nil)", v, err)
	}
	if calls != 1 {
		t.Errorf("compute() ran %d times, want 1", calls)
	}
	if got := c.Size(); got != 1 {
		t.Errorf("Size() = %d, want 1", got)
	}
}

func TestOnceCacheHitReturnsRememberedValueWithoutComputing(t *testing.T) {
	watchdog(t)
	c := NewOnceCache[int]()
	getWithin(t, c, "a", func() (int, error) { return 42, nil })
	called := false
	v, err := getWithin(t, c, "a", func() (int, error) { called = true; return 7, nil })
	if err != nil || v != 42 {
		t.Errorf("second Get(\"a\") = (%d, %v), want the remembered (42, nil)", v, err)
	}
	if called {
		t.Error("compute() ran for a key that already had a value")
	}
}

func TestOnceCacheDifferentKeysGetTheirOwnValues(t *testing.T) {
	watchdog(t)
	c := NewOnceCache[string]()
	for _, tc := range []struct{ key, computed, want string }{
		{"a", "apple", "apple"},
		{"b", "banana", "banana"},
		{"a", "avocado", "apple"},
	} {
		computed := tc.computed
		v, err := getWithin(t, c, tc.key, func() (string, error) { return computed, nil })
		if err != nil || v != tc.want {
			t.Errorf("Get(%q) = (%q, %v), want (%q, nil)", tc.key, v, err, tc.want)
		}
	}
	if got := c.Size(); got != 2 {
		t.Errorf("Size() = %d, want 2", got)
	}
}

func TestOnceCacheFailingComputeRemembersNothing(t *testing.T) {
	watchdog(t)
	c := NewOnceCache[int]()
	boom := errors.New("boom")
	if _, err := getWithin(t, c, "a", func() (int, error) { return 0, boom }); !errors.Is(err, boom) {
		t.Errorf("Get(\"a\") with a failing compute(): got error %v, want the error compute() returned", err)
	}
	if got := c.Size(); got != 0 {
		t.Errorf("Size() = %d after a failed compute(), want 0", got)
	}
	calls := 0
	v, err := getWithin(t, c, "a", func() (int, error) { calls++; return 5, nil })
	if err != nil || v != 5 {
		t.Errorf("Get(\"a\") after the failure = (%d, %v), want (5, nil)", v, err)
	}
	if calls != 1 {
		t.Errorf("compute() ran %d times -- after a failure, the next Get() must compute again", calls)
	}
	if got := c.Size(); got != 1 {
		t.Errorf("Size() = %d, want 1", got)
	}
}

func TestOnceCacheConcurrentGetsOfOneKeyComputeOnce(t *testing.T) {
	watchdog(t)
	// 16 goroutines miss the same key at the same moment; compute() is
	// slow, so a check-then-act gap is wide open.
	const goroutines = 16
	c := NewOnceCache[int]()
	var calls atomic.Int32
	results := make([]int, goroutines)
	errs := make([]error, goroutines)
	runTogether(t, goroutines, func(g int) {
		results[g], errs[g] = c.Get("config", func() (int, error) {
			calls.Add(1)
			time.Sleep(50 * time.Millisecond)
			return 99, nil
		})
	})

	if n := calls.Load(); n != 1 {
		t.Errorf("compute() ran %d times for one key, want exactly 1", n)
	}
	for g := range results {
		if results[g] != 99 || errs[g] != nil {
			t.Errorf("goroutine %d: Get(\"config\") = (%d, %v), want (99, nil)", g, results[g], errs[g])
			break
		}
	}
	if got := c.Size(); got != 1 {
		t.Errorf("Size() = %d, want 1", got)
	}
}

func TestOnceCacheConcurrentGetsOfManyKeysComputeEachOnce(t *testing.T) {
	watchdog(t)
	const goroutines = 8
	const keys = 200
	c := NewOnceCache[int]()
	calls := make([]atomic.Int32, keys)
	var bad firstFailure
	runTogether(t, goroutines, func(g int) {
		for i := 0; i < keys; i++ {
			k := (i + g*25) % keys // each goroutine starts at a different key
			v, err := c.Get(strconv.Itoa(k), func() (int, error) {
				calls[k].Add(1)
				time.Sleep(time.Millisecond)
				return k * 10, nil
			})
			if err != nil || v != k*10 {
				bad.record("Get(%q) = (%d, %v), want (%d, nil)", strconv.Itoa(k), v, err, k*10)
				return
			}
		}
	})

	bad.report(t)
	for k := range calls {
		if n := calls[k].Load(); n != 1 {
			t.Errorf("key %d was computed %d times, want exactly 1", k, n)
			break
		}
	}
	if got := c.Size(); got != keys {
		t.Errorf("Size() = %d, want %d", got, keys)
	}
}

func TestOnceCacheDifferentKeysDoNotWaitForEachOther(t *testing.T) {
	watchdog(t)
	// compute("a") and compute("b") each wait until the other one has
	// started. That only works if both can run at the same time; if a
	// lock is held while compute() runs, one of them gives up after 2 s.
	c := NewOnceCache[string]()
	aStarted, bStarted := newSignal(), newSignal()
	waitForOther := func(started, other *signal) func() (string, error) {
		return func() (string, error) {
			started.fire()
			if other.waitFor(2 * time.Second) {
				return "ok", nil
			}
			return "timed out", nil
		}
	}
	var a, b string
	var aErr, bErr error
	bothReturned := within(10*time.Second, func() {
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			a, aErr = c.Get("a", waitForOther(aStarted, bStarted))
		}()
		go func() {
			defer wg.Done()
			b, bErr = c.Get("b", waitForOther(bStarted, aStarted))
		}()
		wg.Wait()
	})
	if !bothReturned {
		t.Fatal("Get(\"a\") and Get(\"b\") had not both returned after 10 s -- deadlocked?")
	}

	if a != "ok" || aErr != nil {
		t.Errorf("Get(\"a\") = (%q, %v), want (\"ok\", nil) -- compute(\"a\") waited for compute(\"b\")", a, aErr)
	}
	if b != "ok" || bErr != nil {
		t.Errorf("Get(\"b\") = (%q, %v), want (\"ok\", nil) -- compute(\"b\") waited for compute(\"a\")", b, bErr)
	}
}

func TestOnceCacheSizeCountsOnlyFinishedValues(t *testing.T) {
	watchdog(t)
	c := NewOnceCache[int]()
	started, release, getReturned := newSignal(), newSignal(), newSignal()
	defer release.fire() // never leave compute() waiting, whatever happens below
	go func() {
		defer getReturned.fire()
		c.Get("slow", func() (int, error) {
			started.fire()
			release.wait()
			return 1, nil
		})
	}()

	// Wait until compute() is running -- or Get returned without calling it.
	select {
	case <-started.done():
	case <-getReturned.done():
	case <-time.After(callTimeLimit):
	}
	if !started.isSet() {
		t.Fatal("compute() was never called")
	}

	// Ask Size() on another goroutine, so a Size() that waits for the
	// running compute() fails this test instead of hanging it.
	sizeCh := make(chan int, 1)
	go func() { sizeCh <- c.Size() }()
	sizeAnswered := false
	sizeWhileRunning := 0
	select {
	case sizeWhileRunning = <-sizeCh:
		sizeAnswered = true
	case <-time.After(2 * time.Second):
	}
	release.fire()
	if !getReturned.waitFor(callTimeLimit) {
		t.Fatalf("Get(\"slow\") did not return within %v after compute() returned", callTimeLimit)
	}

	if !sizeAnswered {
		t.Error("Size() waited for a running compute()")
	} else if sizeWhileRunning != 0 {
		t.Errorf("Size() = %d while compute() was still running, want 0 "+
			"-- a computation still running counted as remembered", sizeWhileRunning)
	}
	if got := c.Size(); got != 1 {
		t.Errorf("Size() = %d after compute() finished, want 1", got)
	}
}
