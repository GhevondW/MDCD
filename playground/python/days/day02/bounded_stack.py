"""Day 2 -- Contracts & Invariants.

Implement BoundedStack so it enforces its own contract:
  - push(v): precondition "not full" -- if the stack IS full, raise
    IndexError instead of silently corrupting state.
  - pop() / top(): precondition "not empty" -- same idea.
  - Invariant to hold at all times between calls: 0 <= size() <= capacity.
"""


class BoundedStack:
    def __init__(self, capacity: int) -> None:
        self._capacity = capacity
        self._data: list[int] = []

    def full(self) -> bool:
        # TODO
        return False

    def empty(self) -> bool:
        # TODO
        return False

    def size(self) -> int:
        # TODO
        return 0

    def push(self, value: int) -> None:
        # TODO: raise IndexError("full") if full(), otherwise store value.
        pass

    def pop(self) -> int:
        # TODO: raise IndexError("empty") if empty(), otherwise remove and
        # return the top value.
        return 0

    def top(self) -> int:
        # TODO: same precondition as pop(), but don't remove anything.
        return 0
