package day02;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

class IntervalBookerTest {

    @Test
    void bookOnEmptyCalendarSucceeds() {
        IntervalBooker b = new IntervalBooker();
        assertTrue(b.book(10, 20));
        assertEquals(1, b.count());
    }

    @Test
    void adjacentIntervalsDoNotOverlap() {
        // Half-open semantics: [0,10) and [10,20) share only the boundary point.
        IntervalBooker b = new IntervalBooker();
        assertTrue(b.book(0, 10));
        assertTrue(b.book(10, 20));
        assertEquals(2, b.count());
    }

    @Test
    void partialOverlapIsRejectedBothDirections() {
        IntervalBooker b = new IntervalBooker();
        assertTrue(b.book(10, 20));
        assertFalse(b.book(15, 25));   // overlaps the tail
        assertFalse(b.book(5, 15));    // overlaps the head
        assertEquals(1, b.count());
    }

    @Test
    void containmentIsRejectedBothDirections() {
        IntervalBooker b = new IntervalBooker();
        assertTrue(b.book(10, 20));
        assertFalse(b.book(12, 18));   // inside the existing booking
        assertFalse(b.book(5, 25));    // swallows the existing booking
    }

    @Test
    void identicalIntervalIsRejected() {
        IntervalBooker b = new IntervalBooker();
        assertTrue(b.book(10, 20));
        assertFalse(b.book(10, 20));
    }

    @Test
    void rejectedBookChangesNothing() {
        IntervalBooker b = new IntervalBooker();
        assertTrue(b.book(10, 20));
        assertFalse(b.book(15, 25));
        assertEquals(1, b.count());
        assertTrue(b.book(20, 30));    // the slot the rejected book didn't take
    }

    @Test
    void emptyOrReversedIntervalThrows() {
        IntervalBooker b = new IntervalBooker();
        assertThrows(IllegalArgumentException.class, () -> b.book(10, 10));
        assertThrows(IllegalArgumentException.class, () -> b.book(20, 10));
        assertEquals(0, b.count());
    }

    @Test
    void cancelExactMatchFreesTheSlot() {
        IntervalBooker b = new IntervalBooker();
        assertTrue(b.book(10, 20));
        assertTrue(b.cancel(10, 20));
        assertEquals(0, b.count());
        assertTrue(b.book(15, 25));    // the freed range is bookable again
    }

    @Test
    void cancelRequiresExactMatch() {
        IntervalBooker b = new IntervalBooker();
        assertTrue(b.book(10, 20));
        assertFalse(b.cancel(10, 15));   // sub-range: no
        assertFalse(b.cancel(5, 25));    // super-range: no
        assertEquals(1, b.count());
    }

    @Test
    void cancelUnknownIntervalReturnsFalse() {
        IntervalBooker b = new IntervalBooker();
        assertFalse(b.cancel(10, 20));
    }

    @Test
    void firstFreeOnEmptyCalendarIsFrom() {
        IntervalBooker b = new IntervalBooker();
        assertEquals(7, b.firstFree(7, 5));
    }

    @Test
    void firstFreeSkipsBusyRangesAndTooSmallGaps() {
        IntervalBooker b = new IntervalBooker();
        b.book(10, 20);
        b.book(20, 30);
        // [5,15) collides; the only room left starts after the last booking.
        assertEquals(30, b.firstFree(5, 10));
    }

    @Test
    void firstFreeFindsGapOfExactSize() {
        IntervalBooker b = new IntervalBooker();
        b.book(0, 10);
        b.book(15, 25);
        assertEquals(10, b.firstFree(0, 5));   // the gap [10,15) fits exactly
    }

    @Test
    void firstFreeWhenFromFallsInsideABooking() {
        IntervalBooker b = new IntervalBooker();
        b.book(10, 20);
        assertEquals(20, b.firstFree(12, 5));
    }

    @Test
    void firstFreeIgnoresBookingsEntirelyBeforeFrom() {
        IntervalBooker b = new IntervalBooker();
        b.book(0, 10);
        assertEquals(50, b.firstFree(50, 5));
    }

    @Test
    void firstFreeResultIsActuallyBookable() {
        IntervalBooker b = new IntervalBooker();
        b.book(3, 9);
        b.book(12, 40);
        b.book(41, 50);
        long s = b.firstFree(0, 4);
        assertTrue(b.book(s, s + 4));
    }

    @Test
    void nonPositiveDurationThrows() {
        IntervalBooker b = new IntervalBooker();
        assertThrows(IllegalArgumentException.class, () -> b.firstFree(0, 0));
        assertThrows(IllegalArgumentException.class, () -> b.firstFree(0, -3));
    }
}
