import threading
import time
import unittest

from ._concurrency import TimeLimitedTestCase, interleaved, run_together
from .bounded_stack import BoundedStack

THREADS = 8

# The single-threaded contract is the same as Day 2 -- a quick check that
# it still holds.


class BoundedStackTest(TimeLimitedTestCase):
    def test_push_pop_is_lifo(self):
        s = BoundedStack(3)
        s.push(1)
        s.push(2)
        s.push(3)
        self.assertTrue(s.full())
        self.assertEqual(s.top(), 3)
        self.assertEqual(s.pop(), 3)
        self.assertEqual(s.pop(), 2)
        self.assertEqual(s.size(), 1)

    def test_full_and_empty_still_throw(self):
        s = BoundedStack(1)
        self.assertTrue(s.empty())
        with self.assertRaises(IndexError):
            s.pop()
        with self.assertRaises(IndexError):
            s.top()
        s.push(7)
        with self.assertRaises(IndexError):
            s.push(8)
        self.assertEqual(s.size(), 1)
        self.assertEqual(s.top(), 7)

    def test_concurrent_pushes_stop_exactly_at_capacity(self):
        # 8 threads push 400 values into a stack that holds 200. Exactly
        # 200 pushes may succeed, and the stack must then hold exactly the
        # values whose push succeeded. The race only matters at the moment
        # the stack fills up, so play it 20 times.
        capacity = 200
        per_thread = 50
        with interleaved():
            for round_ in range(20):
                s = BoundedStack(capacity)
                pushed = [[] for _ in range(THREADS)]

                def push_values(t):
                    for i in range(per_thread):
                        v = t * 1000 + i
                        try:
                            s.push(v)
                            pushed[t].append(v)
                        except IndexError:
                            pass

                run_together(THREADS, push_values)

                expected = sorted(v for mine in pushed for v in mine)
                self.assertEqual(len(expected), capacity,
                                 f"round {round_}: pushes that succeeded")
                self.assertEqual(s.size(), capacity, f"round {round_}")
                self.assertTrue(s.full(), f"round {round_}")

                held = []
                for _ in range(capacity):
                    if s.empty():
                        break
                    held.append(s.pop())
                self.assertTrue(s.empty(), f"round {round_}")
                self.assertEqual(sorted(held), expected,
                                 f"round {round_}: the stack lost a value or made one up")

    def test_concurrent_pops_return_each_value_once(self):
        values = 4000
        s = BoundedStack(values)
        for v in range(values):
            s.push(v)

        # Pop until the stack says it is empty. (Stop early if more values
        # come out than ever went in -- the stack is making them up.)
        popped = [[] for _ in range(THREADS)]

        def pop_values(t):
            while sum(map(len, popped)) <= values:
                try:
                    popped[t].append(s.pop())
                except IndexError:
                    return

        with interleaved():
            run_together(THREADS, pop_values)

        got = sorted(v for mine in popped for v in mine)
        self.assertEqual(len(got), values, "a value was popped twice, or lost")
        for expected, v in enumerate(got):
            self.assertEqual(v, expected, "a value was popped twice, or lost")
        self.assertTrue(s.empty())

    def test_concurrent_push_and_pop_conserve_values(self):
        # A small stack, 4 pushers and 4 poppers. Pushers retry while it is
        # full, poppers retry while it is empty. Every value pushed must come
        # out exactly once.
        pushers = 4
        per_pusher = 1000
        total = pushers * per_pusher
        give_up_after = 5  # seconds
        s = BoundedStack(16)
        popped = [[] for _ in range(pushers)]
        deadline = time.monotonic() + give_up_after
        gave_up = threading.Event()

        def past_deadline():
            if time.monotonic() > deadline:
                gave_up.set()
            return gave_up.is_set()

        def push_or_pop(t):
            if t < pushers:
                for i in range(per_pusher):
                    while True:
                        try:
                            s.push(t * per_pusher + i)
                            break
                        except IndexError:
                            if past_deadline():
                                return
                            time.sleep(0)  # let another thread run
            else:
                mine = popped[t - pushers]
                while sum(map(len, popped)) < total and not past_deadline():
                    try:
                        mine.append(s.pop())
                    except IndexError:
                        time.sleep(0)  # let another thread run

        with interleaved():
            run_together(2 * pushers, push_or_pop)

        self.assertFalse(gave_up.is_set(),
                         f"gave up after {give_up_after} s -- values were lost")
        got = sorted(v for mine in popped for v in mine)
        self.assertEqual(len(got), total)
        for expected, v in enumerate(got):
            self.assertEqual(v, expected, "a value was popped twice, or lost")
        self.assertTrue(s.empty())


    # ---- try_push / try_pop: check and act as one step ----

    def test_try_push_and_try_pop_follow_the_rules(self):
        s = BoundedStack(2)
        self.assertIsNone(s.try_pop(), "empty stack")
        self.assertTrue(s.try_push(1))
        self.assertTrue(s.try_push(2))
        self.assertFalse(s.try_push(3), "full stack")
        self.assertEqual(s.size(), 2)
        self.assertEqual(s.try_pop(), 2)
        self.assertEqual(s.try_pop(), 1)
        self.assertIsNone(s.try_pop())
        self.assertTrue(s.empty())

    def test_concurrent_try_pop_takes_each_value_once(self):
        # 8 threads empty the stack with try_pop -- the one-call version of
        # "if not s.empty(): v = s.top(); s.pop()".
        values = 4000
        s = BoundedStack(values)
        for v in range(values):
            s.push(v)
        popped = [[] for _ in range(THREADS)]

        def pop_values(t):
            while sum(map(len, popped)) <= values:  # more than values: made up
                v = s.try_pop()
                if v is None:
                    return
                popped[t].append(v)

        with interleaved():
            run_together(THREADS, pop_values)

        got = sorted(v for mine in popped for v in mine)
        self.assertEqual(len(got), values, "a value was popped twice, or lost")
        for expected, v in enumerate(got):
            self.assertEqual(v, expected, "a value was popped twice, or lost")
        self.assertTrue(s.empty())

    def test_concurrent_try_push_stops_exactly_at_capacity(self):
        # Like test_concurrent_pushes_stop_exactly_at_capacity, with try_push.
        capacity = 200
        per_thread = 50
        with interleaved():
            for round_ in range(20):
                s = BoundedStack(capacity)
                pushed = [[] for _ in range(THREADS)]

                def push_values(t):
                    for i in range(per_thread):
                        v = t * 1000 + i
                        if s.try_push(v):
                            pushed[t].append(v)

                run_together(THREADS, push_values)

                expected = sorted(v for mine in pushed for v in mine)
                self.assertEqual(len(expected), capacity,
                                 f"round {round_}: try_push calls that returned True")
                self.assertEqual(s.size(), capacity, f"round {round_}")
                held = []
                for _ in range(capacity):
                    v = s.try_pop()
                    if v is None:
                        break
                    held.append(v)
                self.assertEqual(sorted(held), expected,
                                 f"round {round_}: the stack lost a value or made one up")

    def test_concurrent_try_push_and_try_pop_conserve_values(self):
        # A small stack, 4 pushers and 4 poppers, all using the try_ calls.
        pushers = 4
        per_pusher = 1000
        total = pushers * per_pusher
        give_up_after = 5  # seconds
        s = BoundedStack(16)
        popped = [[] for _ in range(pushers)]
        deadline = time.monotonic() + give_up_after
        gave_up = threading.Event()

        def past_deadline():
            if time.monotonic() > deadline:
                gave_up.set()
            return gave_up.is_set()

        def push_or_pop(t):
            if t < pushers:
                for i in range(per_pusher):
                    while not s.try_push(t * per_pusher + i):
                        if past_deadline():
                            return
                        time.sleep(0)  # let another thread run
            else:
                mine = popped[t - pushers]
                while sum(map(len, popped)) < total and not past_deadline():
                    v = s.try_pop()
                    if v is None:
                        time.sleep(0)  # let another thread run
                    else:
                        mine.append(v)

        with interleaved():
            run_together(2 * pushers, push_or_pop)

        self.assertFalse(gave_up.is_set(),
                         f"gave up after {give_up_after} s -- values were lost")
        got = sorted(v for mine in popped for v in mine)
        self.assertEqual(len(got), total)
        for expected, v in enumerate(got):
            self.assertEqual(v, expected, "a value was popped twice, or lost")


if __name__ == "__main__":
    unittest.main()
