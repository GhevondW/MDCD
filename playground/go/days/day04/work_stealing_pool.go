package day04

// Day 4 (challenge) -- A work-stealing pool.
//
// One shared queue is one lock: every submit and every take fights for
// it. And a pool whose workers WAIT for sub-tasks can deadlock with no
// lock at all (the slides' "parent waits for child" cycle). A
// work-stealing pool fixes both -- this is the idea behind Java's
// ForkJoinPool, Go's own scheduler, Rust's rayon and the .NET pool.
//
//   - Every worker owns a DEQUE of tasks.
//   - A task running on a worker may Fork new tasks; they go onto the
//     BOTTOM of that worker's own deque.
//   - A worker takes its next task from the bottom of its own deque
//     (newest first, LIFO): the most recently forked task is the one
//     whose data is still warm in the cache, and the one a Join is most
//     likely waiting for.
//   - A worker whose deque is empty first takes from the shared
//     INJECTION queue (tasks submitted from outside, FIFO), and if that
//     is empty too, STEALS the oldest task (the TOP, FIFO) of another
//     worker's deque. Oldest first: old tasks tend to be big ones that
//     split into more work, so a thief gets a lot for one steal.
//   - Join does not sleep: while the task it waits for is not done, the
//     worker keeps running other tasks (its own, injected, stolen) --
//     "helping". A pool with ONE worker can therefore compute a
//     recursive fib without deadlock.
//
// Go has no goroutine identity, so a task can't ask "which worker am I
// on?". Instead the pool hands every task its worker: tasks are
// func(w *Worker), and everything that must happen "on the current
// worker" goes through w: w.Fork(...), handle.Join(w). A *Worker is only
// valid inside the task it was passed to, on that task's goroutine.
//
// You write: NewWorkStealingPool, Execute, Shutdown, Steals, Workers,
// Worker.Fork, Worker.helpUntil, and the worker loop and deques behind
// them. Provided below (don't change): Handle[T], Run and Fork, which
// wrap Execute / Worker.Fork and call helpUntil to Join.
//
// The contract:
//
//   - NewWorkStealingPool(workers): starts `workers` goroutines. Panics
//     if workers < 1. Workers() returns that number.
//   - Execute(task): submit from OUTSIDE the pool (any goroutine that is
//     not a worker). The task goes into the injection queue. Returns nil
//     when accepted -- an accepted task always runs exactly once -- or
//     ErrPoolClosed (not run) once Shutdown has been called.
//   - w.Fork(task): called by a running task; pushes onto the bottom of
//     w's own deque. Never fails and never waits -- not even during
//     Shutdown: a forked task is part of work already accepted.
//   - w.helpUntil(done): returns only once done() is true. Until then
//     it runs tasks (own deque newest-first, then injection queue, then
//     stealing). If there is nothing to run right now, yield
//     (runtime.Gosched) or sleep a few microseconds and look again --
//     the task being waited for is running on some other worker. Don't
//     block on a lock that only a *finished task* would release.
//   - A worker with nothing to do and nothing to wait for parks on a
//     sync.Cond instead of spinning. Whoever adds work (Execute, Fork)
//     must wake a sleeper.
//   - Shutdown(): stop accepting EXTERNAL tasks (Execute -> ErrPoolClosed
//     from now on), wait until every accepted task AND everything those
//     tasks forked has finished, then stop the workers and return after
//     they exited. Hint: count tasks that were accepted or forked but
//     not yet finished. Idempotent, safe from many goroutines; don't call
//     it from inside a task.
//   - A task that panics must not kill its worker (recover it; it still
//     counts as finished).
//   - Steals() returns how many tasks were taken from ANOTHER worker's
//     deque so far (not from the injection queue, not from your own
//     deque).
//
// Deques: a mutex per deque is fine (the lock-free Chase-Lev deque is
// the real thing, but that is a research project). The point is that
// the owner and the thieves lock the deque of ONE worker, not a global
// lock, and that the owner takes from one end and the thieves from the
// other.

import (
	"errors"
	"sync/atomic"
)

type WorkStealingPool struct {
	// TODO: choose your own representation.
}

// Worker is one goroutine of the pool, with its own deque. Tasks get
// theirs as an argument.
type Worker struct {
	// TODO: choose your own representation.
}

func NewWorkStealingPool(workers int) *WorkStealingPool {
	// TODO
	return &WorkStealingPool{}
}

func (p *WorkStealingPool) Execute(task func(w *Worker)) error {
	// TODO
	return errNotImplemented
}

func (p *WorkStealingPool) Shutdown() {
	// TODO
}

func (p *WorkStealingPool) Steals() int {
	// TODO
	return 0
}

func (p *WorkStealingPool) Workers() int {
	// TODO
	return 0
}

func (w *Worker) Fork(task func(w *Worker)) {
	// TODO
}

func (w *Worker) helpUntil(done func() bool) {
	// TODO
}

// ---- provided: handles to results ---------------------------------------

// Handle is the read end for the result of a forked or submitted task.
type Handle[T any] struct {
	fin  atomic.Bool
	done chan struct{}
	val  T
	err  error
}

func newHandle[T any]() *Handle[T] { return &Handle[T]{done: make(chan struct{})} }

func (h *Handle[T]) finished() bool { return h.fin.Load() }

func (h *Handle[T]) complete(v T, err error) {
	h.val, h.err = v, err
	h.fin.Store(true)
	close(h.done)
}

// run wraps f so that its result (or panic, as a *PanicError) lands in h.
func (h *Handle[T]) run(f func(w *Worker) T) func(w *Worker) {
	return func(w *Worker) {
		defer func() {
			if r := recover(); r != nil {
				var zero T
				h.complete(zero, &PanicError{Value: r})
			}
		}()
		h.complete(f(w), nil)
	}
}

// Join waits for the task to finish, from INSIDE a task: while it waits,
// w keeps running other tasks. Returns the result, or a *PanicError.
func (h *Handle[T]) Join(w *Worker) (T, error) {
	w.helpUntil(h.finished)
	if !h.finished() {
		var zero T
		return zero, errors.New("helpUntil returned before the task finished")
	}
	return h.val, h.err
}

// Wait blocks until the task has finished, from OUTSIDE the pool (a plain
// goroutine). Don't call it from inside a task: that worker would sleep.
func (h *Handle[T]) Wait() (T, error) {
	<-h.done
	return h.val, h.err
}

// Done is closed once the result is there.
func (h *Handle[T]) Done() <-chan struct{} { return h.done }

// Run runs f on the pool from outside. If the pool refused it
// (ErrPoolClosed), Wait returns that error.
func Run[T any](p *WorkStealingPool, f func(w *Worker) T) *Handle[T] {
	h := newHandle[T]()
	if err := p.Execute(h.run(f)); err != nil {
		var zero T
		h.complete(zero, err)
	}
	return h
}

// Fork runs f as a new task on w's deque and returns a Handle for it.
// Call it from inside a task; Join the handle later.
func Fork[T any](w *Worker, f func(w *Worker) T) *Handle[T] {
	h := newHandle[T]()
	w.Fork(h.run(f))
	return h
}
