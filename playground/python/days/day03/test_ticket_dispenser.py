import unittest

from ._concurrency import TimeLimitedTestCase, interleaved, run_together
from .ticket_dispenser import TicketDispenser

THREADS = 8
PER_THREAD = 1000


class TicketDispenserTest(TimeLimitedTestCase):
    def test_first_ticket_is_one(self):
        d = TicketDispenser()
        self.assertEqual(d.next(), 1)

    def test_tickets_count_up_by_one(self):
        d = TicketDispenser()
        self.assertEqual(d.next(), 1)
        self.assertEqual(d.next(), 2)
        self.assertEqual(d.next(), 3)

    def test_issued_counts_handed_out_tickets(self):
        d = TicketDispenser()
        self.assertEqual(d.issued(), 0)
        d.next()
        d.next()
        self.assertEqual(d.issued(), 2)

    def test_concurrent_tickets_are_unique_and_gapless(self):
        d = TicketDispenser()
        got = [[] for _ in range(THREADS)]

        def take_tickets(t):
            for _ in range(PER_THREAD):
                got[t].append(d.next())

        with interleaved():
            run_together(THREADS, take_tickets)

        tickets = sorted(ticket for mine in got for ticket in mine)
        for i, ticket in enumerate(tickets):
            self.assertEqual(ticket, i + 1,
                             "a ticket was handed out twice, or a number was skipped")
        self.assertEqual(d.issued(), THREADS * PER_THREAD)

    def test_each_thread_sees_its_tickets_increase(self):
        d = TicketDispenser()
        got = [[] for _ in range(THREADS)]

        def take_tickets(t):
            for _ in range(PER_THREAD):
                got[t].append(d.next())

        with interleaved():
            run_together(THREADS, take_tickets)

        for mine in got:
            for earlier, later in zip(mine, mine[1:]):
                self.assertLess(earlier, later, "a later ticket had a smaller number")

    def test_issued_never_goes_backwards_while_tickets_are_taken(self):
        takers = 4
        d = TicketDispenser()
        went_backwards = False

        def take_or_watch(t):
            nonlocal went_backwards
            if t == takers:  # the watcher
                last = 0
                for _ in range(PER_THREAD):
                    now = d.issued()
                    if now < last:
                        went_backwards = True
                    last = now
                return
            for _ in range(PER_THREAD):
                d.next()

        with interleaved():
            run_together(takers + 1, take_or_watch)

        self.assertFalse(went_backwards, "issued() went down while tickets were taken")
        self.assertEqual(d.issued(), takers * PER_THREAD)


if __name__ == "__main__":
    unittest.main()
