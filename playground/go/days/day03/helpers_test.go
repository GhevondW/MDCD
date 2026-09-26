package day03

// Helpers shared by this day's tests.
//
// A deadlocked test would otherwise hang `go test` until its 10-minute
// timeout, so these tests guard against deadlocks in two ways:
//
//   - watchdog(t), at the top of every test, stops the whole run after
//     60 s with a panic that names the test and prints the stack of every
//     goroutine -- look for the ones stuck in sync.(*Mutex).Lock. (A Go
//     mutex is not reentrant: a method that holds the lock must not call
//     another method that takes the same lock.)
//   - runTogether, within, and the xxxWithin helpers fail just the one
//     test, sooner, and leave the stuck goroutines behind.

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const (
	testTimeLimit      = 60 * time.Second // a whole test (watchdog)
	goroutineTimeLimit = 30 * time.Second // all goroutines of one runTogether
	callTimeLimit      = 5 * time.Second  // one call that could deadlock
)

func must(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// watchdog stops the whole test run if the test is still running after
// testTimeLimit, instead of letting a deadlock hang it.
func watchdog(t *testing.T) {
	name := t.Name()
	timer := time.AfterFunc(testTimeLimit, func() {
		debug.SetTraceback("all")
		panic(fmt.Sprintf("%s is still running after %v -- deadlocked? "+
			"The stacks below show where every goroutine is stuck.", name, testTimeLimit))
	})
	t.Cleanup(func() { timer.Stop() })
}

// runTogether runs body(0) ... body(n-1) on n goroutines that all start
// at (nearly) the same moment, so their calls overlap as much as
// possible, and waits for all of them. If they are not all done after
// goroutineTimeLimit -- most likely a deadlock -- the test fails instead
// of hanging.
//
// body runs on other goroutines, so it must not call t.Fatal, t.FailNow
// or must: record what went wrong and check it after runTogether
// returns.
func runTogether(t *testing.T, n int, body func(g int)) {
	t.Helper()
	var ready atomic.Int32
	var start atomic.Bool
	var wg sync.WaitGroup
	wg.Add(n)
	for g := 0; g < n; g++ {
		go func(g int) {
			defer wg.Done()
			ready.Add(1)
			for !start.Load() {
				runtime.Gosched()
			}
			body(g)
		}(g)
	}
	for ready.Load() < int32(n) {
		runtime.Gosched()
	}
	start.Store(true)
	if !within(goroutineTimeLimit, wg.Wait) {
		t.Fatalf("the goroutines were still running after %v -- deadlocked?", goroutineTimeLimit)
	}
}

// within runs f on another goroutine and waits up to d for it to return.
// It reports whether f returned in time; if not, f is left running.
func within(d time.Duration, f func()) bool {
	done := make(chan struct{})
	go func() {
		defer close(done)
		f()
	}()
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-done:
		return true
	case <-timer.C:
		return false
	}
}

// signal is a one-shot "it happened" flag that goroutines can wait for.
// fire() may be called more than once; only the first call counts.
type signal struct {
	once sync.Once
	ch   chan struct{}
}

func newSignal() *signal { return &signal{ch: make(chan struct{})} }

func (s *signal) fire() { s.once.Do(func() { close(s.ch) }) }

// done returns a channel that is closed once fire() has been called.
func (s *signal) done() <-chan struct{} { return s.ch }

// isSet reports, without waiting, whether fire() has been called.
func (s *signal) isSet() bool {
	select {
	case <-s.ch:
		return true
	default:
		return false
	}
}

// wait waits for fire(), however long it takes.
func (s *signal) wait() { <-s.ch }

// waitFor waits up to d for fire(); it reports whether fire() happened.
func (s *signal) waitFor(d time.Duration) bool {
	if s.isSet() {
		return true
	}
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-s.ch:
		return true
	case <-timer.C:
		return false
	}
}

// firstFailure remembers the first thing that went wrong on another
// goroutine, so the test goroutine can report it after runTogether.
type firstFailure struct {
	mu  sync.Mutex
	msg string
}

func (f *firstFailure) record(format string, args ...any) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.msg == "" {
		f.msg = fmt.Sprintf(format, args...)
	}
}

func (f *firstFailure) report(t *testing.T) {
	t.Helper()
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.msg != "" {
		t.Error(f.msg)
	}
}
