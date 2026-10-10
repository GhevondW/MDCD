import threading
import time
import unittest

from ._concurrency import (AtomicCounter, Background, TimeLimitedTestCase,
                           interleaved, run_together)
from .bounded_blocking_queue import BoundedBlockingQueue

WAIT = 3.0  # a call that should return "soon" gets this long


class SingleThreadedTest(TimeLimitedTestCase):
    def test_items_come_out_in_fifo_order(self):
        q = BoundedBlockingQueue(3)
        for v in (1, 2, 3):
            self.assertTrue(q.push(v))
        self.assertEqual(q.size(), 3)
        self.assertEqual(q.pop(), (True, 1))
        self.assertEqual(q.pop(), (True, 2))
        self.assertTrue(q.push(4))
        self.assertEqual(q.pop(), (True, 3))
        self.assertEqual(q.pop(), (True, 4))
        self.assertEqual(q.size(), 0)

    def test_capacity_must_be_at_least_one(self):
        with self.assertRaises(ValueError):
            BoundedBlockingQueue(0)
        with self.assertRaises(ValueError):
            BoundedBlockingQueue(-5)
        self.assertEqual(BoundedBlockingQueue(7).capacity(), 7)

    def test_none_is_an_ordinary_item(self):
        q = BoundedBlockingQueue(2)
        self.assertTrue(q.push(None))
        self.assertEqual(q.try_pop(), (True, None))
        self.assertEqual(q.try_pop(), (False, None))

    def test_try_push_and_try_pop_never_wait(self):
        q = BoundedBlockingQueue(2)
        self.assertEqual(q.try_pop(), (False, None))
        self.assertTrue(q.try_push(1))
        self.assertTrue(q.try_push(2))
        self.assertFalse(q.try_push(3))  # full: nothing changes
        self.assertEqual(q.size(), 2)
        self.assertEqual(q.try_pop(), (True, 1))
        self.assertEqual(q.try_pop(), (True, 2))
        self.assertEqual(q.try_pop(), (False, None))

    def test_close_hands_out_what_is_left_then_nothing(self):
        q = BoundedBlockingQueue(5)
        for v in (1, 2, 3):
            q.push(v)
        self.assertFalse(q.closed())
        q.close()
        self.assertTrue(q.closed())
        self.assertFalse(q.push(9))
        self.assertFalse(q.try_push(9))
        self.assertFalse(q.push_for(9, 0.01))
        self.assertEqual(q.size(), 3)
        self.assertEqual(q.pop(), (True, 1))
        self.assertEqual(q.try_pop(), (True, 2))
        self.assertEqual(q.pop_for(1.0), (True, 3))
        self.assertEqual(q.pop(), (False, None))
        self.assertEqual(q.try_pop(), (False, None))
        self.assertEqual(q.pop_for(1.0), (False, None))

    def test_close_twice_is_fine(self):
        q = BoundedBlockingQueue(1)
        q.close()
        q.close()
        self.assertTrue(q.closed())


class BlockingTest(TimeLimitedTestCase):
    def test_push_waits_while_full_and_goes_on_after_a_pop(self):
        q = BoundedBlockingQueue(2)
        q.push(1)
        q.push(2)
        b = Background(lambda: q.push(3))
        self.assertFalse(b.wait(0.3), "push on a full queue returned at once")
        self.assertEqual(q.size(), 2)
        self.assertEqual(q.pop(), (True, 1))
        self.assertTrue(b.wait(WAIT), "push did not wake up after a pop")
        self.assertTrue(b.result())
        self.assertEqual(q.pop(), (True, 2))
        self.assertEqual(q.pop(), (True, 3))

    def test_pop_waits_while_empty_and_goes_on_after_a_push(self):
        q = BoundedBlockingQueue(2)
        b = Background(q.pop)
        self.assertFalse(b.wait(0.3), "pop on an empty queue returned at once")
        q.push("x")
        self.assertTrue(b.wait(WAIT), "pop did not wake up after a push")
        self.assertEqual(b.result(), (True, "x"))

    def test_close_wakes_every_waiting_consumer(self):
        q = BoundedBlockingQueue(2)
        waiting = [Background(q.pop) for _ in range(4)]
        time.sleep(0.2)
        q.close()
        for b in waiting:
            self.assertTrue(b.wait(WAIT), "a consumer slept through close()")
            self.assertEqual(b.result(), (False, None))

    def test_close_wakes_every_waiting_producer_and_adds_nothing(self):
        q = BoundedBlockingQueue(2)
        q.push("a")
        q.push("b")
        waiting = [Background(lambda i=i: q.push(i)) for i in range(4)]
        time.sleep(0.2)
        q.close()
        for b in waiting:
            self.assertTrue(b.wait(WAIT), "a producer slept through close()")
            self.assertFalse(b.result())
        self.assertEqual(q.size(), 2)
        self.assertEqual(q.pop(), (True, "a"))
        self.assertEqual(q.pop(), (True, "b"))
        self.assertEqual(q.pop(), (False, None))

    def test_each_item_wakes_exactly_one_consumer_of_many(self):
        # 4 consumers sleep; 4 items arrive one by one: every consumer gets
        # one. (A push that wakes nobody leaves a consumer asleep forever.)
        q = BoundedBlockingQueue(8)
        waiting = [Background(q.pop) for _ in range(4)]
        time.sleep(0.2)
        for v in range(4):
            q.push(v)
        got = []
        for b in waiting:
            self.assertTrue(b.wait(WAIT), "a consumer was never woken up")
            ok, item = b.result()
            self.assertTrue(ok)
            got.append(item)
        self.assertEqual(sorted(got), [0, 1, 2, 3])

    def test_each_pop_wakes_a_producer_of_many(self):
        q = BoundedBlockingQueue(1)
        q.push(-1)
        waiting = [Background(lambda i=i: q.push(i)) for i in range(4)]
        time.sleep(0.2)
        got = []
        for _ in range(5):
            b = Background(q.pop)
            self.assertTrue(b.wait(WAIT), "a producer was never woken up")
            ok, item = b.result()
            self.assertTrue(ok)
            got.append(item)
            time.sleep(0.02)
        for b in waiting:
            self.assertTrue(b.wait(WAIT), "a producer was never woken up")
            self.assertTrue(b.result())
        self.assertEqual(sorted(got), [-1, 0, 1, 2, 3])


