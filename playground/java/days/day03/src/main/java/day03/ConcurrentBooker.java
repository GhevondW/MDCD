package day03;

import java.util.List;

/**
 * Day 3 -- Challenge: Interval Booker, shared.
 *
 * Day 2's interval booker, now used by many threads at the same time --
 * many people booking one meeting room at once. Intervals are half-open,
 * as on Day 2: [start, end) includes start but not end, so [0,10) and
 * [10,20) do NOT overlap. The one rule that must always hold: no two
 * recorded bookings overlap.
 *
 * On Day 2, booking "the first free slot" took two calls:
 *
 *     long s = booker.firstFree(from, d);   // check
 *     booker.book(s, s + d);                // act
 *
 * That is check-then-act: another thread can take the slot in between.
 * This booker has operations that find and book as ONE step instead.
 *
 *   book(start, end)     if [start, end) overlaps no booking, record it
 *                        and return true. Otherwise return false and
 *                        change nothing. start must be < end; if not,
 *                        throw IllegalArgumentException.
 *   cancel(start, end)   remove the booking with exactly this start and
 *                        end. Returns whether a booking was removed.
 *   bookFirstFree(from, duration)
 *                        find the smallest start >= from such that
 *                        [start, start + duration) overlaps no booking,
 *                        book it, and return start -- as one step.
 *                        duration must be > 0; if not, throw
 *                        IllegalArgumentException.
 *   bookAll(intervals)   book every interval in the list, or none of
 *                        them. If any one overlaps an existing booking or
 *                        another interval in the same list, return false
 *                        and change nothing. An empty list returns true.
 *                        If any interval has start >= end, throw
 *                        IllegalArgumentException and change nothing.
 *   bookings()           every booking, sorted by start.
 *   count()              how many bookings are recorded.
 *
 * Every method may be called by many threads at the same time.
 */
public class ConcurrentBooker {
    public record Interval(long start, long end) {}

    // TODO: choose your own representation.

    public boolean book(long start, long end) {
        // TODO
        return false;
    }

    public boolean cancel(long start, long end) {
        // TODO
        return false;
    }

    public long bookFirstFree(long from, long duration) {
        // TODO
        return 0;
    }

    public boolean bookAll(List<Interval> intervals) {
        // TODO
        return false;
    }

    public List<Interval> bookings() {
        // TODO
        return List.of();
    }

    public int count() {
        // TODO
        return 0;
    }
}
