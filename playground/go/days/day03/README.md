# Day 3 — Practical Tasks (Go)

Implement the `TODO`s in the `.go` files below; don't change the
`_test.go` files. Problems 1–3 are the core set; 4–6 are challenge
problems. Each file's doc comment starts with the full problem statement.

Every type here is used by many goroutines at the same time. The tests
start real goroutines and hammer your code. A test that passes once
doesn't prove the code is race-free; a test that fails even once proves
it isn't.

Run just this day's tests from `playground/go`:
```
go test ./days/day03/...
```

Each test has a 60-second limit, so a deadlock stops the run with a
panic that names the test and prints where every goroutine is stuck,
instead of hanging it. Look for goroutines stuck in `sync.(*Mutex).Lock`:
a Go mutex is not reentrant, so a method that holds the lock must not
call another method that takes the same lock.

To have Go's race detector check your code too:
```
go test -race ./days/day03/...
```
Every `WARNING: DATA RACE` it prints is a real bug, even when the test
itself passed. (`-race` needs a 64-bit OS and, except on macOS, cgo with
a C compiler such as gcc.) Without `-race`, the Go runtime still catches
some races on its own: `fatal error: concurrent map writes` (or
`concurrent map read and map write`) means two goroutines used a map at
the same time — it stops the whole test run.

## Core

### 1. Ticket Dispenser — Races
`ticket_dispenser.go`

The lecture's broken counter. `Next()` hands out 1, 2, 3, ... — never the
same number twice and never skipping one, no matter how many goroutines
call it at once.

### 2. Bank Account — Check-then-act
`account.go`

`Withdraw` succeeds only if the balance covers it, so the balance never
goes below zero. The check and the subtraction must happen as one step.

### 3. Bounded Stack, shared — Invariants under concurrency
`bounded_stack.go`

Day 2's bounded stack with the same contract, now used by many goroutines
at once. No value may be lost, returned twice, or made up.

## Challenge

### 4. Compute Once per Key — Singleton, fixed
`once_cache.go`

A cache that computes each key's value at most once, even when many
goroutines ask for the same missing key at the same moment. Different
keys must be able to compute at the same time, so one big lock won't do.
The full semantics are in the file's doc comment.

### 5. Thread-Safe Event Bus — Observer, fixed
`event_bus.go`

Day 1's publisher, safe for many goroutines and for callbacks that
subscribe, unsubscribe or publish from inside a callback. Each publish
delivers to a snapshot of the subscribers, and no lock may be held while
a callback runs. The full semantics are in the file's doc comment.

### 6. Interval Booker, shared — Check-then-act
`concurrent_booker.go`

Day 2's interval booker for many goroutines. `BookFirstFree` finds and
books a slot in one step, and `BookAll` books a whole group of intervals
or none of them. The full semantics are in the file's doc comment.
