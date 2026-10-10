package day04

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func mustExecute(t *testing.T, p *ThreadPool, task func()) {
	t.Helper()
	if err := p.Execute(task); err != nil {
		t.Fatalf("Execute returned %v, want nil", err)
	}
}

// blockOneWorker makes a pool with one worker busy: it submits a task that
// waits for release.fire(), and returns once that task is really running
// (so it is no longer counted as waiting in the queue).
func blockOneWorker(t *testing.T, p *ThreadPool) (release *signal) {
	t.Helper()
	started, release := newSignal(), newSignal()
	mustExecute(t, p, func() {
		started.fire()
		release.waitFor(30 * time.Second)
	})
	if !started.waitFor(5 * time.Second) {
		t.Fatal("the first task never started")
	}
	return release
}

func TestPoolRejectsBadArguments(t *testing.T) {
	watchdog(t)
	for _, c := range [][2]int{{0, 0}, {-1, 0}, {1, -1}} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("NewThreadPool(%d, %d) did not panic", c[0], c[1])
				}
			}()
			NewThreadPool(c[0], c[1])
		}()
	}
}

func TestPoolRunsEveryTask(t *testing.T) {
	watchdog(t)
	p := NewThreadPool(4, 0)
	var n atomic.Int32
	for i := 0; i < 1000; i++ {
		mustExecute(t, p, func() { n.Add(1) })
	}
	p.Shutdown()
	if got := n.Load(); got != 1000 {
		t.Errorf("%d of 1000 tasks ran; Shutdown must run every accepted task", got)
	}
}

func TestPoolRunsTasksInParallel(t *testing.T) {
	watchdog(t)
	// Four tasks that each wait until all four are running at once. On a
	// pool that runs fewer than four at a time they never meet.
	p := NewThreadPool(4, 0)
	defer p.Shutdown()
	var arrived atomic.Int32
	all := newSignal()
	for i := 0; i < 4; i++ {
		mustExecute(t, p, func() {
			if arrived.Add(1) == 4 {
				all.fire()
			}
			all.waitFor(10 * time.Second)
		})
	}
	if !all.waitFor(5 * time.Second) {
		t.Fatalf("only %d of 4 tasks were running at the same time: a pool of 4 workers must run 4 tasks in parallel", arrived.Load())
	}
}

func TestPoolNeverRunsMoreTasksThanWorkers(t *testing.T) {
	watchdog(t)
	p := NewThreadPool(3, 0)
	var cur, max atomic.Int32
	for i := 0; i < 200; i++ {
		mustExecute(t, p, func() {
			c := cur.Add(1)
			for {
				m := max.Load()
				if c <= m || max.CompareAndSwap(m, c) {
					break
				}
			}
			time.Sleep(time.Millisecond)
			cur.Add(-1)
		})
	}
	p.Shutdown()
	if m := max.Load(); m > 3 {
		t.Errorf("%d tasks ran at the same time on a pool of 3 workers", m)
	}
}

func TestPoolWithOneWorkerRunsTasksInOrder(t *testing.T) {
	watchdog(t)
	p := NewThreadPool(1, 0)
	var order []int // only the single worker appends
	for i := 0; i < 100; i++ {
		i := i
		mustExecute(t, p, func() { order = append(order, i) })
	}
	p.Shutdown()
	if len(order) != 100 {
		t.Fatalf("%d of 100 tasks ran", len(order))
	}
	for i, v := range order {
		if v != i {
			t.Fatalf("task %d ran at position %d: tasks must start in the order they were accepted", v, i)
		}
	}
}

func TestPoolSubmitReturnsTheResult(t *testing.T) {
	watchdog(t)
	p := NewThreadPool(2, 0)
	defer p.Shutdown()
	var futs []*Future[int]
	for i := 0; i < 20; i++ {
		i := i
		futs = append(futs, Submit(p, func() int { return i * i }))
	}
	for i, f := range futs {
		if v, err := getWithin(t, f, 5*time.Second); err != nil || v != i*i {
			t.Errorf("future %d = (%d, %v), want (%d, nil)", i, v, err, i*i)
		}
	}
}

func TestPoolPanicInASubmittedTaskBecomesAnError(t *testing.T) {
	watchdog(t)
	p := NewThreadPool(1, 0)
	defer p.Shutdown()
	f := Submit(p, func() int { panic("boom") })
	_, err := getWithin(t, f, 5*time.Second)
	var pe *PanicError
	if !errors.As(err, &pe) || pe.Value != "boom" {
		t.Fatalf("Get() error = %v, want a *PanicError carrying \"boom\"", err)
	}
	g := Submit(p, func() int { return 1 })
	if v, err := getWithin(t, g, 5*time.Second); err != nil || v != 1 {
		t.Errorf("after a panic the pool's next task gave (%d, %v), want (1, nil)", v, err)
	}
}

