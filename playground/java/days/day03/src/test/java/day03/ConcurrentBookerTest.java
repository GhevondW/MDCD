package day03;

import day03.ConcurrentBooker.Interval;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.Timeout;

import java.util.ArrayList;
import java.util.Comparator;
import java.util.List;
import java.util.Random;
import java.util.concurrent.atomic.AtomicInteger;

import static day03.RunTogether.runTogether;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

@Timeout(value = 60, threadMode = Timeout.ThreadMode.SEPARATE_THREAD)
class ConcurrentBookerTest {

    // ---- single thread: the rules ----

    @Test
    void bookRejectsOverlapAndAcceptsAdjacent() {
        ConcurrentBooker b = new ConcurrentBooker();
        assertTrue(b.book(10, 20));
        assertTrue(b.book(20, 30));   // adjacent: [10,20) and [20,30) don't overlap
        assertFalse(b.book(15, 25));
        assertFalse(b.book(0, 11));
        assertEquals(2, b.count());
    }

    @Test
    void emptyOrReversedIntervalThrows() {
        ConcurrentBooker b = new ConcurrentBooker();
        assertThrows(IllegalArgumentException.class, () -> b.book(5, 5));
        assertThrows(IllegalArgumentException.class, () -> b.book(6, 5));
        assertEquals(0, b.count());
    }

    @Test
    void cancelNeedsAnExactMatch() {
        ConcurrentBooker b = new ConcurrentBooker();
        assertTrue(b.book(10, 20));
        assertFalse(b.cancel(10, 15));
        assertTrue(b.cancel(10, 20));
        assertFalse(b.cancel(10, 20));
        assertEquals(0, b.count());
    }

    @Test
    void bookingsAreSortedByStart() {
        ConcurrentBooker b = new ConcurrentBooker();
        b.book(50, 60);
        b.book(0, 10);
        b.book(20, 30);
        assertEquals(List.of(new Interval(0, 10), new Interval(20, 30), new Interval(50, 60)),
                b.bookings());
    }

    @Test
    void bookFirstFreeOnEmptyCalendarStartsAtFrom() {
        ConcurrentBooker b = new ConcurrentBooker();
        assertEquals(7, b.bookFirstFree(7, 5));
        assertEquals(List.of(new Interval(7, 12)), b.bookings());
    }

    @Test
    void bookFirstFreeSkipsBusyRangesAndSmallGaps() {
        ConcurrentBooker b = new ConcurrentBooker();
        b.book(0, 10);
        b.book(12, 20);
        assertEquals(20, b.bookFirstFree(0, 5));  // the gap [10,12) is too small
        assertEquals(10, b.bookFirstFree(0, 2));  // ...but fits a booking of 2
        assertEquals(4, b.count());
    }

    @Test
    void bookFirstFreeWhenFromIsInsideABooking() {
        ConcurrentBooker b = new ConcurrentBooker();
        b.book(10, 20);
        assertEquals(20, b.bookFirstFree(15, 5));
    }

    @Test
    void bookFirstFreeRejectsNonPositiveDuration() {
        ConcurrentBooker b = new ConcurrentBooker();
        assertThrows(IllegalArgumentException.class, () -> b.bookFirstFree(0, 0));
        assertThrows(IllegalArgumentException.class, () -> b.bookFirstFree(0, -1));
        assertEquals(0, b.count());
    }

    @Test
    void bookAllBooksEveryInterval() {
        ConcurrentBooker b = new ConcurrentBooker();
        assertTrue(b.bookAll(List.of(new Interval(30, 40), new Interval(0, 10), new Interval(10, 20))));
        assertEquals(List.of(new Interval(0, 10), new Interval(10, 20), new Interval(30, 40)),
                b.bookings());
    }

    @Test
    void bookAllIsAllOrNothing() {
        ConcurrentBooker b = new ConcurrentBooker();
        assertTrue(b.book(10, 20));
        assertFalse(b.bookAll(List.of(new Interval(0, 5), new Interval(15, 25), new Interval(30, 40))));
        assertEquals(1, b.count(), "a failed bookAll must book nothing");
        assertTrue(b.book(0, 5));
        assertTrue(b.book(30, 40));
    }

