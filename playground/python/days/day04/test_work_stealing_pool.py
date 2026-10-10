import threading
import time
import unittest

from ._concurrency import (AtomicCounter, Background, TimeLimitedTestCase,
                           run_together)
from .work_stealing_pool import WorkStealingPool

WAIT = 3.0  # a result that should come "soon" gets this long


def fib(pool, n):
    # Fork one half, compute the other half here, then wait for the fork.
    if n < 2:
        return n
    left = pool.submit(lambda: fib(pool, n - 1))
    right = fib(pool, n - 2)
    return left.result() + right


def total(pool, numbers, lo, hi):
    if hi - lo <= 500:
        return sum(numbers[lo:hi])
    mid = (lo + hi) // 2
    left = pool.submit(lambda: total(pool, numbers, lo, mid))
    right = total(pool, numbers, mid, hi)
    return left.result() + right


class BasicsTest(TimeLimitedTestCase):
    def test_arguments_are_checked(self):
        with self.assertRaises(ValueError):
            WorkStealingPool(0)
        pool = WorkStealingPool(3)
        self.assertEqual(pool.worker_count(), 3)
        pool.shutdown()

    def test_submit_returns_the_result(self):
        with WorkStealingPool(2) as pool:
            self.assertEqual(pool.submit(lambda: 6 * 7).result(WAIT), 42)

    def test_an_exception_travels_through_the_handle(self):
        def boom():
            raise KeyError("boom")

        with WorkStealingPool(2) as pool:
            with self.assertRaises(KeyError):
                pool.submit(boom).result(WAIT)
            self.assertEqual(pool.submit(lambda: 1).result(WAIT), 1)

    def test_a_raising_task_does_not_kill_its_worker(self):
        def boom():
            raise RuntimeError("boom")

        with WorkStealingPool(1) as pool:
            for _ in range(5):
                pool.execute(boom)
            self.assertEqual(pool.submit(lambda: "alive").result(WAIT), "alive")

    def test_in_worker_tells_workers_from_the_rest(self):
        with WorkStealingPool(2) as pool:
            self.assertFalse(pool.in_worker())
            self.assertTrue(pool.submit(pool.in_worker).result(WAIT))
            with WorkStealingPool(1) as other:
                # a worker of one pool is not a worker of another
                self.assertFalse(pool.submit(other.in_worker).result(WAIT))


class HelpingJoinTest(TimeLimitedTestCase):
    def test_fib_needs_no_free_worker_to_wait_for_a_child(self):
        # Waiting for a child inside a task would deadlock a plain pool
        # with few workers. Here the waiting worker runs the child itself.
        for workers in (1, 2, 4):
            with self.subTest(workers=workers):
                with WorkStealingPool(workers) as pool:
                    self.assertEqual(pool.submit(lambda: fib(pool, 16)).result(WAIT), 987)

    def test_divide_and_conquer_sum(self):
        numbers = list(range(20000))
        for workers in (1, 3):
            with self.subTest(workers=workers):
                with WorkStealingPool(workers) as pool:
                    got = pool.submit(lambda: total(pool, numbers, 0, len(numbers))).result(WAIT)
                self.assertEqual(got, sum(numbers))

    def test_one_worker_never_steals(self):
        with WorkStealingPool(1) as pool:
            pool.submit(lambda: fib(pool, 12)).result(WAIT)
            self.assertEqual(pool.steals(), 0)

    def test_two_parents_waiting_at_once_on_two_workers(self):
        # The lecture's "deadlock with no locks at all".
        with WorkStealingPool(2) as pool:
            def parent():
                child = pool.submit(lambda: 1)
                return child.result() + 1

            a = pool.submit(parent)
            b = pool.submit(parent)
            self.assertEqual(a.result(WAIT), 2)
            self.assertEqual(b.result(WAIT), 2)

    def test_waiting_on_a_task_a_thief_is_running(self):
        # The child is taken by another worker and runs for a while: the
        # waiting worker must wait for it, not give up or run it twice.
        runs = AtomicCounter()
        with WorkStealingPool(2) as pool:
            def parent():
                def child():
                    runs.add()
                    time.sleep(0.3)
                    return "child"

                h = pool.submit(child)
                time.sleep(0.1)  # let the other worker steal it
                return h.result()

            self.assertEqual(pool.submit(parent).result(WAIT), "child")
        self.assertEqual(runs.value, 1)


