import random
import unittest

from ._concurrency import TimeLimitedTestCase, interleaved, run_together
from .concurrent_booker import ConcurrentBooker, Interval


def no_overlaps(sorted_bookings):
    """No two bookings in a sorted list may overlap."""
    return all(later[0] >= earlier[1]
               for earlier, later in zip(sorted_bookings, sorted_bookings[1:]))


def sorted_by_start(bookings):
    return all(earlier[0] <= later[0] for earlier, later in zip(bookings, bookings[1:]))


class ConcurrentBookerTest(TimeLimitedTestCase):
    # ---- single thread: the rules ----

    def test_book_rejects_overlap_and_accepts_adjacent(self):
        b = ConcurrentBooker()
        self.assertTrue(b.book(10, 20))
        self.assertTrue(b.book(20, 30))  # adjacent: [10,20) and [20,30) don't overlap
        self.assertFalse(b.book(15, 25))
        self.assertFalse(b.book(0, 11))
        self.assertEqual(b.count(), 2)

    def test_empty_or_reversed_interval_throws(self):
        b = ConcurrentBooker()
        with self.assertRaises(ValueError):
            b.book(5, 5)
        with self.assertRaises(ValueError):
            b.book(6, 5)
        self.assertEqual(b.count(), 0)

    def test_cancel_needs_an_exact_match(self):
        b = ConcurrentBooker()
        self.assertTrue(b.book(10, 20))
        self.assertFalse(b.cancel(10, 15))
        self.assertTrue(b.cancel(10, 20))
        self.assertFalse(b.cancel(10, 20))
        self.assertEqual(b.count(), 0)

    def test_bookings_are_sorted_by_start(self):
        b = ConcurrentBooker()
        b.book(50, 60)
        b.book(0, 10)
        b.book(20, 30)
        self.assertEqual(b.bookings(), [Interval(0, 10), Interval(20, 30), Interval(50, 60)])

    def test_book_first_free_on_empty_calendar_starts_at_from(self):
        b = ConcurrentBooker()
        self.assertEqual(b.book_first_free(7, 5), 7)
        self.assertEqual(b.bookings(), [Interval(7, 12)])

    def test_book_first_free_skips_busy_ranges_and_small_gaps(self):
        b = ConcurrentBooker()
        b.book(0, 10)
        b.book(12, 20)
        self.assertEqual(b.book_first_free(0, 5), 20)  # the gap [10,12) is too small
        self.assertEqual(b.book_first_free(0, 2), 10)  # ...but fits a booking of 2
        self.assertEqual(b.count(), 4)

    def test_book_first_free_when_from_is_inside_a_booking(self):
        b = ConcurrentBooker()
        b.book(10, 20)
        self.assertEqual(b.book_first_free(15, 5), 20)

    def test_book_first_free_rejects_non_positive_duration(self):
        b = ConcurrentBooker()
        with self.assertRaises(ValueError):
            b.book_first_free(0, 0)
        with self.assertRaises(ValueError):
            b.book_first_free(0, -1)
        self.assertEqual(b.count(), 0)

    def test_book_all_books_every_interval(self):
        b = ConcurrentBooker()
        self.assertTrue(b.book_all([Interval(30, 40), Interval(0, 10), Interval(10, 20)]))
        self.assertEqual(b.bookings(), [Interval(0, 10), Interval(10, 20), Interval(30, 40)])

    def test_book_all_is_all_or_nothing(self):
        b = ConcurrentBooker()
        self.assertTrue(b.book(10, 20))
        self.assertFalse(b.book_all([Interval(0, 5), Interval(15, 25), Interval(30, 40)]))
        self.assertEqual(b.count(), 1, "a failed book_all must book nothing")
        self.assertTrue(b.book(0, 5))
        self.assertTrue(b.book(30, 40))

    def test_book_all_rejects_overlap_inside_the_list(self):
        b = ConcurrentBooker()
        self.assertFalse(b.book_all([Interval(0, 10), Interval(5, 15)]))
        self.assertEqual(b.count(), 0)

    def test_book_all_of_empty_list_succeeds(self):
        b = ConcurrentBooker()
        self.assertTrue(b.book_all([]))
        self.assertEqual(b.count(), 0)

    def test_book_all_with_invalid_interval_throws_and_books_nothing(self):
        b = ConcurrentBooker()
        with self.assertRaises(ValueError):
            b.book_all([Interval(0, 5), Interval(7, 7)])
        self.assertEqual(b.count(), 0)

    # ---- many threads ----

    def test_concurrent_books_of_overlapping_slots_have_one_winner(self):
        # Thread t books [t, t+16). Every pair of these overlaps, so in each
        # round exactly one thread may win.
        threads = 16
        with interleaved():
            for round_ in range(100):
                b = ConcurrentBooker()
                won = [False] * threads

                def book_mine(t):
                    won[t] = b.book(t, t + 16)

                run_together(threads, book_mine)
                self.assertEqual(sum(won), 1, f"round {round_}")
                self.assertEqual(b.count(), 1, f"round {round_}")

    def test_concurrent_book_first_free_tiles_the_calendar(self):
        # 8 threads each take the first free 10-unit slot 50 times. Nobody
        # cancels, so the slots must fill [0, 4000) with no gaps and no
        # overlaps: 0, 10, 20, ...
        threads = 8
        per_thread = 50
        b = ConcurrentBooker()
        starts = [[] for _ in range(threads)]

        def take_slots(t):
            for _ in range(per_thread):
                starts[t].append(b.book_first_free(0, 10))

        with interleaved():
            run_together(threads, take_slots)

        got = sorted(s for mine in starts for s in mine)
        self.assertEqual(len(got), threads * per_thread)
        for i, start in enumerate(got):
            self.assertEqual(start, i * 10,
                             "two threads got the same slot, or a slot was skipped")
        booked = b.bookings()
        self.assertEqual(len(booked), len(got))
        self.assertTrue(no_overlaps(booked))

    def test_concurrent_book_all_has_one_winner_and_no_partial_bookings(self):
        # Every thread's group contains the shared slot [50,60), plus two
        # slots of its own. Exactly one group can win each round -- and the
        # losers must leave none of their own slots behind.
        threads = 8
        with interleaved():
            for round_ in range(100):
                b = ConcurrentBooker()
                won = [False] * threads

                def book_group(t):
                    mine = 100 + 20 * t
                    won[t] = b.book_all([Interval(mine, mine + 10),
                                         Interval(50, 60),
                                         Interval(1000 + mine, 1010 + mine)])

                run_together(threads, book_group)
                self.assertEqual(sum(won), 1, f"round {round_}")
                mine = 100 + 20 * won.index(True)
                self.assertEqual(
                    b.bookings(),
                    [Interval(50, 60), Interval(mine, mine + 10),
                     Interval(1000 + mine, 1010 + mine)],
                    f"round {round_}: a losing book_all left a partial booking")

    def test_a_failed_book_all_is_never_seen_half_done(self):
        # [50,60) is taken, so book_all([[0,10), [50,60)]) must always fail
        # and book nothing -- which leaves [0,10) free for the one thread
        # that books and cancels it over and over. A book_all that books
        # [0,10) first and takes it back once [50,60) turns out to be taken
        # is not all-or-nothing: for a moment, that thread finds [0,10)
        # taken.
        group_threads = 3
        tries = 500
        b = ConcurrentBooker()
        self.assertTrue(b.book(50, 60))
        group_won = []
        blocked = []

        def work(t):
            for _ in range(tries):
                if t < group_threads:
                    if b.book_all([Interval(0, 10), Interval(50, 60)]):
                        group_won.append(t)
                elif b.book(0, 10):
                    b.cancel(0, 10)
                else:
                    blocked.append(t)

        with interleaved():
            run_together(group_threads + 1, work)

        self.assertEqual(group_won, [], "book_all succeeded although [50,60) was taken")
        self.assertEqual(len(blocked), 0,
                         "book(0, 10) failed: a book_all that failed held [0,10) for a moment")
        self.assertEqual(b.bookings(), [Interval(50, 60)])

    def test_concurrent_mixed_operations_keep_the_calendar_consistent(self):
        # Threads book, cancel and book_first_free at random. Afterwards the
        # calendar must hold exactly the bookings that were made and not
        # cancelled -- and none may overlap.
        threads = 8
        ops = 200
        b = ConcurrentBooker()
        kept = [[] for _ in range(threads)]

        def random_ops(t):
            rng = random.Random(1234 + t)
            mine = kept[t]
            for _ in range(ops):
                op = rng.randrange(3)
                if op == 0:
                    s = rng.randrange(1000)
                    length = 1 + rng.randrange(20)
                    if b.book(s, s + length):
                        mine.append(Interval(s, s + length))
                elif op == 1:
                    length = 1 + rng.randrange(20)
                    s = b.book_first_free(rng.randrange(1000), length)
                    mine.append(Interval(s, s + length))
                elif mine:
                    k = rng.randrange(len(mine))
                    self.assertTrue(b.cancel(mine[k].start, mine[k].end),
                                    "could not cancel a booking this thread made")
                    del mine[k]

        with interleaved():
            run_together(threads, random_ops)

        expected = sorted((i for mine in kept for i in mine), key=lambda i: i.start)
        booked = b.bookings()
        self.assertTrue(sorted_by_start(booked), "bookings() is not sorted by start")
        self.assertTrue(no_overlaps(booked), "two bookings overlap")
        self.assertEqual(booked, expected)
        self.assertEqual(b.count(), len(expected))


if __name__ == "__main__":
    unittest.main()