    @Test
    void bookAllRejectsOverlapInsideTheList() {
        ConcurrentBooker b = new ConcurrentBooker();
        assertFalse(b.bookAll(List.of(new Interval(0, 10), new Interval(5, 15))));
        assertEquals(0, b.count());
    }

    @Test
    void bookAllOfEmptyListSucceeds() {
        ConcurrentBooker b = new ConcurrentBooker();
        assertTrue(b.bookAll(List.of()));
        assertEquals(0, b.count());
    }

    @Test
    void bookAllWithInvalidIntervalThrowsAndBooksNothing() {
        ConcurrentBooker b = new ConcurrentBooker();
        assertThrows(IllegalArgumentException.class,
                () -> b.bookAll(List.of(new Interval(0, 5), new Interval(7, 7))));
        assertEquals(0, b.count());
    }

    // ---- many threads ----

    @Test
    void concurrentBooksOfOverlappingSlotsHaveOneWinner() throws Exception {
        // Thread t books [t, t+16). Every pair of these overlaps, so in each
        // round exactly one thread may win.
        int threads = 16;
        for (int round = 0; round < 100; round++) {
            ConcurrentBooker b = new ConcurrentBooker();
            AtomicInteger wins = new AtomicInteger();
            runTogether(threads, t -> {
                if (b.book(t, t + 16)) wins.incrementAndGet();
            });
            int r = round;
            assertEquals(1, wins.get(), () -> "round " + r + ": book() calls that returned true");
            assertEquals(1, b.count(), () -> "round " + r + ": bookings recorded");
        }
    }

    @Test
    void concurrentBookFirstFreeTilesTheCalendar() throws Exception {
        // 8 threads each take the first free 10-unit slot 200 times. Nobody
        // cancels, so the slots must fill [0, 16000) with no gaps and no
        // overlaps: 0, 10, 20, ...
        int threads = 8;
        int perThread = 200;
        ConcurrentBooker b = new ConcurrentBooker();
        long[][] starts = new long[threads][perThread];
        runTogether(threads, t -> {
            for (int i = 0; i < perThread; i++) starts[t][i] = b.bookFirstFree(0, 10);
        });

        List<Long> all = new ArrayList<>();
        for (long[] mine : starts) {
            for (long s : mine) all.add(s);
        }
        all.sort(null);
        assertEquals(threads * perThread, all.size());
        for (int i = 0; i < all.size(); i++) {
            assertEquals(i * 10L, all.get(i), "two threads got the same slot, or a slot was skipped");
        }
        List<Interval> booked = b.bookings();
        assertEquals(all.size(), booked.size());
        assertTrue(noOverlaps(booked), "two bookings overlap");
    }

    @Test
    void concurrentBookAllHasOneWinnerAndNoPartialBookings() throws Exception {
        // Every thread's group contains the shared slot [50,60), plus two
        // slots of its own. Exactly one group can win each round -- and the
        // losers must leave none of their own slots behind.
        int threads = 8;
        for (int round = 0; round < 100; round++) {
            ConcurrentBooker b = new ConcurrentBooker();
            boolean[] won = new boolean[threads];
            runTogether(threads, t -> {
                long mine = 100 + 20L * t;
                won[t] = b.bookAll(List.of(
                        new Interval(mine, mine + 10),
                        new Interval(50, 60),
                        new Interval(1000 + mine, 1010 + mine)));
            });

            int winners = 0;
            int w = -1;
            for (int t = 0; t < threads; t++) {
                if (won[t]) {
                    winners++;
                    w = t;
                }
            }
            int r = round;
            assertEquals(1, winners, () -> "round " + r + ": bookAll() calls that returned true");
            long mine = 100 + 20L * w;
            assertEquals(List.of(
                            new Interval(50, 60),
                            new Interval(mine, mine + 10),
                            new Interval(1000 + mine, 1010 + mine)),
                    b.bookings(),
                    () -> "round " + r + ": a losing bookAll left a partial booking");
        }
    }

