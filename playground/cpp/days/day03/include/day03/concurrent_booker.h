#pragma once
// Day 3 -- Challenge: Interval Booker, shared
//
// Day 2's interval booker, now used by many threads at the same time --
// many people booking one meeting room at once. Intervals are half-open,
// as on Day 2: [start, end) includes start but not end, so [0,10) and
// [10,20) do NOT overlap. The one rule that must always hold: no two
// recorded bookings overlap.
//
// On Day 2, booking "the first free slot" took two calls:
//
//     long long s = booker.firstFree(from, d);   // check
//     booker.book(s, s + d);                     // act
//
// That is check-then-act across two calls -- an API race: another
// thread can take the slot in between, even if each call is locked.
// This booker has operations that find and book as ONE step instead.
//
//   book(start, end)     if [start, end) overlaps no booking, record it
//                        and return true. Otherwise return false and
//                        change nothing. start must be < end; if not,
//                        throw std::invalid_argument.
//   cancel(start, end)   remove the booking with exactly this start and
//                        end. Returns whether a booking was removed.
//   bookFirstFree(from, duration)
//                        find the smallest start >= from such that
//                        [start, start + duration) overlaps no booking,
//                        book it, and return start -- as one step.
//                        duration must be > 0; if not, throw
//                        std::invalid_argument.
//   bookAll(intervals)   book every interval in the list, or none of
//                        them. If any one overlaps an existing booking or
//                        another interval in the same list, return false
//                        and change nothing. An empty list returns true.
//                        If any interval has start >= end, throw
//                        std::invalid_argument and change nothing.
//   bookings()           every booking, sorted by start.
//   count()              how many bookings are recorded.
//
// Every method may be called by many threads at the same time.

#include <cstddef>
#include <stdexcept>
#include <vector>

struct Interval {
    long long start;
    long long end;

    bool operator==(const Interval& o) const {
        return start == o.start && end == o.end;
    }
};

class ConcurrentBooker {
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

    long long bookFirstFree(long long from, long long duration) {
        (void)from;
        (void)duration;
        // TODO
        return 0;
    }

    bool bookAll(const std::vector<Interval>& intervals) {
        (void)intervals;
        // TODO
        return false;
    }

    std::vector<Interval> bookings() const {
        // TODO
        return {};
    }

    std::size_t count() const {
        // TODO
        return 0;
    }

private:
    // TODO: choose your own representation.
};
