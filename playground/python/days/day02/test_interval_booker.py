import unittest

from .interval_booker import IntervalBooker


class IntervalBookerTest(unittest.TestCase):
    def test_book_on_empty_calendar_succeeds(self):
        b = IntervalBooker()
        self.assertTrue(b.book(10, 20))
        self.assertEqual(b.count(), 1)

    def test_adjacent_intervals_do_not_overlap(self):
        # Half-open semantics: [0,10) and [10,20) share only the boundary point.
        b = IntervalBooker()
        self.assertTrue(b.book(0, 10))
        self.assertTrue(b.book(10, 20))
        self.assertEqual(b.count(), 2)

    def test_partial_overlap_is_rejected_both_directions(self):
        b = IntervalBooker()
        self.assertTrue(b.book(10, 20))
        self.assertFalse(b.book(15, 25))  # overlaps the tail
        self.assertFalse(b.book(5, 15))  # overlaps the head
        self.assertEqual(b.count(), 1)

    def test_containment_is_rejected_both_directions(self):
        b = IntervalBooker()
        self.assertTrue(b.book(10, 20))
        self.assertFalse(b.book(12, 18))  # inside the existing booking
        self.assertFalse(b.book(5, 25))  # swallows the existing booking

    def test_identical_interval_is_rejected(self):
        b = IntervalBooker()
        self.assertTrue(b.book(10, 20))
        self.assertFalse(b.book(10, 20))

    def test_rejected_book_changes_nothing(self):
        b = IntervalBooker()
        self.assertTrue(b.book(10, 20))
        self.assertFalse(b.book(15, 25))
        self.assertEqual(b.count(), 1)
        self.assertTrue(b.book(20, 30))  # the slot the rejected book didn't take

    def test_empty_or_reversed_interval_throws(self):
        b = IntervalBooker()
        with self.assertRaises(ValueError):
            b.book(10, 10)
        with self.assertRaises(ValueError):
            b.book(20, 10)
        self.assertEqual(b.count(), 0)

    def test_cancel_exact_match_frees_the_slot(self):
        b = IntervalBooker()
        self.assertTrue(b.book(10, 20))
        self.assertTrue(b.cancel(10, 20))
        self.assertEqual(b.count(), 0)
        self.assertTrue(b.book(15, 25))  # the freed range is bookable again

    def test_cancel_requires_exact_match(self):
        b = IntervalBooker()
        self.assertTrue(b.book(10, 20))
        self.assertFalse(b.cancel(10, 15))  # sub-range: no
        self.assertFalse(b.cancel(5, 25))  # super-range: no
        self.assertEqual(b.count(), 1)

    def test_cancel_unknown_interval_returns_false(self):
        b = IntervalBooker()
        self.assertFalse(b.cancel(10, 20))

    def test_first_free_on_empty_calendar_is_from(self):
        b = IntervalBooker()
        self.assertEqual(b.first_free(7, 5), 7)

    def test_first_free_skips_busy_ranges_and_too_small_gaps(self):
        b = IntervalBooker()
        b.book(10, 20)
        b.book(20, 30)
        # [5,15) collides; the only room left starts after the last booking.
        self.assertEqual(b.first_free(5, 10), 30)

    def test_first_free_finds_gap_of_exact_size(self):
        b = IntervalBooker()
        b.book(0, 10)
        b.book(15, 25)
        self.assertEqual(b.first_free(0, 5), 10)  # the gap [10,15) fits exactly

    def test_first_free_when_from_falls_inside_a_booking(self):
        b = IntervalBooker()
        b.book(10, 20)
        self.assertEqual(b.first_free(12, 5), 20)

    def test_first_free_ignores_bookings_entirely_before_from(self):
        b = IntervalBooker()
        b.book(0, 10)
        self.assertEqual(b.first_free(50, 5), 50)

    def test_first_free_result_is_actually_bookable(self):
        b = IntervalBooker()
        b.book(3, 9)
        b.book(12, 40)
        b.book(41, 50)
        s = b.first_free(0, 4)
        self.assertTrue(b.book(s, s + 4))

    def test_non_positive_duration_throws(self):
        b = IntervalBooker()
        with self.assertRaises(ValueError):
            b.first_free(0, 0)
        with self.assertRaises(ValueError):
            b.first_free(0, -3)


if __name__ == "__main__":
    unittest.main()
