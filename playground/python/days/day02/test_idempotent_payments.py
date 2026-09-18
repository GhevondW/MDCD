import unittest

from .idempotent_payments import PaymentProcessor


class IdempotentPaymentsTest(unittest.TestCase):
    def test_first_charge_is_new(self):
        p = PaymentProcessor()
        r = p.charge("orderA", 100)
        self.assertTrue(r.was_new)
        self.assertEqual(r.amount, 100)
        self.assertEqual(r.running_total, 100)
        self.assertEqual(p.total_charged(), 100)

    def test_repeat_key_is_not_charged(self):
        p = PaymentProcessor()
        p.charge("orderA", 100)
        r = p.charge("orderA", 100)
        self.assertFalse(r.was_new)
        self.assertEqual(r.amount, 100)
        self.assertEqual(r.running_total, 100)  # unchanged -- not double-charged
        self.assertEqual(p.total_charged(), 100)

    def test_different_keys_accumulate(self):
        p = PaymentProcessor()
        p.charge("orderA", 100)
        p.charge("orderB", 250)
        self.assertEqual(p.total_charged(), 350)

    def test_repeat_ignores_new_amount_argument(self):
        # A caller retrying a request might resend the same amount; a repeat
        # must still report the ORIGINAL amount, never a re-derived one.
        p = PaymentProcessor()
        p.charge("x", 5)
        r = p.charge("x", 5)
        self.assertEqual(r.amount, 5)
        self.assertEqual(p.total_charged(), 5)

    def test_many_repeats_still_charge_once(self):
        p = PaymentProcessor()
        p.charge("x", 5)
        p.charge("x", 5)
        p.charge("x", 5)
        p.charge("y", 10)
        self.assertEqual(p.total_charged(), 15)

    def test_repeat_is_idempotent_even_with_other_charges_between(self):
        p = PaymentProcessor()
        p.charge("A", 100)
        p.charge("B", 50)
        r = p.charge("A", 100)  # repeat of A, not adjacent to the first
        self.assertFalse(r.was_new)
        self.assertEqual(r.amount, 100)
        self.assertEqual(r.running_total, 150)  # reflects total so far
        p.charge("C", 25)
        self.assertEqual(p.total_charged(), 175)  # A never double-counted

    def test_zero_amount_charge_is_still_tracked_by_key(self):
        p = PaymentProcessor()
        first = p.charge("free-trial", 0)
        self.assertTrue(first.was_new)
        self.assertEqual(first.amount, 0)
        repeat = p.charge("free-trial", 0)
        self.assertFalse(repeat.was_new)
        self.assertEqual(p.total_charged(), 0)


if __name__ == "__main__":
    unittest.main()