class TimedTest(TimeLimitedTestCase):
    def test_pop_for_times_out_on_an_empty_queue(self):
        q = BoundedBlockingQueue(2)
        start = time.monotonic()
        self.assertEqual(q.pop_for(0.15), (False, None))
        elapsed = time.monotonic() - start
        self.assertGreaterEqual(elapsed, 0.12, "gave up before the timeout")
        self.assertLess(elapsed, 3.0)

    def test_push_for_times_out_on_a_full_queue(self):
        q = BoundedBlockingQueue(1)
        q.push(1)
        start = time.monotonic()
        self.assertFalse(q.push_for(2, 0.15))
        elapsed = time.monotonic() - start
        self.assertGreaterEqual(elapsed, 0.12, "gave up before the timeout")
        self.assertLess(elapsed, 3.0)
        self.assertEqual(q.size(), 1)
        self.assertEqual(q.pop(), (True, 1))
        self.assertEqual(q.try_pop(), (False, None))  # 2 was not added

    def test_zero_timeout_is_try(self):
        q = BoundedBlockingQueue(1)
        self.assertEqual(q.pop_for(0), (False, None))
        self.assertTrue(q.push_for("a", 0))
        self.assertFalse(q.push_for("b", 0))
        self.assertEqual(q.pop_for(-1), (True, "a"))

    def test_pop_for_returns_as_soon_as_an_item_arrives(self):
        q = BoundedBlockingQueue(2)
        threading.Timer(0.1, lambda: q.push("late")).start()
        start = time.monotonic()
        self.assertEqual(q.pop_for(30.0), (True, "late"))
        self.assertLess(time.monotonic() - start, WAIT)

    def test_push_for_returns_as_soon_as_room_appears(self):
        q = BoundedBlockingQueue(1)
        q.push(1)
        threading.Timer(0.1, q.pop).start()
        start = time.monotonic()
        self.assertTrue(q.push_for(2, 30.0))
        self.assertLess(time.monotonic() - start, WAIT)

    def test_close_ends_a_timed_wait_at_once(self):
        q = BoundedBlockingQueue(1)
        threading.Timer(0.1, q.close).start()
        start = time.monotonic()
        self.assertEqual(q.pop_for(30.0), (False, None))
        self.assertLess(time.monotonic() - start, WAIT)

        q2 = BoundedBlockingQueue(1)
        q2.push(1)
        threading.Timer(0.1, q2.close).start()
        start = time.monotonic()
        self.assertFalse(q2.push_for(2, 30.0))
        self.assertLess(time.monotonic() - start, WAIT)

    def test_the_timeout_covers_the_whole_call(self):
        # Two consumers wait up to 0.6 s. One item arrives after 0.3 s: one
        # consumer gets it, the other must still give up at 0.6 s -- being
        # woken up, finding nothing, and starting the 0.6 s over would make
        # it wait until 0.9 s.
        q = BoundedBlockingQueue(2)
        results = [None, None]
        ends = [0.0, 0.0]
        start = time.monotonic()

        def consumer(i):
            results[i] = q.pop_for(0.6)
            ends[i] = time.monotonic() - start

        threading.Timer(0.3, lambda: q.push("only")).start()
        run_together(2, consumer)
        oks = sorted(r[0] for r in results)
        self.assertEqual(oks, [False, True])
        loser = ends[0] if results[0][0] is False else ends[1]
        self.assertGreaterEqual(loser, 0.5)
        self.assertLess(loser, 0.8, "the timer started over after a wake-up")


