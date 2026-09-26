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

Even with every method locked, a caller who writes

    if not s.empty():
        v = s.top()
        s.pop()

has an API race: another thread can run between the calls. So the stack
also offers calls that check and act as ONE step:

  - try_push(v): if there is room, push v and return True. If the stack
    is full, return False and change nothing. Never raises.
  - try_pop(): if the stack is not empty, remove the top value and return
    it. If it is empty, return None. Never raises.
"""

from typing import Optional


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

    def try_push(self, value: int) -> bool:
        # TODO
        return False

    def try_pop(self) -> Optional[int]:
        # TODO
        return None
