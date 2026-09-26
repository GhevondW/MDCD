"""Day 3 -- Challenge: Compute Once per Key.

Day 1's Singleton had a line to remember:

    if Singleton._instance is None:
        Singleton._instance = Singleton()

Two threads can both see None, and both create an instance. This task
is the same problem for many keys: a cache where the value for each key
is computed at most once.

  get(key, compute)  if `key` already has a value, return it.
                     Otherwise call compute() (it takes no arguments)
                     to make one, remember it, and return it.
  size()             how many keys have a remembered value. A
                     computation that is still running does not count
                     yet.

The rules that make it hard:

  1. compute() runs AT MOST ONCE per key -- even when many threads ask
     for the same missing key at the same moment. One of them runs
     compute(); the others wait for it and return the same value.
  2. Different keys never wait for each other: while compute() runs
     for key "a", a get() for key "b" can run its own compute() at the
     same time. So you cannot hold one lock while compute() runs.
  3. If compute() raises, nothing is remembered for that key, and the
     same exception reaches the caller that ran compute(). Threads that
     were waiting for that same computation get the same exception. The
     next get() for that key calls compute() again.

Every method may be called by many threads at the same time.
"""

from typing import Callable, Generic, TypeVar

V = TypeVar("V")


class OnceCache(Generic[V]):
    def __init__(self) -> None:
        # TODO: choose your own representation.
        pass

    def get(self, key: str, compute: Callable[[], V]) -> V:
        # TODO
        return None

    def size(self) -> int:
        # TODO
        return 0
