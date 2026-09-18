package day02;

/**
 * Day 2 -- Challenge: Interval Booker.
 *
 * A booking calendar. A booking takes the time interval [start, end):
 * `start` is included, `end` is not. So [0,10) and [10,20) do NOT
 * overlap -- one ends exactly where the other starts.
 *
 * The one rule that must always hold: no two recorded bookings overlap.
 *
 *   book(start, end)          if the interval overlaps no existing
 *                             booking, record it and return true.
 *                             Otherwise return false and change nothing.
 *                             start must be < end; if not, throw
 *                             IllegalArgumentException.
 *   cancel(start, end)        remove a booking. Only an EXACT match
 *                             counts: the same start and the same end.
 *                             Returns whether a booking was removed.
 *   firstFree(from, duration) find the earliest time a booking of length
 *                             `duration` could start: the smallest
 *                             start >= from such that
 *                             [start, start + duration) overlaps no
 *                             existing booking. duration must be > 0; if
 *                             not, throw IllegalArgumentException.
 *   count()                   how many bookings are recorded.
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
