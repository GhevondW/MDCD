# Day 4 — Practical Tasks (Python)

Implement the `TODO`s in the `.py` files below; don't change the
`test_*.py` files or `_concurrency.py`. Problems 1–2 are the core set; 3
is the challenge problem. Each module's docstring starts with the full
problem statement.

The tests start real threads and hammer your code, and they check things
that only show up when threads wait for each other: a lost wake-up, a
wake-up for nobody, a worker that never stops. A test that passes once
doesn't prove the code is right; a test that fails even once proves it
isn't.

Run just this day's tests from `playground/python`:
```
python3 -m unittest discover -s days/day04 -t . -p "test_*.py" -v
```
(`-t .` matters — without it, the relative imports in the test files fail.)
One problem at a time, e.g. `python3 -m unittest days.day04.test_thread_pool -v`.

A fresh clone **fails** every test here and takes a minute or two to do
it, because the tests wait for things that never happen. Once your code is
right, the whole day runs in a few seconds.

## Rules

Build everything from `threading` (`Lock`, `Condition`, `Event`, ...).
Not allowed: `queue.Queue`, `concurrent.futures.ThreadPoolExecutor`,
`multiprocessing.pool.ThreadPool`, `asyncio`. A plain `collections.deque`
or list for the data is fine, and so is `concurrent.futures.Future` as the
holder for a result (it is already used in the code you are given).

Waiting threads must **sleep** — a loop that keeps asking "is it my turn
yet?" burns a whole core and is the busy waiting from the start of the
lecture. The tests cannot see that, so check your own code: a program with
blocked consumers should use almost no CPU while it waits.

## Why the tests slow your code down

CPython runs only one thread at a time (the GIL) and switches threads
only at certain points, so a race can hide for a long time. For the queue
and the thread pool, the tests run your code in `interleaved()` mode: before
each line of your solution files, the thread briefly gives up the GIL, so
another thread can run in between — the way real threads on a multi-core
machine can. The GIL is not a lock for your data: you still need a lock.
(The GIL also means Problem 3 cannot make CPU-bound Python code faster;
its tests use tasks that sleep, which is what shows the pool is working.)

A deadlocked test fails with a message showing where a thread is stuck. If
a test hangs anyway, the run stops after 30 seconds and shows where that
test is stuck.

## Core

### 1. Bounded Blocking Queue — Waiting, wake-ups, shutdown
`bounded_blocking_queue.py`

The structure from the lecture: a queue of limited size, many producers and
many consumers. `push` waits while it is full, `pop` waits while it is
empty. Items come out in the order they went in. `close()` wakes everyone
who is waiting; consumers first take what is left, then get "nothing".
Also: `try_push` / `try_pop` that never wait, and `push_for` / `pop_for`
that give up after a time limit — one deadline for the whole call.

### 2. Thread Pool — Workers, backpressure, shutdown
`thread_pool.py`

A few worker threads and one queue of tasks. `execute` hands over a task
and waits if the queue is full; `try_submit` says no instead of waiting.
`shutdown()` finishes every accepted task, then stops the workers. A task
that raises must not kill its worker. `submit` (already written) gives
back a `Future` with the result or the exception.

## Challenge

### 3. Work-Stealing Pool — Per-worker queues, stealing, helping
`work_stealing_pool.py`

Task 2's pool has one queue, one lock, and deadlocks when tasks wait for
tasks. Here every worker has its own deque (newest task first for the
owner), an idle worker steals the oldest task of another, and a task that
waits for a result keeps running other tasks meanwhile — so recursive
"fork, then wait" code works even on one worker. Shutdown has to finish
work that running tasks create while it is shutting down.
