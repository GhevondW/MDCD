"""Day 2 -- Challenge: Interval Booker.

A booking calendar over HALF-OPEN intervals [start, end): the invariant
is that no two recorded bookings ever overlap. [0,10) and [10,20) do
not overlap.

  book(start, end)             record the booking and return True iff it
                               overlaps no existing booking; otherwise
                               return False and change nothing.
                               Contract: start < end, else raise
                               ValueError.
  cancel(start, end)           remove a booking by EXACT match only.
  first_free(from_, duration)  the smallest start >= from_ such that
                               [start, start + duration) overlaps no
                               existing booking.
                               Contract: duration > 0, else raise
                               ValueError.
  count()                      number of recorded bookings.
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
