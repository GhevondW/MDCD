"""Day 4 -- A work-stealing pool (challenge).

Task 2's pool has one queue and one lock for everybody, and a task that
waits for another task can block it forever. This pool fixes both, the
way Java's ForkJoinPool, Rust's rayon and Go's scheduler do:

  - Every worker has its OWN deque of tasks. A worker takes the NEWEST
    task from its own deque (last in, first out): the task it just
    created is the one whose data is still warm.
  - A worker whose deque is empty STEALS: it takes the OLDEST task from
    another worker's deque. (The oldest is usually the biggest piece of
    work, and the owner is farthest from needing it.)
  - A task that waits for another task does not sit idle: while it waits,
    its worker keeps running other tasks. That is why recursive
    "fork two subtasks, wait for both" code works even with ONE worker.

Build it from `threading` only -- no `concurrent.futures` executors, no
`queue.Queue`. A lock per deque, or one lock for the lot, is fine: the
point is the shape of the algorithm, not lock-free tricks. (Python's GIL
means this pool will not run CPU-bound code faster; the tests use tasks
that sleep.)

  - WorkStealingPool(workers): starts `workers` threads (at least 1,
    otherwise raise ValueError).

  - execute(task): `task` is a function with no arguments.
      * Called from one of this pool's workers (a task creating more
        work): put the task on that worker's own deque.
      * Called from any other thread: put the task on a shared
        "injection" queue, oldest first.
    Raises RuntimeError if the pool is shutting down -- but only for
    callers that are NOT workers: a running task may always add more
    work, because that work is part of what the pool already promised to
    finish.
  - How a worker picks its next task: first the newest task of its own
    deque; then the oldest task of the injection queue; then it tries to
    steal the oldest task of another worker's deque. If there is nothing
    anywhere, it sleeps -- no busy loops -- until something arrives.
  - in_worker() -> bool: is the calling thread one of this pool's
    workers?
  - help_until(done): called (by Handle.result, below) on a worker
    thread, with a function `done()` that says whether the awaited task
    has finished. Until done() is true, run other tasks, found in the same
    order as above. When there is momentarily nothing to run (the awaited
    task is being run by someone else), wait briefly -- 1 millisecond is
    fine -- and look again; do not spin.
  - steals() -> int: how many tasks were taken from another worker's
    deque so far. Tasks from the injection queue, or from a worker's own
    deque, are not steals.
  - worker_count() -> int.
  - shutdown(): stop accepting work from outside, finish EVERY task that
    was accepted -- including those that running tasks add while the pool
    is shutting down -- then wait until all workers have stopped. Calling
    it again (or from several threads) is fine. The pool is a context
    manager too.

A task that raises must not kill its worker: swallow the exception.
(submit(), below, catches errors and hands them over through the Future.)

Already written for you:

  - submit(fn) -> Handle: runs fn() as a task and returns a Handle.
  - Handle.result(timeout=None): the value fn returned -- or raises what
    fn raised. On a worker thread it first calls help_until, so waiting
    never idles a worker; anywhere else it just waits.

For example, with any number of workers, even one:

    def fib(n):
        if n < 2:
            return n
        left = pool.submit(lambda: fib(n - 1))
        right = fib(n - 2)
        return left.result() + right

Things to think about: when is it safe to stop the workers -- how do you
know that no task will ever appear again? What can a thief that looks at
another worker's deque race with? How does a sleeping worker find out
that a task has been added to somebody else's deque?
"""

import threading
import time
from collections import deque
from concurrent.futures import Future
from typing import Any, Callable, Optional


class WorkStealingPool:
    def __init__(self, workers: int) -> None:
        # TODO: validate the argument, then choose your own representation
        # and start the worker threads.
        pass

    def execute(self, task: Callable[[], Any]) -> None:
        # TODO
        pass

    def in_worker(self) -> bool:
        # TODO
        return False

    def help_until(self, done: Callable[[], bool]) -> None:
        # TODO
        pass

    def steals(self) -> int:
        # TODO
        return 0

    def worker_count(self) -> int:
        # TODO
        return 0

    def shutdown(self) -> None:
        # TODO
        pass

    # ---- already written for you ------------------------------------

    def submit(self, fn: Callable[[], Any]) -> "Handle":
        future = Future()

        def run() -> None:
            if not future.set_running_or_notify_cancel():
                return
            try:
                result = fn()
            except BaseException as e:  # travels to whoever asks the Handle
                future.set_exception(e)
            else:
                future.set_result(result)

        self.execute(run)
        return Handle(self, future)

    def __enter__(self) -> "WorkStealingPool":
        return self

    def __exit__(self, *exc_info: Any) -> None:
        self.shutdown()


class Handle:
    """The read end of a submitted task (already written for you)."""

    def __init__(self, pool: WorkStealingPool, future: Future) -> None:
        self._pool = pool
        self._future = future

    def done(self) -> bool:
        return self._future.done()

    def result(self, timeout: Optional[float] = None) -> Any:
        if self._pool.in_worker():
            self._pool.help_until(self._future.done)
        return self._future.result(timeout)
