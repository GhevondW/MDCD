package day04

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func mustExecuteWS(t *testing.T, p *WorkStealingPool, task func(w *Worker)) {
	t.Helper()
	if err := p.Execute(task); err != nil {
		t.Fatalf("Execute returned %v, want nil", err)
	}
}

// fib forks fib(n-1), computes fib(n-2) itself, then joins. On one worker
// this only works if Join helps: the forked task sits in the worker's own
// deque, and the worker must run it instead of sleeping.
func fib(w *Worker, n int) int {
	if n < 2 {
		return n
	}
	h := Fork(w, func(w *Worker) int { return fib(w, n-1) })
	b := fib(w, n-2)
	a, err := h.Join(w)
	if err != nil {
		panic(err)
	}
	return a + b
}

func sum(w *Worker, xs []int) int {
	if len(xs) <= 1000 {
		s := 0
		for _, x := range xs {
			s += x
		}
		return s
	}
	mid := len(xs) / 2
	left := Fork(w, func(w *Worker) int { return sum(w, xs[:mid]) })
	right := sum(w, xs[mid:])
	l, err := left.Join(w)
	if err != nil {
		panic(err)
	}
	return l + right
}

func TestStealingPoolWorkersCount(t *testing.T) {
	watchdog(t)
	p := NewWorkStealingPool(3)
	defer p.Shutdown()
	if p.Workers() != 3 {
		t.Errorf("Workers() = %d, want 3", p.Workers())
	}
	func() {
		defer func() {
			if recover() == nil {
				t.Error("NewWorkStealingPool(0) did not panic")
			}
		}()
		NewWorkStealingPool(0)
	}()
}

func TestStealingPoolRunsExternalTasks(t *testing.T) {
	watchdog(t)
	p := NewWorkStealingPool(4)
	var n atomic.Int32
	for i := 0; i < 1000; i++ {
		mustExecuteWS(t, p, func(w *Worker) { n.Add(1) })
	}
	p.Shutdown()
	if got := n.Load(); got != 1000 {
		t.Errorf("%d of 1000 tasks ran", got)
	}
}

func TestStealingPoolFibOnOneWorkerDoesNotDeadlock(t *testing.T) {
	watchdog(t)
	// The lecture's "deadlock with no locks": a parent waits for its
	// child, and the only worker is the parent. Join must run the child.
	p := NewWorkStealingPool(1)
	defer p.Shutdown()
	h := Run(p, func(w *Worker) int { return fib(w, 20) })
	v, err := waitWithin(t, h, 20*time.Second)
	if err != nil || v != 6765 {
		t.Errorf("fib(20) = (%d, %v), want (6765, nil)", v, err)
	}
	if s := p.Steals(); s != 0 {
		t.Errorf("Steals() = %d on a pool with one worker: there is nobody to steal from", s)
	}
}

func TestStealingPoolFibOnManyWorkers(t *testing.T) {
	watchdog(t)
	for _, workers := range []int{2, 4} {
		p := NewWorkStealingPool(workers)
		h := Run(p, func(w *Worker) int { return fib(w, 24) })
		v, err := waitWithin(t, h, 30*time.Second)
		if err != nil || v != 46368 {
			t.Errorf("%d workers: fib(24) = (%d, %v), want (46368, nil)", workers, v, err)
		}
		p.Shutdown()
	}
}

func TestStealingPoolParallelSum(t *testing.T) {
	watchdog(t)
	xs := make([]int, 1_000_000)
	for i := range xs {
		xs[i] = i + 1
	}
	want := len(xs) * (len(xs) + 1) / 2
	for _, workers := range []int{1, 4} {
		p := NewWorkStealingPool(workers)
		h := Run(p, func(w *Worker) int { return sum(w, xs) })
		v, err := waitWithin(t, h, 30*time.Second)
		if err != nil || v != want {
			t.Errorf("%d workers: sum = (%d, %v), want (%d, nil)", workers, v, err, want)
		}
		p.Shutdown()
	}
}

