# Day 3 — Practical Tasks

Implement the `TODO`s in the headers under `include/day03/`; don't change
the test files. Problems 1–3 are the core set; 4–6 are challenge problems.
Each header starts with the full problem statement.

Every class here is used by many threads at the same time. The tests start
real threads and hammer your code. A test that passes once doesn't prove
the code is race-free; a test that fails even once proves it isn't.

Run just this day's tests from `playground/cpp/build`:
```
ctest -C Release -R "^day03_"
```

Each test has a 60-second limit, so a deadlock fails the test instead of
hanging it.

To have ThreadSanitizer check your code too (Linux/macOS, clang or gcc),
configure a separate build with:
```
cmake -S . -B build-tsan -DCMAKE_CXX_FLAGS="-fsanitize=thread -g" -DCMAKE_EXE_LINKER_FLAGS="-fsanitize=thread"
cmake --build build-tsan
cd build-tsan && ctest -R "^day03_" --output-on-failure
```
Every `WARNING: ThreadSanitizer: data race` it prints is a real bug, even
when the test itself passed.

## Core

### 1. Ticket Dispenser — Races
`include/day03/ticket_dispenser.h`

The lecture's broken counter. `next()` hands out 1, 2, 3, ... — never the
same number twice and never skipping one, no matter how many threads call
it at once.

### 2. Bank Account — Check-then-act
`include/day03/account.h`

`withdraw` succeeds only if the balance covers it, so the balance never
goes below zero. The check and the subtraction must happen as one step.

### 3. Bounded Stack, shared — Invariants under concurrency
`include/day03/bounded_stack.h`

Day 2's bounded stack with the same contract, now used by many threads at
once. No value may be lost, returned twice, or made up.

## Challenge

### 4. Compute Once per Key — Singleton, fixed
`include/day03/once_cache.h`

A cache that computes each key's value at most once, even when many
threads ask for the same missing key at the same moment. Different keys
must be able to compute at the same time, so one big lock won't do. The
full semantics are in the header.

### 5. Thread-Safe Event Bus — Observer, fixed
`include/day03/event_bus.h`

Day 1's publisher, safe for many threads and for callbacks that
subscribe, unsubscribe or publish from inside a callback. Each publish
delivers to a snapshot of the subscribers, and no lock may be held while
a callback runs. The full semantics are in the header.

### 6. Interval Booker, shared — Check-then-act
`include/day03/concurrent_booker.h`

Day 2's interval booker for many threads. `bookFirstFree` finds and
books a slot in one step, and `bookAll` books a whole group of intervals
or none of them. The full semantics are in the header.