class OrderTest(TimeLimitedTestCase):
    def test_a_worker_runs_its_own_newest_task_first(self):
        order = []
        all_ran = threading.Event()

        def record(name):
            order.append(name)
            if len(order) == 3:
                all_ran.set()

        with WorkStealingPool(1) as pool:
            def parent():
                for name in "ABC":
                    pool.execute(lambda name=name: record(name))

            pool.execute(parent)
            self.assertTrue(all_ran.wait(WAIT), "forked tasks never ran")
        self.assertEqual(order, ["C", "B", "A"])

    def test_a_thief_takes_the_oldest_task(self):
        order = []
        all_ran = threading.Event()

        def record(name):
            order.append(name)
            if len(order) == 3:
                all_ran.set()

        with WorkStealingPool(2) as pool:
            def parent():
                for name in "ABC":
                    pool.execute(lambda name=name: record(name))
                # Stay busy (without running anything) until the other
                # worker has stolen all three.
                all_ran.wait(WAIT)

            pool.execute(parent)
            self.assertTrue(all_ran.wait(WAIT), "the idle worker never stole")
            self.assertEqual(order, ["A", "B", "C"])
            self.assertGreaterEqual(pool.steals(), 3)

    def test_tasks_from_outside_run_oldest_first(self):
        order = []
        gate = threading.Event()
        started = threading.Event()
        all_ran = threading.Event()

        def first():
            started.set()
            gate.wait(WAIT)

        def record(i):
            order.append(i)
            if len(order) == 5:
                all_ran.set()

        with WorkStealingPool(1) as pool:
            pool.execute(first)
            self.assertTrue(started.wait(WAIT))
            for i in range(5):
                pool.execute(lambda i=i: record(i))
            gate.set()
            self.assertTrue(all_ran.wait(WAIT))
        self.assertEqual(order, [0, 1, 2, 3, 4])
        self.assertEqual(pool.steals(), 0, "the injection queue is not a steal")


class StealingTest(TimeLimitedTestCase):
    def test_idle_workers_take_over_work_that_one_worker_created(self):
        # One task creates 100 tasks that each sleep 10 ms: 1 s of work.
        # They all sit in ONE worker's deque; the other 3 workers must
        # steal, so the whole thing takes about a quarter of that.
        with WorkStealingPool(4) as pool:
            def parent():
                children = [pool.submit(lambda: time.sleep(0.01)) for _ in range(100)]
                for c in children:
                    c.result()

            start = time.monotonic()
            pool.submit(parent).result(WAIT)
            elapsed = time.monotonic() - start
            self.assertGreater(pool.steals(), 0, "nobody stole anything")
            self.assertLess(elapsed, 0.7,
                            "took %.2f s; 1 s of work on 4 workers should take ~0.25 s" % elapsed)

    def test_many_outside_tasks_with_many_outside_threads(self):
        ran = AtomicCounter()
        pool = WorkStealingPool(4)

        def body(t):
            handles = [pool.submit(lambda: (ran.add(), t)[1]) for _ in range(100)]
            for h in handles:
                if h.result(WAIT) != t:
                    raise AssertionError("got another task's result")

        run_together(4, body, timeout=60)
        pool.shutdown()
        self.assertEqual(ran.value, 400)


class ShutdownTest(TimeLimitedTestCase):
    def test_outside_execute_after_shutdown_is_rejected(self):
        pool = WorkStealingPool(2)
        pool.shutdown()
        with self.assertRaises(RuntimeError):
            pool.execute(lambda: None)
        with self.assertRaises(RuntimeError):
            pool.submit(lambda: None)

    def test_shutdown_twice_and_from_many_threads(self):
        pool = WorkStealingPool(3)
        pool.shutdown()
        pool.shutdown()
        pool2 = WorkStealingPool(3)
        pool2.execute(lambda: time.sleep(0.05))
        run_together(4, lambda t: pool2.shutdown())

    def test_shutdown_really_stops_the_worker_threads(self):
        before = threading.active_count()
        pool = WorkStealingPool(4)
        self.assertEqual(threading.active_count(), before + 4)
        pool.submit(lambda: None).result(WAIT)
        pool.shutdown()
        self.assertEqual(threading.active_count(), before,
                         "worker threads are still alive after shutdown()")

    def test_idle_workers_sleep_through_nothing_but_wake_for_shutdown(self):
        pool = WorkStealingPool(4)
        b = Background(pool.shutdown)
        self.assertTrue(b.wait(WAIT), "idle workers slept through shutdown()")
        b.result()

    def test_shutdown_finishes_work_that_tasks_create_while_it_runs(self):
        # A binary tree of tasks, 2^8 - 1 = 255 of them: every task forks two
        # more (with execute, no waiting). shutdown() is called right after
        # the root went in, so most of the tree is created AFTER shutdown()
        # began -- all of it must still run.
        ran = AtomicCounter()

        def node(depth):
            ran.add()
            time.sleep(0.001)
            if depth < 7:
                pool.execute(lambda: node(depth + 1))
                pool.execute(lambda: node(depth + 1))

        pool = WorkStealingPool(3)
        pool.execute(lambda: node(0))
        b = Background(pool.shutdown)
        self.assertTrue(b.wait(15), "shutdown() never returned")
        b.result()
        self.assertEqual(ran.value, 255, "work created during shutdown was dropped")

    def test_shutdown_waits_for_a_running_task(self):
        gate = threading.Event()
        started = threading.Event()
        ran = AtomicCounter()

        def slow():
            started.set()
            gate.wait(WAIT)
            ran.add()

        pool = WorkStealingPool(2)
        pool.execute(slow)
        self.assertTrue(started.wait(WAIT))
        b = Background(pool.shutdown)
        self.assertFalse(b.wait(0.3), "shutdown() returned with a task still running")
        gate.set()
        self.assertTrue(b.wait(WAIT))
        b.result()
        self.assertEqual(ran.value, 1)


if __name__ == "__main__":
    unittest.main()