func TestStealingPoolOwnerTakesNewestFirst(t *testing.T) {
	watchdog(t)
	// One worker. A task forks A, B, C and returns without joining. The
	// worker then takes from the BOTTOM of its own deque: newest first.
	p := NewWorkStealingPool(1)
	var mu sync.Mutex
	var order []string
	note := func(s string) {
		mu.Lock()
		order = append(order, s)
		mu.Unlock()
	}
	mustExecuteWS(t, p, func(w *Worker) {
		for _, name := range []string{"A", "B", "C"} {
			name := name
			w.Fork(func(w *Worker) { note(name) })
		}
	})
	p.Shutdown()
	mu.Lock()
	defer mu.Unlock()
	if len(order) != 3 || order[0] != "C" || order[1] != "B" || order[2] != "A" {
		t.Errorf("forked tasks ran in order %v, want [C B A] (a worker takes its own newest task first)", order)
	}
}

func TestStealingPoolThiefTakesOldestFirst(t *testing.T) {
	watchdog(t)
	// Two workers. One runs a task that forks A, B, C and then just waits
	// (without joining, so it does not help). The other worker is idle: it
	// has to STEAL, and it takes from the TOP of the first worker's
	// deque: oldest first.
	p := NewWorkStealingPool(2)
	defer p.Shutdown()
	var mu sync.Mutex
	var order []string
	var ran atomic.Int32
	allRan := newSignal()
	mustExecuteWS(t, p, func(w *Worker) {
		for _, name := range []string{"A", "B", "C"} {
			name := name
			w.Fork(func(w *Worker) {
				mu.Lock()
				order = append(order, name)
				mu.Unlock()
				if ran.Add(1) == 3 {
					allRan.fire()
				}
			})
		}
		allRan.waitFor(10 * time.Second) // a plain wait: this worker stays busy
	})
	if !allRan.waitFor(10 * time.Second) {
		t.Fatal("the forked tasks never ran on the idle worker: an idle worker must steal")
	}
	mu.Lock()
	defer mu.Unlock()
	if len(order) != 3 || order[0] != "A" || order[1] != "B" || order[2] != "C" {
		t.Errorf("stolen tasks ran in order %v, want [A B C] (a thief takes the oldest task first)", order)
	}
	if s := p.Steals(); s != 3 {
		t.Errorf("Steals() = %d, want 3", s)
	}
}

func TestStealingPoolSkewedWorkIsSpreadByStealing(t *testing.T) {
	watchdog(t)
	// ONE task forks 200 children that each sleep 2 ms. They all start on
	// that one worker's deque; the other three workers must steal them.
	const children = 200
	const nap = 2 * time.Millisecond
	serial := sleepCost(nap, 20) * children

	p := NewWorkStealingPool(4)
	defer p.Shutdown()
	start := time.Now()
	h := Run(p, func(w *Worker) int {
		hs := make([]*Handle[int], children)
		for i := range hs {
			hs[i] = Fork(w, func(w *Worker) int { time.Sleep(nap); return 1 })
		}
		total := 0
		for _, c := range hs {
			v, err := c.Join(w)
			if err != nil {
				panic(err)
			}
			total += v
		}
		return total
	})
	v, err := waitWithin(t, h, 30*time.Second)
	took := time.Since(start)
	if err != nil || v != children {
		t.Fatalf("result = (%d, %v), want (%d, nil)", v, err, children)
	}
	if p.Steals() == 0 {
		t.Error("Steals() = 0: the other workers never stole from the busy one")
	}
	if limit := serial * 6 / 10; took > limit {
		t.Errorf("took %v; one worker alone would need about %v, and 4 workers should be well under %v", took, serial, limit)
	}
}

func TestStealingPoolNeverRunsMoreTasksThanWorkers(t *testing.T) {
	watchdog(t)
	p := NewWorkStealingPool(3)
	var cur, peak atomic.Int32
	for i := 0; i < 200; i++ {
		mustExecuteWS(t, p, func(w *Worker) {
			c := cur.Add(1)
			for {
				m := peak.Load()
				if c <= m || peak.CompareAndSwap(m, c) {
					break
				}
			}
			time.Sleep(time.Millisecond)
			cur.Add(-1)
		})
	}
	p.Shutdown()
	if m := peak.Load(); m > 3 {
		t.Errorf("%d tasks ran at the same time on a pool of 3 workers", m)
	}
}

