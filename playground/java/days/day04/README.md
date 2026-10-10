# Day 4 — Practical Tasks (Java)

Implement the `TODO`s in `src/main/java/day04/`; don't change
`src/test/`. Problems 1–2 are the core set; 3 is the challenge problem.
Each class's Javadoc starts with the full problem statement.

Day 3 was about keeping threads out of each other's way. Day 4 is about
threads that have to wait for each other: for room in a queue, for work to
arrive, for a result. A waiting thread must **sleep**, not spin — build
everything from a lock and conditions (`synchronized` with
`wait`/`notifyAll`, or `ReentrantLock` with `Condition`s). The ready-made
queues, executors and pools in `java.util.concurrent`
(`ArrayBlockingQueue`, `ExecutorService`, `ForkJoinPool`, ...) are exactly
what you are building, so don't use them. (`Future`/`FutureTask` as a
holder for a result is fine, and the tests use the rest of the package.)

Run just this day's tests from `playground/java`:
```
mvn -pl days/day04 test
```

Each test has a 60-second limit, so a deadlock fails the test instead of
hanging it. The tests start real threads and hammer your code; a test that
passes once doesn't prove the code is correct, so run them several times.
On a fresh clone they fail, and some take a few seconds each to give up
waiting for a task that never runs — that is expected.

Some tests measure time (a timeout that must expire, four sleeping tasks
that must overlap). The margins are wide, but a heavily loaded machine can
still trip them: if one fails once and never again, look at the message
before you suspect your code.

## Core

### 1. Bounded Blocking Queue — Waiting, wake-ups, shutdown
`src/main/java/day04/BoundedBlockingQueue.java`

The lecture's queue for many producers and many consumers: one lock, two
conditions. `push` waits while the queue is full, `pop` waits while it is
empty, `pushFor` and `popFor` wait only up to a deadline, `tryPush` and
`tryPop` never wait. `close()` is the shutdown story: no more items go in,
the ones inside can still be drained, and every sleeping thread must wake
up. Watch for `if` instead of `while` around a wait, for waking the wrong
condition, and for a timeout that starts over every time the thread is
woken.

### 2. Thread Pool — Workers, backpressure, orderly shutdown
`src/main/java/day04/ThreadPool.java`

A fixed set of worker threads running tasks from a queue that may be
bounded: `execute` waits for room (backpressure), `trySubmit` refuses
instead. A task that throws must not kill its worker. `shutdown()` runs
every task already accepted, stops the workers, and waits until they are
gone. `submit` (which returns a `Future`) is given to you.

## Challenge

### 3. Work-Stealing Pool — Fewer fights, and a join that doesn't sleep
`src/main/java/day04/WorkStealingPool.java`

The lecture's last slide. Every worker gets its own deque: it pushes and
pops at the newest end, while an idle worker steals from the oldest end of
someone else's. Tasks from outside go to a shared injection queue.
`join()` on a worker doesn't sleep — it keeps running other tasks until the
awaited one is done, so recursive fork/join code (fib, parallel sum) works
even with one worker, where a plain pool deadlocks. The hard part is
`shutdown()`: tasks that are still being spawned by running tasks must
finish too, and the workers may only stop when no task is left anywhere.
The full semantics are in the class Javadoc.

How would you measure that this beats one shared queue? Count the steals
(`steals()`), and try the tests' skewed workload with a pool that has just
one lock for everything.
