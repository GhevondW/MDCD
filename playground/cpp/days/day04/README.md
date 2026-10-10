# Day 4 — Practical Tasks

Implement the `TODO`s in the headers under `include/day04/`; don't change
the test files. Problems 1–2 are the core set; 3 is the challenge. Each
header starts with the full problem statement.

The three tasks build on each other, like the slides: the queue, then the
thread pool that runs on one, then a faster kind of pool. Each one is
self-contained, so you can start any of them. Don't reach for ready-made
thread-safe queues or thread pools — building them is the point.

The tests start real threads, and most of them check that a call WAITS (and
wakes up when it should), so a solution that spins in a loop or loses one
wake-up fails them. A test that passes once doesn't prove the code is
correct; a test that fails even once proves it isn't.

Run just this day's tests from `playground/cpp/build`:
```
ctest -C Release -R "^day04_"
```

Each test has a 60-second limit, so a deadlock fails the test instead of
hanging it. The tests that start many threads print a message and stop the
test run if a thread is stuck for 30 seconds.

To have ThreadSanitizer check your code too (Linux/macOS, clang or gcc),
configure a separate build with:
```
cmake -S . -B build-tsan -DCMAKE_CXX_FLAGS="-fsanitize=thread -g" -DCMAKE_EXE_LINKER_FLAGS="-fsanitize=thread"
cmake --build build-tsan
cd build-tsan && ctest -R "^day04_" --output-on-failure
```
Every `WARNING: ThreadSanitizer: data race` it prints is a real bug, even
when the test itself passed.

## Core

### 1. Bounded Blocking Queue — Waiting
`include/day04/bounded_blocking_queue.h`

The slides' queue: a FIFO with a fixed capacity for many producers and many
consumers. `push` waits while it is full, `pop` waits while it is empty —
asleep on a condition variable, not spinning. It has non-waiting
(`tryPush`, `tryPop`) and timed (`pushFor`, `popFor`) versions, and
`close()`, which wakes every sleeper and lets the consumers drain what is
left. The full semantics are in the header.

### 2. Thread Pool — Tasks, not threads
`include/day04/thread_pool.h`

A fixed set of worker threads that run tasks from a queue. The queue can
have a limit: then `execute` waits (backpressure) and `trySubmit` says no.
`shutdown()` runs everything already accepted, then stops the workers. A
task that throws must not kill its worker. `submit()`, which returns a
`std::future`, is already written — it calls your `execute()`.

## Challenge

### 3. Work-Stealing Pool — Fast, and no deadlock
`include/day04/work_stealing_pool.h`

The pool from the end of the slides. Every worker has its own deque: it
takes its newest task, and a worker with nothing to do steals the oldest
task of another worker. A worker that waits for a result keeps running
other tasks, so recursive code (a task per call of `fib`) works even on a
pool with ONE worker — the plain pool from task 2 would deadlock on it.
Shutdown has to notice that tasks can still fork more tasks. The full
semantics are in the header.