func TestPoolPanicInARawTaskDoesNotKillTheWorker(t *testing.T) {
	watchdog(t)
	// The one and only worker runs a task that panics. If the panic is not
	// recovered the whole test program dies; if the worker is lost, the
	// next task never runs.
	p := NewThreadPool(1, 0)
	defer p.Shutdown()
	mustExecute(t, p, func() { panic("raw boom") })
	f := Submit(p, func() string { return "still here" })
	if v, err := getWithin(t, f, 5*time.Second); err != nil || v != "still here" {
		t.Errorf("task after a panic: (%q, %v), want (\"still here\", nil)", v, err)
	}
}

func TestPoolTrySubmitRejectsWhenTheQueueIsFull(t *testing.T) {
	watchdog(t)
	p := NewThreadPool(1, 1)
	release := blockOneWorker(t, p)
	var ran atomic.Int32
	if err := p.TrySubmit(func() { ran.Add(1) }); err != nil {
		t.Fatalf("TrySubmit with room in the queue returned %v", err)
	}
	if !within(callTimeLimit, func() {
		if err := p.TrySubmit(func() { ran.Add(100) }); !errors.Is(err, ErrQueueFull) {
			t.Errorf("TrySubmit on a full queue returned %v, want ErrQueueFull", err)
		}
	}) {
		t.Fatal("TrySubmit waited on a full queue")
	}
	release.fire()
	p.Shutdown()
	if got := ran.Load(); got != 1 {
		t.Errorf("tasks that ran = %d (1 = only the accepted one), the rejected task must not run", got)
	}
}

func TestPoolExecuteWaitsWhenTheQueueIsFull(t *testing.T) {
	watchdog(t)
	p := NewThreadPool(1, 1)
	release := blockOneWorker(t, p)
	var ran atomic.Int32
	mustExecute(t, p, func() { ran.Add(1) }) // fills the queue
	returned := newSignal()
	var err error
	go func() {
		err = p.Execute(func() { ran.Add(1) })
		returned.fire()
	}()
	if returned.waitFor(150 * time.Millisecond) {
		t.Fatal("Execute returned although the queue was full: it must wait for room")
	}
	release.fire()
	if !returned.waitFor(5 * time.Second) {
		t.Fatal("Execute never returned after the queue drained")
	}
	if err != nil {
		t.Fatalf("Execute returned %v, want nil", err)
	}
	p.Shutdown()
	if got := ran.Load(); got != 2 {
		t.Errorf("%d queued tasks ran, want 2", got)
	}
}

func TestPoolShutdownRunsQueuedTasksAndRefusesNewOnes(t *testing.T) {
	watchdog(t)
	p := NewThreadPool(1, 0)
	release := blockOneWorker(t, p)
	var ran atomic.Int32
	for i := 0; i < 50; i++ {
		mustExecute(t, p, func() { ran.Add(1) })
	}
	shutdownDone := newSignal()
	go func() { p.Shutdown(); shutdownDone.fire() }()

	// Once Shutdown has begun, new tasks are refused.
	deadline := time.Now().Add(5 * time.Second)
	for {
		err := p.TrySubmit(func() {})
		if errors.Is(err, ErrPoolClosed) {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("TrySubmit still returns %v long after Shutdown was called", err)
		}
		time.Sleep(time.Millisecond)
	}
	if err := p.Execute(func() {}); !errors.Is(err, ErrPoolClosed) {
		t.Errorf("Execute after Shutdown returned %v, want ErrPoolClosed", err)
	}
	if shutdownDone.waitFor(150 * time.Millisecond) {
		t.Fatal("Shutdown returned while a task was still running and 50 were still queued")
	}
	release.fire()
	if !shutdownDone.waitFor(5 * time.Second) {
		t.Fatal("Shutdown did not return after the tasks finished")
	}
	if got := ran.Load(); got < 50 {
		t.Errorf("only %d of the 50 queued tasks ran: Shutdown must run every accepted task", got)
	}
}

