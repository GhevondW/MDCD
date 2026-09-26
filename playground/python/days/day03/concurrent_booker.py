"""Day 3 -- Challenge: Interval Booker, shared.

Day 2's interval booker, now used by many threads at the same time --
many people booking one meeting room at once. Intervals are half-open,
as on Day 2: [start, end) includes start but not end, so [0,10) and
[10,20) do NOT overlap. The one rule that must always hold: no two
recorded bookings overlap.

On Day 2, booking "the first free slot" took two calls:

    s = booker.first_free(from_, d)   # check
    booker.book(s, s + d)             # act

That is check-then-act: another thread can take the slot in between.
This booker has operations that find and book as ONE step instead.

  book(start, end)     if [start, end) overlaps no booking, record it
                       and return True. Otherwise return False and
                       change nothing. start must be < end; if not,
                       raise ValueError.
  cancel(start, end)   remove the booking with exactly this start and
                       end. Returns whether a booking was removed.
  book_first_free(from_, duration)
                       find the smallest start >= from_ such that
                       [start, start + duration) overlaps no booking,
                       book it, and return start -- as one step.
                       duration must be > 0; if not, raise ValueError.
  book_all(intervals)  book every interval in the list, or none of
                       them. If any one overlaps an existing booking or
                       another interval in the same list, return False
                       and change nothing. An empty list returns True.
                       If any interval has start >= end, raise
                       ValueError and change nothing.
  bookings()           every booking, as a list of Interval, sorted by
                       start.
  count()              how many bookings are recorded.

Every method may be called by many threads at the same time.
"""

from typing import List, NamedTuple


class Interval(NamedTuple):
    start: int
    end: int


class ConcurrentBooker:
    def __init__(self) -> None:
        # TODO: choose your own representation.
        pass

    def book(self, start: int, end: int) -> bool:
        # TODO
        return False

    def cancel(self, start: int, end: int) -> bool:
        # TODO
        return False

    def book_first_free(self, from_: int, duration: int) -> int:
        # TODO
        return 0

    def book_all(self, intervals: List[Interval]) -> bool:
        # TODO
        return False

    def bookings(self) -> List[Interval]:
        # TODO
        return []

    def count(self) -> int:
        # TODO
        return 0
