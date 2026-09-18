"""Day 2 -- Challenge: Interval Booker.

A booking calendar. A booking takes the time interval [start, end):
`start` is included, `end` is not. So [0,10) and [10,20) do NOT overlap
-- one ends exactly where the other starts.

The one rule that must always hold: no two recorded bookings overlap.

  book(start, end)             if the interval overlaps no existing
                               booking, record it and return True.
                               Otherwise return False and change nothing.
                               start must be < end; if not, raise
                               ValueError.
  cancel(start, end)           remove a booking. Only an EXACT match
                               counts: the same start and the same end.
                               Returns whether a booking was removed.
  first_free(from_, duration)  find the earliest time a booking of
                               length `duration` could start: the
                               smallest start >= from_ such that
                               [start, start + duration) overlaps no
                               existing booking. duration must be > 0;
                               if not, raise ValueError.
  count()                      how many bookings are recorded.
"""


class IntervalBooker:
    def __init__(self) -> None:
        # TODO: choose your own representation.
        pass

    def book(self, start: int, end: int) -> bool:
        # TODO
        return False

    def cancel(self, start: int, end: int) -> bool:
        # TODO
        return False

    def first_free(self, from_: int, duration: int) -> int:
        # TODO
        return 0

    def count(self) -> int:
        # TODO
        return 0