    @Test
    void aFailedBookAllIsNeverSeenHalfDone() throws Exception {
        // [50,60) is taken, so bookAll([[0,10), [50,60)]) must always fail and
        // book nothing -- which leaves [0,10) free for the one thread that
        // books and cancels it over and over. A bookAll that books [0,10)
        // first and takes it back once [50,60) turns out to be taken is not
        // all-or-nothing: for a moment, that thread finds [0,10) taken.
        int groupThreads = 3;
        int tries = 20000;
        ConcurrentBooker b = new ConcurrentBooker();
        assertTrue(b.book(50, 60));
        AtomicInteger groupWins = new AtomicInteger();
        AtomicInteger blocked = new AtomicInteger();
        runTogether(groupThreads + 1, t -> {
            for (int i = 0; i < tries; i++) {
                if (t < groupThreads) {
                    if (b.bookAll(List.of(new Interval(0, 10), new Interval(50, 60)))) {
                        groupWins.incrementAndGet();
                    }
                } else if (b.book(0, 10)) {
                    b.cancel(0, 10);
                } else {
                    blocked.incrementAndGet();
                }
            }
        });

        assertEquals(0, groupWins.get(), "bookAll succeeded although [50,60) was taken");
        assertEquals(0, blocked.get(),
                "book(0, 10) failed: a bookAll that failed held [0,10) for a moment");
        assertEquals(List.of(new Interval(50, 60)), b.bookings());
    }

    @Test
    void concurrentMixedOperationsKeepTheCalendarConsistent() throws Exception {
        // Threads book, cancel and bookFirstFree at random. Afterwards the
        // calendar must hold exactly the bookings that were made and not
        // cancelled -- and none may overlap.
        int threads = 8;
        int ops = 3000;
        ConcurrentBooker b = new ConcurrentBooker();
        List<List<Interval>> kept = new ArrayList<>();
        for (int t = 0; t < threads; t++) kept.add(new ArrayList<>());
        runTogether(threads, t -> {
            Random rng = new Random(1234 + t);
            List<Interval> mine = kept.get(t);
            for (int i = 0; i < ops; i++) {
                int op = rng.nextInt(3);
                if (op == 0) {
                    long s = rng.nextInt(5000);
                    long len = 1 + rng.nextInt(20);
                    if (b.book(s, s + len)) mine.add(new Interval(s, s + len));
                } else if (op == 1) {
                    long len = 1 + rng.nextInt(20);
                    long s = b.bookFirstFree(rng.nextInt(5000), len);
                    mine.add(new Interval(s, s + len));
                } else if (!mine.isEmpty()) {
                    int k = rng.nextInt(mine.size());
                    Interval gone = mine.remove(k);
                    assertTrue(b.cancel(gone.start(), gone.end()),
                            "could not cancel a booking this thread made");
                }
            }
        });

        List<Interval> expected = new ArrayList<>();
        for (List<Interval> m : kept) expected.addAll(m);
        expected.sort(Comparator.comparingLong(Interval::start));
        List<Interval> booked = b.bookings();
        assertTrue(sortedByStart(booked), "bookings() is not sorted by start");
        assertTrue(noOverlaps(booked), "two bookings overlap");
        assertEquals(expected, booked);
        assertEquals(expected.size(), b.count());
    }

    // No two bookings in a sorted list may overlap.
    private static boolean noOverlaps(List<Interval> sorted) {
        for (int i = 1; i < sorted.size(); i++) {
            if (sorted.get(i).start() < sorted.get(i - 1).end()) return false;
        }
        return true;
    }

    private static boolean sortedByStart(List<Interval> list) {
        for (int i = 1; i < list.size(); i++) {
            if (list.get(i).start() < list.get(i - 1).start()) return false;
        }
        return true;
    }
}
