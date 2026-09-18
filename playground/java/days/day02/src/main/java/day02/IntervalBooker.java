package day02;

/**
 * Day 2 -- Challenge: Interval Booker.
 *
 * A booking calendar over HALF-OPEN intervals [start, end): the invariant
 * is that no two recorded bookings ever overlap. [0,10) and [10,20) do
 * not overlap.
 *
 *   book(start, end)          record the booking and return true iff it
 *                             overlaps no existing booking; otherwise
 *                             return false and change nothing.
 *                             Contract: start < end, else throw
 *                             IllegalArgumentException.
 *   cancel(start, end)        remove a booking by EXACT match only.
 *   firstFree(from, duration) the smallest start >= from such that
 *                             [start, start + duration) overlaps no
 *                             existing booking.
 *                             Contract: duration > 0, else throw
 *                             IllegalArgumentException.
 *   count()                   number of recorded bookings.
 */
public class IntervalBooker {
    // TODO: choose your own representation.

    public boolean book(long start, long end) {
        // TODO
        return false;
    }

    public boolean cancel(long start, long end) {
        // TODO
        return false;
    }

    public long firstFree(long from, long duration) {
        // TODO
        return 0;
    }

    public int count() {
        // TODO
        return 0;
    }
}
