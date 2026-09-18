import unittest

from .expiring_payments import ExpiringPaymentProcessor, Status


class ExpiringPaymentsTest(unittest.TestCase):
    def test_first_charge_is_new(self):
        p = ExpiringPaymentProcessor(10)
        r = p.charge("k", 100, 0)
        self.assertEqual(r.status, Status.NEW)
        self.assertEqual(r.amount, 100)
        self.assertEqual(r.running_total, 100)
        self.assertEqual(p.total_charged(), 100)

    def test_same_amount_within_ttl_is_replay(self):
        p = ExpiringPaymentProcessor(10)
        p.charge("k", 100, 0)
        r = p.charge("k", 100, 5)
        self.assertEqual(r.status, Status.REPLAY)
        self.assertEqual(r.amount, 100)
        self.assertEqual(r.running_total, 100)
        self.assertEqual(p.total_charged(), 100)

    def test_different_amount_within_ttl_is_conflict(self):
        p = ExpiringPaymentProcessor(10)
        p.charge("k", 100, 0)
        r = p.charge("k", 250, 5)
        self.assertEqual(r.status, Status.CONFLICT)
        self.assertEqual(r.amount, 100)  # the RECORDED amount, not 250
        self.assertEqual(r.running_total, 100)
        self.assertEqual(p.total_charged(), 100)  # nothing was charged

    def test_expired_key_is_charged_again(self):
        p = ExpiringPaymentProcessor(10)
        p.charge("k", 100, 0)
        r = p.charge("k", 100, 20)
        self.assertEqual(r.status, Status.NEW)
        self.assertEqual(p.total_charged(), 200)

    def test_expiry_boundary_is_inclusive(self):
        # Recorded at 0 with ttl 10 -> expired exactly at now == 10.
        p = ExpiringPaymentProcessor(10)
        p.charge("k", 100, 0)
        self.assertEqual(p.charge("k", 100, 9).status, Status.REPLAY)
        self.assertEqual(p.charge("k", 100, 10).status, Status.NEW)
        self.assertEqual(p.total_charged(), 200)

    def test_replay_does_not_extend_expiry(self):
        p = ExpiringPaymentProcessor(10)
        p.charge("k", 100, 0)
        self.assertEqual(p.charge("k", 100, 9).status, Status.REPLAY)
        # If the replay at 9 had refreshed the record, the key would still be
        # active at 10. It must not be.
        self.assertEqual(p.charge("k", 100, 10).status, Status.NEW)

    def test_conflict_changes_nothing(self):
        p = ExpiringPaymentProcessor(10)
        p.charge("k", 100, 0)
        p.charge("k", 250, 5)  # conflict
        r = p.charge("k", 100, 9)  # original amount again
        self.assertEqual(r.status, Status.REPLAY)  # record still 100
        self.assertEqual(r.amount, 100)
        self.assertEqual(p.total_charged(), 100)
        self.assertEqual(p.charge("k", 100, 10).status, Status.NEW)  # expiry not shifted

    def test_after_expiry_different_amount_is_new_not_conflict(self):
        p = ExpiringPaymentProcessor(10)
        p.charge("k", 100, 0)
        r = p.charge("k", 250, 15)
        self.assertEqual(r.status, Status.NEW)
        self.assertEqual(r.amount, 250)
        self.assertEqual(p.total_charged(), 350)
        # The record was replaced: 250 is now the recorded amount.
        self.assertEqual(p.charge("k", 250, 16).status, Status.REPLAY)
        self.assertEqual(p.charge("k", 100, 16).status, Status.CONFLICT)

    def test_keys_are_independent(self):
        p = ExpiringPaymentProcessor(10)
        p.charge("a", 100, 0)
        p.charge("b", 50, 8)
        self.assertEqual(p.total_charged(), 150)
        self.assertEqual(p.charge("a", 100, 12).status, Status.NEW)  # a expired
        self.assertEqual(p.charge("b", 50, 12).status, Status.REPLAY)  # b still active

    def test_active_keys_counts_only_unexpired(self):
        p = ExpiringPaymentProcessor(10)
        p.charge("a", 1, 0)  # expires at 10
        p.charge("b", 2, 5)  # expires at 15
        p.charge("c", 3, 8)  # expires at 18
        self.assertEqual(p.active_keys(0), 3)
        self.assertEqual(p.active_keys(12), 2)
        self.assertEqual(p.active_keys(15), 1)
        self.assertEqual(p.active_keys(18), 0)

    def test_zero_ttl_means_every_charge_is_new(self):
        p = ExpiringPaymentProcessor(0)
        self.assertEqual(p.charge("k", 100, 5).status, Status.NEW)
        self.assertEqual(p.charge("k", 100, 5).status, Status.NEW)
        self.assertEqual(p.total_charged(), 200)
        self.assertEqual(p.active_keys(5), 0)

    def test_replay_reports_current_total_not_total_at_first_charge(self):
        p = ExpiringPaymentProcessor(100)
        p.charge("a", 100, 0)
        p.charge("b", 50, 1)
        r = p.charge("a", 100, 2)
        self.assertEqual(r.status, Status.REPLAY)
        self.assertEqual(r.running_total, 150)

    def test_total_is_sum_of_new_charges_only(self):
        p = ExpiringPaymentProcessor(10)
        p.charge("a", 100, 0)  # new
        p.charge("a", 100, 1)  # replay
        p.charge("a", 999, 2)  # conflict
        p.charge("b", 50, 3)  # new
        p.charge("a", 25, 12)  # expired -> new
        self.assertEqual(p.total_charged(), 175)


if __name__ == "__main__":
    unittest.main()