func TestPoolShutdownWakesAnExecuteThatWaitsForRoom(t *testing.T) {
	watchdog(t)
	p := NewThreadPool(1, 1)
	release := blockOneWorker(t, p)
	mustExecute(t, p, func() {}) // fills the queue
	var ranLate atomic.Bool
	var err error
	returned := newSignal()
	go func() {
		err = p.Execute(func() { ranLate.Store(true) })
		returned.fire()
	}()
	time.Sleep(100 * time.Millisecond)
	shutdownDone := newSignal()
	go func() { p.Shutdown(); shutdownDone.fire() }()
	// The worker is still blocked, so the queue stays full: the only way
	// out for the waiting Execute is the shutdown.
	if !returned.waitFor(5 * time.Second) {
		t.Fatal("an Execute waiting for room kept waiting after Shutdown")
	}
	if !errors.Is(err, ErrPoolClosed) {
		t.Errorf("that Execute returned %v, want ErrPoolClosed", err)
	}
	release.fire()
	if !shutdownDone.waitFor(5 * time.Second) {
		t.Fatal("Shutdown did not return")
	}
	if ranLate.Load() {
		t.Error("a task whose Execute returned ErrPoolClosed was run anyway")
	}
}

func TestPoolShutdownWaitsForRunningTasks(t *testing.T) {
	watchdog(t)
	p := NewThreadPool(2, 0)
	var finished atomic.Bool
	started := newSignal()
	mustExecute(t, p, func() {
		started.fire()
		time.Sleep(200 * time.Millisecond)
		finished.Store(true)
	})
	started.wait()
	p.Shutdown()
	if !finished.Load() {
		t.Error("Shutdown returned while a task was still running")
	}
}

func TestPoolShutdownTwiceAndFromManyGoroutines(t *testing.T) {
	watchdog(t)
	p := NewThreadPool(3, 0)
	var n atomic.Int32
	for i := 0; i < 100; i++ {
		mustExecute(t, p, func() { time.Sleep(100 * time.Microsecond); n.Add(1) })
	}
	runTogether(t, 6, func(g int) {
		p.Shutdown()
		if got := n.Load(); got != 100 {
			t.Errorf("a Shutdown call returned when only %d of 100 tasks had run", got)
		}
	})
	p.Shutdown()
	if err := p.Execute(func() {}); !errors.Is(err, ErrPoolClosed) {
		t.Errorf("Execute after Shutdown returned %v, want ErrPoolClosed", err)
	}
	f := Submit(p, func() int { return 1 })
	if _, err := getWithin(t, f, 5*time.Second); !errors.Is(err, ErrPoolClosed) {
		t.Errorf("Submit on a shut-down pool: Get() error = %v, want ErrPoolClosed", err)
	}
}

func TestPoolTasksCanSubmitMoreTasks(t *testing.T) {
	watchdog(t)
	// Each task submits two children and does not wait for them; 2^11 - 1
	// tasks in all. (Waiting for a child inside a task is the deadlock
	// from the lecture -- problem 3 fixes that.)
	p := NewThreadPool(4, 0)
	var n atomic.Int32
	var spawn func(depth int)
	spawn = func(depth int) {
		n.Add(1)
		if depth == 0 {
			return
		}
		for c := 0; c < 2; c++ {
			if err := p.Execute(func() { spawn(depth - 1) }); err != nil {
				t.Errorf("Execute from inside a task returned %v", err)
			}
		}
	}
	mustExecute(t, p, func() { spawn(10) })
	deadline := time.Now().Add(10 * time.Second)
	for n.Load() != 2047 && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	p.Shutdown()
	if got := n.Load(); got != 2047 {
		t.Errorf("%d of 2047 tasks ran", got)
	}
}

func TestPoolManySubmittersAtOnce(t *testing.T) {
	watchdog(t)
	p := NewThreadPool(4, 8)
	var n atomic.Int32
	runTogether(t, 8, func(g int) {
		for i := 0; i < 1000; i++ {
			if err := p.Execute(func() { n.Add(1) }); err != nil {
				t.Errorf("Execute returned %v", err)
				return
			}
		}
	})
	p.Shutdown()
	if got := n.Load(); got != 8000 {
		t.Errorf("%d of 8000 tasks ran", got)
	}
}

func TestPoolEveryAcceptedTaskRunsEvenWhenShutdownRaces(t *testing.T) {
	watchdog(t)
	// Submitters and a Shutdown run at the same moment. Whatever Execute
	// accepted (returned nil for) must have run once Shutdown returns --
	// nothing accepted may be lost, and nothing may run twice.
	for round := 0; round < 30; round++ {
		p := NewThreadPool(3, 4)
		var accepted, ran atomic.Int32
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(time.Duration(round%5) * 200 * time.Microsecond)
			p.Shutdown()
		}()
		runTogether(t, 6, func(g int) {
			for i := 0; i < 200; i++ {
				if p.Execute(func() { ran.Add(1) }) == nil {
					accepted.Add(1)
				}
			}
		})
		wg.Wait()
		p.Shutdown()
		if a, r := accepted.Load(), ran.Load(); a != r {
			t.Fatalf("round %d: %d tasks were accepted but %d ran", round, a, r)
		}
	}
}
