import time
import unittest

from ._concurrency import RUN_TIMEOUT, TimeLimitedTestCase, interleaved, run_together
from .account import Account

THREADS = 8
PER_THREAD = 1000


class AccountTest(TimeLimitedTestCase):
    def test_starts_with_initial_balance(self):
        a = Account(100)
        self.assertEqual(a.balance(), 100)

    def test_negative_initial_balance_throws(self):
        with self.assertRaises(ValueError):
            Account(-1)

    def test_deposit_adds_and_withdraw_subtracts(self):
        a = Account(0)
        a.deposit(50)
        self.assertTrue(a.withdraw(20))
        self.assertEqual(a.balance(), 30)

    def test_withdraw_exact_balance_succeeds(self):
        a = Account(40)
        self.assertTrue(a.withdraw(40))
        self.assertEqual(a.balance(), 0)

    def test_withdraw_more_than_balance_fails_and_changes_nothing(self):
        a = Account(30)
        self.assertFalse(a.withdraw(31))
        self.assertEqual(a.balance(), 30)

    def test_non_positive_amounts_throw_and_change_nothing(self):
        a = Account(10)
        with self.assertRaises(ValueError):
            a.deposit(0)
        with self.assertRaises(ValueError):
            a.deposit(-5)
        with self.assertRaises(ValueError):
            a.withdraw(0)
        with self.assertRaises(ValueError):
            a.withdraw(-5)
        self.assertEqual(a.balance(), 10)

    def test_concurrent_withdrawals_never_overdraw(self):
        # 8 threads try to take 800 in total from an account holding 200.
        # Exactly 200 withdrawals of 1 can succeed. The race only matters at
        # the moment the money runs out, so play it 20 times.
        with interleaved():
            for round_ in range(20):
                a = Account(200)
                succeeded = [0] * THREADS

                def withdraw_ones(t):
                    for _ in range(100):
                        if a.withdraw(1):
                            succeeded[t] += 1

                run_together(THREADS, withdraw_ones)
                self.assertEqual(sum(succeeded), 200, f"round {round_}")
                self.assertEqual(a.balance(), 0, f"round {round_}")

    def test_concurrent_deposits_are_never_lost(self):
        a = Account(0)

        def deposit_ones(t):
            for _ in range(PER_THREAD):
                a.deposit(1)

        with interleaved():
            run_together(THREADS, deposit_ones)

        self.assertEqual(a.balance(), THREADS * PER_THREAD)

    def test_concurrent_deposits_and_withdrawals_balance_out(self):
        # 4 threads deposit 3 at a time, 4 threads withdraw 2 at a time, and
        # one watcher checks that the balance is never seen below zero.
        a = Account(0)
        withdrawn = [0] * THREADS
        finished = [False] * THREADS
        saw_negative = False
        give_up = time.monotonic() + RUN_TIMEOUT  # in case the workers hang

        def work_or_watch(t):
            nonlocal saw_negative
            if t == THREADS:  # the watcher
                while not all(finished) and time.monotonic() < give_up:
                    if a.balance() < 0:
                        saw_negative = True
                return
            for _ in range(PER_THREAD):
                if t < 4:
                    a.deposit(3)
                elif a.withdraw(2):
                    withdrawn[t] += 2
            finished[t] = True

        with interleaved():
            run_together(THREADS + 1, work_or_watch)

        self.assertFalse(saw_negative, "the balance was seen below zero")
        self.assertEqual(a.balance(), 4 * PER_THREAD * 3 - sum(withdrawn))


if __name__ == "__main__":
    unittest.main()