class ManyThreadsTest(TimeLimitedTestCase):
    def test_a_woken_consumer_may_find_the_item_gone(self):
        # Consumers sleep in pop(). A fifth thread pushes an item and takes
        # it straight back with try_pop(), so a consumer that is woken up
        # often finds the queue empty again: it has to go back to sleep.
        q = BoundedBlockingQueue(4)
        taken = [[] for _ in range(3)]
        barger_got = []
        rounds = 300

        def body(t):
            if t < 3:
                while True:
                    ok, item = q.pop()
                    if not ok:
                        if not q.closed():
                            raise AssertionError("pop() gave up on an open queue")
                        return
                    taken[t].append(item)
            else:
                for i in range(rounds):
                    q.push(i)
                    ok, item = q.try_pop()
                    if ok:
                        barger_got.append(item)
                q.close()

        with interleaved():
            run_together(4, body, timeout=60)
        everything = sorted(sum(taken, []) + barger_got)
        self.assertEqual(everything, list(range(rounds)))

    def test_a_woken_timed_consumer_keeps_its_deadline(self):
        # The same barging, for pop_for: the consumer is woken up again and
        # again without getting an item, and must still give up on time.
        q = BoundedBlockingQueue(4)
        stop = threading.Event()

        def barger():
            until = time.monotonic() + 1.5
            while not stop.is_set() and time.monotonic() < until:
                q.push("x")
                q.try_pop()
                time.sleep(0.02)

        t = threading.Thread(target=barger, daemon=True)
        t.start()
        start = time.monotonic()
        ok, item = q.pop_for(0.6)
        elapsed = time.monotonic() - start
        stop.set()
        t.join(WAIT)
        if not ok:  # (if the consumer won an item the race tells us nothing)
            self.assertGreaterEqual(elapsed, 0.5)
            self.assertLess(elapsed, 1.0, "the timer started over after a wake-up")

    def _run_mpmc(self, producers, consumers, per_producer, capacity):
        q = BoundedBlockingQueue(capacity)
        taken = [[] for _ in range(consumers)]
        finished = AtomicCounter()
        too_big = AtomicCounter()

        def body(t):
            if t < producers:
                for i in range(per_producer):
                    if not q.push((t, i)):
                        raise AssertionError("push returned False, queue not closed")
                finished.add()
                if finished.value == producers:
                    q.close()  # the last producer to finish closes the queue
            else:
                mine = taken[t - producers]
                while True:
                    ok, item = q.pop()
                    if not ok:
                        if not q.closed():
                            raise AssertionError("pop() gave up on an open queue")
                        break
                    mine.append(item)
                    if q.size() > capacity:
                        too_big.add()

        with interleaved():
            run_together(producers + consumers, body, timeout=60)

        self.assertEqual(too_big.value, 0, "size() went above capacity()")
        everything = [item for mine in taken for item in mine]
        expected = [(t, i) for t in range(producers) for i in range(per_producer)]
        self.assertEqual(sorted(everything), expected,
                         "an item was lost, handed out twice, or made up")
        # Items of one producer must stay in order for any one consumer.
        for mine in taken:
            last = {}
            for t, i in mine:
                self.assertGreater(i, last.get(t, -1), "FIFO order broken")
                last[t] = i

    def test_four_producers_four_consumers(self):
        self._run_mpmc(4, 4, 150, 8)

    def test_capacity_one_forces_a_hand_over_for_every_item(self):
        self._run_mpmc(3, 3, 100, 1)

    def test_one_producer_one_consumer_ping_pong(self):
        self._run_mpmc(1, 1, 300, 1)

    def test_try_calls_under_contention(self):
        q = BoundedBlockingQueue(4)
        pushed = [[] for _ in range(4)]
        popped = [[] for _ in range(4)]

        def body(t):
            for i in range(150):
                if t % 2 == 0:
                    if q.try_push((t, i)):
                        pushed[t].append((t, i))
                else:
                    ok, item = q.try_pop()
                    if ok:
                        popped[t].append(item)

        with interleaved():
            run_together(4, body, timeout=60)
        rest = []
        while True:
            ok, item = q.try_pop()
            if not ok:
                break
            rest.append(item)
        got = sorted(sum(popped, []) + rest)
        self.assertEqual(got, sorted(sum(pushed, [])))


if __name__ == "__main__":
    unittest.main()
