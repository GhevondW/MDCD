"""Day 3 -- The bounded stack, shared.

Day 2's bounded stack, with one change: many threads use it at the
same time. The contract is the same:

  - push(v): only allowed when the stack is not full. If it is already
    full, raise IndexError and do not change anything.
  - pop(): only allowed when the stack is not empty. If it is empty,
    raise IndexError. Otherwise remove and return the top value.
  - top(): same rule as pop(), but the value stays on the stack.
  - full(), empty(), size(): report the current state.

Day 2 said: the invariant 0 <= size() <= capacity may be broken for a
moment *inside* a method, as long as it is restored before the method
returns. With two threads, another thread can arrive in exactly that
moment. No thread may ever see or cause a half-done push or pop: no
value lost, no value returned twice, no value made up.
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
