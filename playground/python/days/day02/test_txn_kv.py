import unittest

from .txn_kv import TxnKV


class TxnKVTest(unittest.TestCase):
    def test_set_then_get(self):
        kv = TxnKV()
        kv.set("a", "1")
        self.assertIsNotNone(kv.get("a"))
        self.assertEqual(kv.get("a"), "1")

    def test_get_missing_key_returns_nothing(self):
        kv = TxnKV()
        self.assertIsNone(kv.get("nope"))

    def test_overwrite_replaces_value(self):
        kv = TxnKV()
        kv.set("a", "1")
        kv.set("a", "2")
        self.assertEqual(kv.get("a"), "2")

    def test_remove_returns_true_only_if_key_was_visible(self):
        kv = TxnKV()
        kv.set("a", "1")
        self.assertTrue(kv.remove("a"))
        self.assertIsNone(kv.get("a"))
        self.assertFalse(kv.remove("a"))
        self.assertFalse(kv.remove("never-set"))

    def test_commit_without_open_transaction_fails(self):
        kv = TxnKV()
        self.assertFalse(kv.commit())

    def test_rollback_without_open_transaction_fails(self):
        kv = TxnKV()
        self.assertFalse(kv.rollback())

    def test_set_inside_transaction_is_visible_before_commit(self):
        kv = TxnKV()
        kv.begin()
        kv.set("a", "1")
        self.assertIsNotNone(kv.get("a"))
        self.assertEqual(kv.get("a"), "1")

    def test_rollback_discards_every_change_of_the_transaction(self):
        kv = TxnKV()
        kv.set("a", "old")
        kv.begin()
        kv.set("a", "new")
        kv.set("b", "1")
        self.assertTrue(kv.rollback())
        self.assertEqual(kv.get("a"), "old")
        self.assertIsNone(kv.get("b"))

    def test_commit_makes_changes_permanent(self):
        kv = TxnKV()
        kv.begin()
        kv.set("a", "1")
        self.assertTrue(kv.commit())
        self.assertIsNotNone(kv.get("a"))
        self.assertEqual(kv.get("a"), "1")
        self.assertFalse(kv.rollback())  # nothing left open to roll back
        self.assertEqual(kv.get("a"), "1")

    def test_remove_inside_transaction_hides_outer_value(self):
        kv = TxnKV()
        kv.set("a", "base")
        kv.begin()
        self.assertTrue(kv.remove("a"))
        self.assertIsNone(kv.get("a"))  # gone as seen from inside

    def test_rollback_restores_value_removed_inside_transaction(self):
        kv = TxnKV()
        kv.set("a", "base")
        kv.begin()
        kv.remove("a")
        kv.rollback()
        self.assertIsNotNone(kv.get("a"))
        self.assertEqual(kv.get("a"), "base")

    def test_committed_remove_deletes_from_the_store(self):
        kv = TxnKV()
        kv.set("a", "base")
        kv.begin()
        kv.remove("a")
        self.assertTrue(kv.commit())
        self.assertIsNone(kv.get("a"))

    def test_nested_inner_rollback_keeps_outer_changes(self):
        kv = TxnKV()
        kv.begin()
        kv.set("outer", "1")
        kv.begin()
        kv.set("inner", "2")
        self.assertTrue(kv.rollback())  # discards only the inner transaction
        self.assertIsNone(kv.get("inner"))
        self.assertIsNotNone(kv.get("outer"))
        self.assertEqual(kv.get("outer"), "1")

    def test_inner_commit_merges_into_parent_not_into_store(self):
        kv = TxnKV()
        kv.begin()
        kv.begin()
        kv.set("a", "1")
        self.assertTrue(kv.commit())  # inner commit: `a` now belongs to the OUTER txn
        self.assertEqual(kv.get("a"), "1")
        self.assertTrue(kv.rollback())  # outer rollback must take `a` down with it
        self.assertIsNone(kv.get("a"))

    def test_inner_committed_remove_stays_inside_parent(self):
        kv = TxnKV()
        kv.set("a", "base")
        kv.begin()
        kv.begin()
        kv.remove("a")
        self.assertTrue(kv.commit())  # remove merged into the outer txn
        self.assertIsNone(kv.get("a"))
        self.assertTrue(kv.rollback())  # outer rollback cancels the remove too
        self.assertIsNotNone(kv.get("a"))
        self.assertEqual(kv.get("a"), "base")

    def test_set_after_remove_in_same_transaction_resurrects_key(self):
        kv = TxnKV()
        kv.set("a", "base")
        kv.begin()
        kv.remove("a")
        kv.set("a", "reborn")
        self.assertEqual(kv.get("a"), "reborn")
        kv.commit()
        self.assertEqual(kv.get("a"), "reborn")

    def test_value_is_visible_through_multiple_untouched_layers(self):
        kv = TxnKV()
        kv.set("a", "base")
        kv.begin()
        kv.begin()
        kv.begin()
        self.assertIsNotNone(kv.get("a"))
        self.assertEqual(kv.get("a"), "base")


if __name__ == "__main__":
    unittest.main()