func TestStealingPoolShutdownWaitsForForkedWork(t *testing.T) {
	watchdog(t)
	// One task forks 100 children, each of which forks 10 more -- and
	// nobody waits for any of it. Shutdown is called right away. It must
	// not return (or stop the workers) until all 1 + 100 + 1000 tasks ran.
	p := NewWorkStealingPool(4)
	var n atomic.Int32
	mustExecuteWS(t, p, func(w *Worker) {
		n.Add(1)
		for i := 0; i < 100; i++ {
			w.Fork(func(w *Worker) {
				n.Add(1)
				for j := 0; j < 10; j++ {
					w.Fork(func(w *Worker) {
						time.Sleep(50 * time.Microsecond)
						n.Add(1)
					})
				}
			})
		}
	})
	p.Shutdown()
	if got := n.Load(); got != 1101 {
		t.Errorf("%d of 1101 tasks had run when Shutdown returned: it must wait for forked tasks too", got)
	}
}

func TestStealingPoolRefusesExternalTasksAfterShutdown(t *testing.T) {
	watchdog(t)
	p := NewWorkStealingPool(2)
	p.Shutdown()
	p.Shutdown() // twice is fine
	if err := p.Execute(func(w *Worker) {}); !errors.Is(err, ErrPoolClosed) {
		t.Errorf("Execute after Shutdown returned %v, want ErrPoolClosed", err)
	}
	h := Run(p, func(w *Worker) int { return 1 })
	if _, err := waitWithin(t, h, 5*time.Second); !errors.Is(err, ErrPoolClosed) {
		t.Errorf("Run on a shut-down pool: Wait() error = %v, want ErrPoolClosed", err)
	}
}

func TestStealingPoolShutdownFromManyGoroutines(t *testing.T) {
	watchdog(t)
	p := NewWorkStealingPool(3)
	var n atomic.Int32
	for i := 0; i < 100; i++ {
		mustExecuteWS(t, p, func(w *Worker) { time.Sleep(100 * time.Microsecond); n.Add(1) })
	}
	runTogether(t, 6, func(g int) {
		p.Shutdown()
		if got := n.Load(); got != 100 {
			t.Errorf("a Shutdown call returned when only %d of 100 tasks had run", got)
		}
	})
}

func TestStealingPoolPanicDoesNotKillAWorker(t *testing.T) {
	watchdog(t)
	p := NewWorkStealingPool(1)
	defer p.Shutdown()
	mustExecuteWS(t, p, func(w *Worker) { panic("raw boom") })

	h := Run(p, func(w *Worker) int { panic("handle boom") })
	_, err := waitWithin(t, h, 5*time.Second)
	var pe *PanicError
	if !errors.As(err, &pe) || pe.Value != "handle boom" {
		t.Errorf("Wait() error = %v, want a *PanicError carrying \"handle boom\"", err)
	}

	// A panicking child joined by its parent reaches the parent as an error.
	g := Run(p, func(w *Worker) string {
		c := Fork(w, func(w *Worker) int { panic("child boom") })
		if _, err := c.Join(w); err != nil {
			return "joined: " + err.Error()
		}
		return "no error?"
	})
	if v, _ := waitWithin(t, g, 5*time.Second); v != "joined: task panicked: child boom" {
		t.Errorf("parent got %q", v)
	}
}

func TestStealingPoolManySubmittersAtOnce(t *testing.T) {
	watchdog(t)
	p := NewWorkStealingPool(4)
	var n atomic.Int32
	runTogether(t, 8, func(g int) {
		for i := 0; i < 500; i++ {
			// each external task also forks and joins one child
			err := p.Execute(func(w *Worker) {
				h := Fork(w, func(w *Worker) int { n.Add(1); return 1 })
				h.Join(w)
			})
			if err != nil {
				t.Errorf("Execute returned %v", err)
				return
			}
		}
	})
	p.Shutdown()
	if got := n.Load(); got != 4000 {
		t.Errorf("%d of 4000 forked tasks ran", got)
	}
}
