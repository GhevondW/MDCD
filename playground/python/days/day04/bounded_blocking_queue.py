"""Day 4 -- The bounded blocking queue (many producers, many consumers).

The structure from the lecture: some state, a lock, and condition
variables. Build it from `threading.Lock` / `threading.Condition` only --
`queue.Queue`, `collections.deque` with `maxlen` tricks, semaphores and
the like are off the table. (A plain `collections.deque` or `list` for
the items is fine.)

  - BoundedBlockingQueue(capacity): capacity must be at least 1,
    otherwise raise ValueError.

  - push(item) -> bool: while the queue is full, wait. Then add the item
    and return True. If the queue is closed -- or becomes closed while
    push is waiting -- return False and do not add the item.
  - pop() -> (ok, item): while the queue is empty, wait. Then remove the
    oldest item and return (True, item). If the queue is closed and
    nothing is left, return (False, None). A closed queue still hands out
    the items that are already in it.

    Why a pair? `None` is a perfectly good item, so "no item" cannot be
    signalled with None. Every call that removes an item returns
    (ok, item).

  - try_push(item) -> bool: never waits. True if the item was added;
    False if the queue is full or closed.
  - try_pop() -> (ok, item): never waits. (True, item) if there was an
    item; (False, None) if the queue is empty.

  - push_for(item, timeout) -> bool and pop_for(timeout) -> (ok, item):
    like push and pop, but give up after `timeout` seconds (a float) and
    return False / (False, None). The limit is for the WHOLE call: being
    woken up and finding the queue still full (or empty) again must not
    start the timer over. A close() ends the wait at once. A timeout of
    0 or less behaves like try_push / try_pop.

  - close(): no more items may be added. Wake up EVERY thread that is
    waiting, in push or pop: producers return False, consumers first take
    what is left and then get (False, None). Calling close() again does
    nothing.
  - closed() -> bool, size() -> int, capacity() -> int.

Order: items come out in the order they went in (FIFO). Many threads call
everything at the same time; nothing may be lost, handed out twice, or
made up, and size() must never be seen above capacity().

Threads that wait must SLEEP -- no loops that keep asking "is it my turn
yet?". Things to think about: after being woken up, is the thing you
waited for still true (another thread may have been faster, and wake-ups
can come for no reason)? Whom does push wake up, and whom does pop? What
must close() wake up?
"""

import threading
import time
from collections import deque
from typing import Any, Tuple


class BoundedBlockingQueue:
    def __init__(self, capacity: int) -> None:
        # TODO: validate capacity, then choose your own representation.
        pass

    def push(self, item: Any) -> bool:
        # TODO
        return False

    def pop(self) -> Tuple[bool, Any]:
        # TODO
        return (False, None)

    def try_push(self, item: Any) -> bool:
        # TODO
        return False

    def try_pop(self) -> Tuple[bool, Any]:
        # TODO
        return (False, None)

    def push_for(self, item: Any, timeout: float) -> bool:
        # TODO
        return False

    def pop_for(self, timeout: float) -> Tuple[bool, Any]:
        # TODO
        return (False, None)

    def close(self) -> None:
        # TODO
        pass

    def closed(self) -> bool:
        # TODO
        return False

    def size(self) -> int:
        # TODO
        return 0

    def capacity(self) -> int:
        # TODO
        return 0
