package day04

// Day 4 -- The bounded blocking queue (MPMC).
//
// The lecture's queue: a fixed capacity, many producer goroutines and
// many consumer goroutines, and goroutines that WAIT (sleep, no CPU)
// instead of spinning when they cannot go on. It is also the queue
// behind every thread pool.
//
// Build it from one sync.Mutex and sync.Cond values -- the "state + a
// lock + a waiting room" recipe from the slides. Do NOT build it on a Go
// channel (a buffered channel already IS a bounded blocking queue; that
// would be the whole exercise done for you), and do not busy-wait.
//
// The contract:
//
//   - NewBoundedBlockingQueue[T](capacity): panics if capacity <= 0.
//   - Push(v): if the queue is full, wait until there is room. Returns
//     true once v is in the queue. Returns false -- and v is NOT added --
//     if the queue is closed, whether it was closed before the call or
//     while Push was waiting.
//   - Pop(): if the queue is empty, wait until there is an item. Returns
//     (item, true). Returns (zero value, false) only when the queue is
//     closed AND empty: after Close, consumers still get every item that
//     was already in the queue ("drain"), then false.
//   - Items come out in the order they went in (FIFO). No item is lost,
//     delivered twice, or made up, and Len() never exceeds Cap().
//   - TryPush(v) / TryPop(): never wait. TryPush returns false if the
//     queue is full or closed; TryPop returns (zero value, false) if the
//     queue is empty. (Day 3's "check and act in one call".)
//   - PushFor(v, d) / PopFor(d): like Push / Pop, but give up after d in
//     total and return false. d is a DEADLINE for the whole call, not a
//     per-wake-up timeout: a goroutine that is woken up and finds it
//     still has to wait must not start a fresh d. d <= 0 behaves like
//     TryPush / TryPop. Closing the queue ends the wait early.
//   - Close(): no more pushes are accepted; every goroutine waiting in
//     Push, Pop, PushFor or PopFor must wake up and return (producers
//     false; consumers drain what is left first). Safe to call more than
//     once and from many goroutines at once.
//   - Closed(), Len(), Cap(): report the current state.
//
// Hints:
//   - Two sync.Cond values on the same mutex -- "not full" and "not
//     empty" -- keep wake-ups targeted. Re-check the condition in a
//     `for` loop around Wait, never an `if`: a woken goroutine can find
//     that someone else already took the item.
//   - Close must Broadcast (wake everyone), not Signal.
//   - sync.Cond has no timed wait. A common trick for PushFor / PopFor:
//     before waiting, start `t := time.AfterFunc(remaining, func() {
//     lock; cond.Broadcast(); unlock })`, Stop it when you are done, and
//     after every wake-up compare time.Now() with the deadline you
//     computed once at the start of the call.
//   - A Go mutex is not reentrant: don't call a locking method from
//     inside another one.

import "time"

type BoundedBlockingQueue[T any] struct {
	// TODO: choose your own representation.
}

func NewBoundedBlockingQueue[T any](capacity int) *BoundedBlockingQueue[T] {
	// TODO
	return &BoundedBlockingQueue[T]{}
}

func (q *BoundedBlockingQueue[T]) Push(v T) bool {
	// TODO
	return false
}

func (q *BoundedBlockingQueue[T]) Pop() (T, bool) {
	// TODO
	var zero T
	return zero, false
}

func (q *BoundedBlockingQueue[T]) TryPush(v T) bool {
	// TODO
	return false
}

func (q *BoundedBlockingQueue[T]) TryPop() (T, bool) {
	// TODO
	var zero T
	return zero, false
}

func (q *BoundedBlockingQueue[T]) PushFor(v T, d time.Duration) bool {
	// TODO
	return false
}

func (q *BoundedBlockingQueue[T]) PopFor(d time.Duration) (T, bool) {
	// TODO
	var zero T
	return zero, false
}

func (q *BoundedBlockingQueue[T]) Close() {
	// TODO
}

func (q *BoundedBlockingQueue[T]) Closed() bool {
	// TODO
	return false
}

func (q *BoundedBlockingQueue[T]) Len() int {
	// TODO
	return 0
}

func (q *BoundedBlockingQueue[T]) Cap() int {
	// TODO
	return 0
}
