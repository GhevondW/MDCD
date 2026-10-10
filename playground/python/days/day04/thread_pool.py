"""Day 4 -- A thread pool.

A few worker threads, many tasks. Build the workers and the task queue
yourself, from `threading` -- `concurrent.futures.ThreadPoolExecutor`,
`multiprocessing.pool.ThreadPool` and `queue.Queue` are off the table.
(You may reuse your own BoundedBlockingQueue from bounded_blocking_queue.py
if you like; this task does not need it, and its tests do not depend on
it.)

  - ThreadPool(workers, queue_capacity=0): starts `workers` threads right
    away. `workers` must be at least 1 and `queue_capacity` at least 0,
    otherwise raise ValueError. queue_capacity is how many tasks may WAIT
    in the queue (tasks that are already running do not count). 0 means
    unlimited.

  - execute(task): `task` is a function with no arguments. Hands it to a
    worker. If the queue is full, wait until there is room (that is
    backpressure: a full queue slows down whoever fills it). Raises
    RuntimeError if the pool has been shut down -- also if the shutdown
    happens while execute is waiting for room; in that case the task does
    not run.
  - try_submit(task) -> bool: never waits. True if the task was accepted;
    False if the queue is full or the pool has been shut down.
  - shutdown(): stop accepting tasks (execute raises RuntimeError from
    now on), let the workers run EVERY task that was accepted before,
    then wait until all the worker threads have finished. Calling it again
    (or from several threads at once) is fine and does nothing more.
    The pool is also a context manager: leaving a `with` block calls
    shutdown().
  - worker_count() -> int.

A task that raises must not kill its worker: swallow the exception and
carry on with the next task. (submit(), below, already catches errors and
hands them to the caller through the Future, so this is about tasks given
to execute() directly.)

Already written for you:

  - submit(fn, *args, **kwargs) -> Future: runs fn(*args, **kwargs) on a
    worker. The returned `concurrent.futures.Future` gets the return
    value -- or the exception. It calls execute(), so it blocks while the
    queue is full and raises RuntimeError after shutdown().

Tasks may call execute() themselves (a task can hand out more work), but
a task that WAITS for the result of another task of the same pool can
block that pool forever -- see "Deadlock with no locks at all" in the
lecture. Task 3 is the pool that fixes that.

Things to think about: what should a worker do when the queue is empty,
and how does it find out that it is time to stop? Whom must shutdown()
wake up? What if the queue is full and shutdown is called?
"""

import threading
from collections import deque
from concurrent.futures import Future
from typing import Any, Callable


class ThreadPool:
    def __init__(self, workers: int, queue_capacity: int = 0) -> None:
        # TODO: validate the arguments, then choose your own representation
        # and start the worker threads.
        pass

    def execute(self, task: Callable[[], Any]) -> None:
        # TODO
        pass

    def try_submit(self, task: Callable[[], Any]) -> bool:
        # TODO
        return False

    def shutdown(self) -> None:
        # TODO
        pass

    def worker_count(self) -> int:
        # TODO
        return 0

    # ---- already written for you ------------------------------------

    def submit(self, fn: Callable[..., Any], *args: Any, **kwargs: Any) -> Future:
        future = Future()

        def run() -> None:
            if not future.set_running_or_notify_cancel():
                return
            try:
                result = fn(*args, **kwargs)
            except BaseException as e:  # travels to whoever asks the Future
                future.set_exception(e)
            else:
                future.set_result(result)

        self.execute(run)
        return future

    def __enter__(self) -> "ThreadPool":
        return self

    def __exit__(self, *exc_info: Any) -> None:
        self.shutdown()
