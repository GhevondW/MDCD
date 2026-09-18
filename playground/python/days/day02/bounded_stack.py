"""Day 2 -- Contracts & Invariants.

A stack of ints with a fixed capacity, given once in the constructor.
The stack must protect its own rules:

  - push(v): only allowed when the stack is not full. If it is already
    full, raise IndexError and do not change anything.
  - pop(): only allowed when the stack is not empty. If it is empty,
    raise IndexError. Otherwise remove and return the top value.
  - top(): same rule as pop(), but the value stays on the stack.
  - full(), empty(), size(): report the current state.

At any moment between calls, 0 <= size() <= capacity must hold.
"""


class BoundedStack:
    def __init__(self, capacity: int) -> None:
        # TODO: choose your own representation.
        pass

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
        # TODO
        pass

    def pop(self) -> int:
        # TODO
        return 0

    def top(self) -> int:
        # TODO
        return 0
