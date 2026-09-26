import random
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


    # ---- transfer and total: two accounts locked at once ----

    def test_transfer_moves_money(self):
        a, b = Account(100), Account(50)
        self.assertTrue(a.transfer(b, 30))
        self.assertEqual(a.balance(), 70)
        self.assertEqual(b.balance(), 80)

    def test_transfer_without_enough_money_changes_nothing(self):
        a, b = Account(10), Account(0)
        self.assertFalse(a.transfer(b, 11))
        self.assertEqual(a.balance(), 10)
        self.assertEqual(b.balance(), 0)

    def test_transfer_rejects_bad_arguments(self):
        a, b = Account(10), Account(0)
        for amount in (0, -5):
            with self.assertRaises(ValueError):
                a.transfer(b, amount)
        with self.assertRaises(ValueError, msg="a transfer to the same account"):
            a.transfer(a, 1)
        self.assertEqual(a.balance(), 10)
        self.assertEqual(b.balance(), 0)

    def test_total_adds_both_balances(self):
        a, b = Account(100), Account(50)
        self.assertEqual(Account.total(a, b), 150)
        self.assertEqual(Account.total(b, a), 150)
        with self.assertRaises(ValueError):
            Account.total(a, a)

    def test_opposite_transfers_do_not_deadlock(self):
        # Half the threads move money from x to y, the other half from y to
        # x -- the slides' "Two transfers, two locks".
        x, y = Account(1000), Account(1000)

        def move(t):
            for _ in range(300):
                if t % 2 == 0:
                    x.transfer(y, 1)
                else:
                    y.transfer(x, 1)

        with interleaved():
            run_together(4, move, msg="opposite transfers never finished")

        self.assertEqual(x.balance() + y.balance(), 2000,
                         "money appeared or disappeared")
        self.assertGreaterEqual(x.balance(), 0)
        self.assertGreaterEqual(y.balance(), 0)

    def test_total_never_sees_a_half_done_transfer(self):
        # Transfers move money back and forth between x and y while a watcher
        # keeps asking for the total. A transfer done as two steps (take from
        # x, then give to y) shows the watcher money in flight.
        movers = 4
        x, y = Account(1000), Account(1000)
        finished = [False] * movers
        wrong_totals = []
        give_up = time.monotonic() + RUN_TIMEOUT  # in case the movers hang

        def move_or_watch(t):
            if t == movers:  # the watcher
                while not all(finished) and time.monotonic() < give_up:
                    seen = Account.total(x, y)
                    if seen != 2000:
                        wrong_totals.append(seen)
                return
            for _ in range(300):
                if t % 2 == 0:
                    x.transfer(y, 7)
                else:
                    y.transfer(x, 7)
            finished[t] = True

        with interleaved():
            run_together(movers + 1, move_or_watch)

        self.assertEqual(wrong_totals, [], "total() saw money that was in flight")
        self.assertEqual(Account.total(x, y), 2000)

    def test_concurrent_transfers_among_many_accounts_keep_the_money(self):
        # 8 threads move random amounts between random pairs of 6 accounts,
        # in both directions. No deadlock, no money made or lost, no account
        # below zero.
        accounts = [Account(1000) for _ in range(6)]

        def move_randomly(t):
            rng = random.Random(99 + t)
            for _ in range(300):
                src = rng.randrange(len(accounts))
                dst = rng.randrange(len(accounts) - 1)
                if dst >= src:
                    dst += 1  # any account but src
                accounts[src].transfer(accounts[dst], 1 + rng.randrange(50))

        with interleaved():
            run_together(THREADS, move_randomly)

        for a in accounts:
            self.assertGreaterEqual(a.balance(), 0, "an account went below zero")
        self.assertEqual(sum(a.balance() for a in accounts), 1000 * len(accounts),
                         "money appeared or disappeared")


if __name__ == "__main__":
    unittest.main()
