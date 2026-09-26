# Day 3 — Practical Tasks (Python)

Implement the `TODO`s in the `.py` files below; don't change the
`test_*.py` files or `_concurrency.py`. Problems 1–3 are the core set;
4–6 are challenge problems. Each module's docstring starts with the full
problem statement.

Every class here is used by many threads at the same time. The tests start
real threads and hammer your code. A test that passes once doesn't prove
the code is race-free; a test that fails even once proves it isn't.

Run just this day's tests from `playground/python`:
```
python3 -m unittest discover -s days/day03 -t . -p "test_*.py" -v
```
(`-t .` matters — without it, the relative imports in the test files fail.)

## Why the tests slow your code down

CPython runs only one thread at a time (the GIL), and it switches
threads only at certain points. So a race that breaks C++, Java or Go
code almost every time may break Python code only once in a while.
Rare, not impossible — and a bug that shows up only once in a while is
the hardest kind to find. To make races show up in the tests,
`_concurrency.py` runs your code in `interleaved()` mode: before each
line of your solution files, the thread briefly gives up the GIL, so
another thread can run in between — the way real threads on a
multi-core machine can. Your code runs slower in the tests because of
this, but it behaves the same. The GIL is not a lock for your data: you
still need a lock.

A deadlocked test fails with a message showing where a thread is stuck.
If a test hangs anyway (say, a thread waits for a lock it already
holds), the run stops after 30 seconds and shows where that test is
stuck.

## Core

### 1. Ticket Dispenser — Races
`ticket_dispenser.py`

The lecture's broken counter. `next()` hands out 1, 2, 3, ... — never the
same number twice and never skipping one, no matter how many threads call
it at once.

### 2. Bank Account — Check-then-act
`account.py`

`withdraw` succeeds only if the balance covers it, so the balance never
goes below zero. The check and the subtraction must happen as one step.

### 3. Bounded Stack, shared — Invariants under concurrency
`bounded_stack.py`

Day 2's bounded stack with the same contract, now used by many threads at
once. No value may be lost, returned twice, or made up.

## Challenge

### 4. Compute Once per Key — Singleton, fixed
`once_cache.py`

A cache that computes each key's value at most once, even when many
threads ask for the same missing key at the same moment. Different keys
must be able to compute at the same time, so one big lock won't do. The
full semantics are in the module docstring.

### 5. Thread-Safe Event Bus — Observer, fixed
`event_bus.py`

Day 1's publisher, safe for many threads and for callbacks that
subscribe, unsubscribe or publish from inside a callback. Each publish
delivers to a snapshot of the subscribers, and no lock may be held while
a callback runs. The full semantics are in the module docstring.

### 6. Interval Booker, shared — Check-then-act
`concurrent_booker.py`

Day 2's interval booker for many threads. `book_first_free` finds and
books a slot in one step, and `book_all` books a whole group of intervals
or none of them. The full semantics are in the module docstring.
