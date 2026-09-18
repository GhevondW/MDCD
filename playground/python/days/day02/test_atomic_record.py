import unittest

from .atomic_record import RecordStore


class AtomicRecordTest(unittest.TestCase):
    def test_get_on_unknown_key_is_none(self):
        store = RecordStore()
        self.assertIsNone(store.get("ghost"))

    def test_valid_set_commits_both_fields(self):
        store = RecordStore()
        self.assertTrue(store.set("acct1", 10, 5))
        r = store.get("acct1")
        self.assertIsNotNone(r)
        self.assertEqual(r.a, 10)
        self.assertEqual(r.b, 5)

    def test_invalid_set_is_rejected(self):
        store = RecordStore()
        self.assertFalse(store.set("acct2", -5, -5))
        self.assertIsNone(store.get("acct2"))

    def test_rejected_set_leaves_existing_record_untouched(self):
        # The core atomicity check: a rejected update on an EXISTING key
        # must not partially apply -- field a must not change even though
        # it would have been "valid" applied in isolation.
        store = RecordStore()
        store.set("acct1", 10, 5)
        self.assertFalse(store.set("acct1", -20, 3))  # -20 + 3 < 0 -> reject
        r = store.get("acct1")
        self.assertIsNotNone(r)
        self.assertEqual(r.a, 10)  # unchanged, not -20
        self.assertEqual(r.b, 5)  # unchanged

    def test_zero_sum_is_valid(self):
        store = RecordStore()
        self.assertTrue(store.set("a", 0, 0))
        r = store.get("a")
        self.assertIsNotNone(r)
        self.assertEqual(r.a, 0)
        self.assertEqual(r.b, 0)

    def test_negative_field_with_compensating_positive_is_valid(self):
        # The rule is a + b >= 0, NOT "both fields individually
        # non-negative" -- a single negative field with a large enough
        # positive one is fine.
        store = RecordStore()
        self.assertTrue(store.set("acct3", -5, 10))
        r = store.get("acct3")
        self.assertIsNotNone(r)
        self.assertEqual(r.a, -5)
        self.assertEqual(r.b, 10)

    def test_valid_set_overwrites_previous_valid_record(self):
        store = RecordStore()
        store.set("acct1", 10, 5)
        self.assertTrue(store.set("acct1", 1, 2))  # second valid set replaces
        r = store.get("acct1")
        self.assertEqual(r.a, 1)
        self.assertEqual(r.b, 2)

    def test_keys_are_independent(self):
        store = RecordStore()
        store.set("acct1", 10, 5)
        store.set("acct2", 1, 1)
        self.assertFalse(store.set("acct2", -100, 0))  # reject on acct2 only
        r1 = store.get("acct1")
        r2 = store.get("acct2")
        self.assertEqual(r1.a, 10)
        self.assertEqual(r1.b, 5)
        self.assertEqual(r2.a, 1)  # acct2's rejection didn't touch acct1
        self.assertEqual(r2.b, 1)


if __name__ == "__main__":
    unittest.main()
