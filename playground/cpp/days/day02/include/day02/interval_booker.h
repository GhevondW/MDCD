#pragma once
// Day 2 -- Challenge: Interval Booker
//
// A booking calendar over HALF-OPEN intervals [start, end): the invariant
// is that no two recorded bookings ever overlap. [0,10) and [10,20) do
// not overlap.
//
//   book(start, end)          record the booking and return true iff it
//                             overlaps no existing booking; otherwise
//                             return false and change nothing.
//                             Contract: start < end, else throw
//                             std::invalid_argument.
//   cancel(start, end)        remove a booking by EXACT match only.
//   firstFree(from, duration) the smallest start >= from such that
//                             [start, start + duration) overlaps no
//                             existing booking.
//                             Contract: duration > 0, else throw
//                             std::invalid_argument.
//   count()                   number of recorded bookings.

#include <cstddef>
#include <stdexcept>

class IntervalBooker {
public:
    bool book(long long start, long long end) {
        (void)start;
        (void)end;
        // TODO
        return false;
    }

    bool cancel(long long start, long long end) {
        (void)start;
        (void)end;
        // TODO
        return false;
    }

    long long firstFree(long long from, long long duration) const {
        (void)from;
        (void)duration;
        // TODO
        return 0;
    }

    std::size_t count() const {
        // TODO
        return 0;
    }

private:
    // TODO: choose your own representation.
};
