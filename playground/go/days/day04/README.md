# Day 4 — Practical Tasks (Go)

Implement the `TODO`s in the `.go` files below; don't change the
`_test.go` files. Problems 1–2 are the core set; 3 is the challenge.
Each file's doc comment starts with the full problem statement.

Today's types make goroutines *wait* — for room, for an item, for a task —
without burning CPU. Build them from `sync.Mutex` and `sync.Cond` (the
"state + a lock + a waiting room" recipe from the lecture). Using a Go
channel as the queue is off limits: a buffered channel already *is* a
bounded blocking queue, which would do the exercise for you. The tests
start real goroutines and hammer your code, and many of them check
*timing* ("this call must still be waiting after 150 ms", "this must give
up after 100 ms"), so a solution that spins or sleeps in a loop will not
only be wasteful — it can fail them.

Run just this day's tests from `playground/go`:
```
go test ./days/day04/...
```

Each test has a 60-second limit, so a deadlock stops the run with a
panic that names the test and prints where every goroutine is stuck,
instead of hanging it. Look for goroutines stuck in `sync.(*Cond).Wait`
(nobody is going to wake them: a lost wake-up, or `Signal` where you
needed `Broadcast`) or in `sync.(*Mutex).Lock` (a Go mutex is not
reentrant).

To have Go's race detector check your code too:
```
go test -race ./days/day04/...
```
Every `WARNING: DATA RACE` it prints is a real bug, even when the test
itself passed. (`-race` needs a 64-bit OS and, except on macOS, cgo with
a C compiler such as gcc.)

A fresh clone compiles but **fails** most tests — expected, the `TODO`s
aren't filled in yet.

## Core

### 1. Bounded Blocking Queue — Waiting with a condition variable
`bounded_blocking_queue.go`

The structure to know cold: a fixed-capacity FIFO queue for many
producers and many consumers. `Push` waits while it is full, `Pop` waits
while it is empty, `TryPush` / `TryPop` never wait, and `PushFor` /
`PopFor` wait up to a deadline. `Close()` is the shutdown story: no more
pushes, every sleeper wakes up, consumers drain what is left and then
get `false`. `sync.Cond` has no timed wait — the file's doc comment
shows the trick.

### 2. Thread Pool — Workers, a task queue, shutdown
`thread_pool.go`

A fixed set of worker goroutines running tasks from a queue. `Execute`
waits when the queue is full (backpressure), `TrySubmit` refuses, and
`Shutdown` runs every accepted task before it stops the workers. A task
that panics must not take its worker down. `Future` and `Submit` are
given: they turn `Execute` into "call this and give me the result".

## Challenge

### 3. Work-Stealing Pool — Faster, and no deadlock
`work_stealing_pool.go`

One shared queue is one lock everyone fights for — and a task that waits
for a child task can starve the pool of workers (the lecture's
"deadlock with no locks at all"). Give every worker its own deque: it
takes its own newest task, and a worker that runs dry *steals* the oldest
task of another. `Join` doesn't sleep: while it waits, the worker runs
other tasks, so a recursive `fib` works even on a pool of one worker.
Go has no goroutine identity, so tasks receive their `*Worker` as an
argument — the file's doc comment explains how `Fork` and `Join` use it.
Shutdown must wait for tasks that were forked after it was called.
