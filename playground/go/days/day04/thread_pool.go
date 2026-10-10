package day04

// Day 4 -- A thread pool (goroutine pool).
//
// A fixed number of worker goroutines take tasks from a queue and run
// them -- "keep a few goroutines busy with many tasks". You write the
// pool itself; the result handle (Future, Submit) is given to you below.
//
// Build it from sync.Mutex / sync.Cond (or your own BoundedBlockingQueue
// from problem 1 -- but you don't have to; this file must work on its
// own). Do NOT use a buffered channel as the task queue, and do not
// busy-wait.
//
// The contract:
//
//   - NewThreadPool(workers, queueCapacity): starts `workers` goroutines
//     right away. queueCapacity is the largest number of tasks that may
//     be WAITING (accepted, but not yet picked up by a worker); tasks a
//     worker is running don't count. 0 means unbounded. Panics if
//     workers < 1 or queueCapacity < 0.
//   - Execute(task): hand a task to the pool. If the queue is full, wait
//     until there is room (backpressure). Returns nil once the task is
//     accepted -- and an accepted task is ALWAYS run, exactly once.
//     Returns ErrPoolClosed (task not run) if the pool is shut down,
//     whether before the call or while Execute was waiting for room.
//   - TrySubmit(task): never waits. Returns nil if accepted,
//     ErrQueueFull if the queue is full, ErrPoolClosed if shut down.
//   - Shutdown(): stop accepting tasks (Execute / TrySubmit return
//     ErrPoolClosed from now on), let the workers run EVERY task that was
//     already accepted, then stop them, and return only when all workers
//     have exited. Safe to call more than once and from many goroutines.
//     Don't call it from inside a task of the same pool (it would wait
//     for itself).
//   - A task that panics must not kill its worker, and must not stop the
//     pool: recover it and carry on with the next task.
//   - Tasks start in the order they were accepted (FIFO); with one worker
//     they therefore also run in that order.
//   - At most `workers` tasks run at the same time, and with enough
//     tasks all workers are busy at once.
//
// Provided (don't change): Future[T], PanicError and Submit, which wrap
// Execute and give the caller the task's result -- or its panic as an
// error -- the way std::future / Java's Future do.

import (
	"errors"
	"fmt"
)

var (
	ErrPoolClosed = errors.New("pool is shut down")
	ErrQueueFull  = errors.New("queue is full")

	// errNotImplemented is what the unfinished skeleton returns. Delete
	// it once nothing uses it any more.
	errNotImplemented = errors.New("TODO: not implemented yet")
)

type ThreadPool struct {
	// TODO: choose your own representation.
}

func NewThreadPool(workers, queueCapacity int) *ThreadPool {
	// TODO
	return &ThreadPool{}
}

func (p *ThreadPool) Execute(task func()) error {
	// TODO
	return errNotImplemented
}

func (p *ThreadPool) TrySubmit(task func()) error {
	// TODO
	return errNotImplemented
}

func (p *ThreadPool) Shutdown() {
	// TODO
}

// ---- provided: results, and panics as errors ---------------------------

// PanicError is what a Future reports when its task panicked.
type PanicError struct{ Value any }

func (e *PanicError) Error() string { return fmt.Sprintf("task panicked: %v", e.Value) }

// Future is the read end of a one-time channel for one value -- or one
// error -- from a task.
type Future[T any] struct {
	done chan struct{}
	val  T
	err  error
}

func newFuture[T any]() *Future[T] { return &Future[T]{done: make(chan struct{})} }

func (f *Future[T]) complete(v T, err error) {
	f.val, f.err = v, err
	close(f.done)
}

// Get waits until the task has finished and returns its result. err is a
// *PanicError if the task panicked, or the error from Execute (for
// example ErrPoolClosed) if the pool refused it.
func (f *Future[T]) Get() (T, error) {
	<-f.done
	return f.val, f.err
}

// Done is closed once the result is there.
func (f *Future[T]) Done() <-chan struct{} { return f.done }

// Submit runs f on the pool and returns a Future for its result. (A
// package-level function, because Go methods can't have type parameters.)
func Submit[T any](p *ThreadPool, f func() T) *Future[T] {
	fut := newFuture[T]()
	err := p.Execute(func() {
		defer func() {
			if r := recover(); r != nil {
				var zero T
				fut.complete(zero, &PanicError{Value: r})
			}
		}()
		fut.complete(f(), nil)
	})
	if err != nil {
		var zero T
		fut.complete(zero, err)
	}
	return fut
}
