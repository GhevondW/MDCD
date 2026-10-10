import threading
import time
import unittest
from concurrent.futures import TimeoutError as FutureTimeout

from ._concurrency import (AtomicCounter, Background, TimeLimitedTestCase,
                           interleaved, run_together)
from .thread_pool import ThreadPool

WAIT = 3.0  # a result that should come "soon" gets this long


def wait_for(test, event, what):
    test.assertTrue(event.wait(WAIT), what)


class BasicsTest(TimeLimitedTestCase):
    def test_arguments_are_checked(self):
        with self.assertRaises(ValueError):
            ThreadPool(0)
        with self.assertRaises(ValueError):
            ThreadPool(-1)
        with self.assertRaises(ValueError):
            ThreadPool(1, -1)
        pool = ThreadPool(3)
        self.assertEqual(pool.worker_count(), 3)
        pool.shutdown()

    def test_submit_returns_the_result(self):
        with ThreadPool(2) as pool:
            f = pool.submit(lambda: 6 * 7)
            g = pool.submit(pow, 2, 10)
            self.assertEqual(f.result(WAIT), 42)
            self.assertEqual(g.result(WAIT), 1024)

    def test_an_exception_travels_through_the_future(self):
        with ThreadPool(2) as pool:
            def boom():
                raise KeyError("boom")

            f = pool.submit(boom)
            with self.assertRaises(KeyError):
                f.result(WAIT)
            # the pool still works
            self.assertEqual(pool.submit(lambda: 1).result(WAIT), 1)

    def test_every_task_runs_exactly_once(self):
        counts = [0] * 500
        lock = threading.Lock()

        def bump(i):
            with lock:
                counts[i] += 1

        with ThreadPool(4) as pool:
            futures = [pool.submit(bump, i) for i in range(500)]
            for f in futures:
                f.result(WAIT)
        self.assertEqual(counts, [1] * 500)

    def test_tasks_run_on_the_workers_not_on_the_callers_thread(self):
        idents = set()
        lock = threading.Lock()

        def note():
            with lock:
                idents.add(threading.get_ident())

        with ThreadPool(3) as pool:
            for f in [pool.submit(note) for _ in range(200)]:
                f.result(WAIT)
        self.assertNotIn(threading.get_ident(), idents)
        self.assertLessEqual(len(idents), 3, "more threads than workers")

    def test_workers_run_tasks_at_the_same_time(self):
        # 4 tasks each wait for the other 3 to arrive. With 4 real workers
        # they all meet; with fewer, the barrier times out.
        barrier = threading.Barrier(4, timeout=WAIT)
        with ThreadPool(4) as pool:
            futures = [pool.submit(barrier.wait) for _ in range(4)]
            for f in futures:
                f.result(WAIT * 2)

    def test_a_task_may_hand_out_more_work(self):
        done = AtomicCounter()
        all_done = threading.Event()

        def child():
            done.add()
            if done.value == 11:
                all_done.set()

        def parent():
            for _ in range(10):
                pool.execute(child)
            done.add()
            if done.value == 11:
                all_done.set()

        with ThreadPool(2) as pool:
            pool.execute(parent)
            wait_for(self, all_done, "tasks created by a task never ran")


class ErrorsAndShutdownTest(TimeLimitedTestCase):
    def test_a_raising_task_does_not_kill_its_worker(self):
        def boom():
            raise RuntimeError("boom")

        with ThreadPool(1) as pool:
            for _ in range(5):
                pool.execute(boom)
            self.assertEqual(pool.submit(lambda: "alive").result(WAIT), "alive")

    def test_shutdown_runs_every_accepted_task_and_waits_for_them(self):
        gate = threading.Event()
        started = threading.Event()
        ran = AtomicCounter()

        def first():
            started.set()
            gate.wait(WAIT)
            ran.add()

        pool = ThreadPool(1)
        pool.execute(first)
        wait_for(self, started, "the first task never started")
        for _ in range(50):
            pool.execute(ran.add)
        b = Background(pool.shutdown)
        self.assertFalse(b.wait(0.3), "shutdown() returned while tasks were still queued")
        gate.set()
        self.assertTrue(b.wait(WAIT), "shutdown() never returned")
        b.result()
        self.assertEqual(ran.value, 51, "an accepted task was dropped")

    def test_execute_and_try_submit_after_shutdown(self):
        pool = ThreadPool(2)
        pool.shutdown()
        with self.assertRaises(RuntimeError):
            pool.execute(lambda: None)
        with self.assertRaises(RuntimeError):
            pool.submit(lambda: None)
        self.assertFalse(pool.try_submit(lambda: None))

    def test_shutdown_twice_and_from_many_threads(self):
        pool = ThreadPool(3)
        pool.shutdown()
        pool.shutdown()
        pool2 = ThreadPool(3)
        pool2.execute(lambda: time.sleep(0.05))
        run_together(4, lambda t: pool2.shutdown())

    def test_shutdown_really_stops_the_worker_threads(self):
        before = threading.active_count()
        pool = ThreadPool(4)
        self.assertEqual(threading.active_count(), before + 4)
        pool.submit(lambda: None).result(WAIT)
        pool.shutdown()
        self.assertEqual(threading.active_count(), before,
                         "worker threads are still alive after shutdown()")

    def test_shutdown_with_nothing_to_do_returns(self):
        pool = ThreadPool(4)
        b = Background(pool.shutdown)
        self.assertTrue(b.wait(WAIT), "idle workers slept through shutdown()")
        b.result()


class BoundedQueueTest(TimeLimitedTestCase):
    def _blocked_pool(self):
        """A 1-worker pool with room for 2 waiting tasks, whose worker is
        stuck in a task until `gate` is set, and whose queue is full."""
        gate = threading.Event()
        started = threading.Event()
        pool = ThreadPool(1, queue_capacity=2)

        def first():
            started.set()
            gate.wait(WAIT)

        pool.execute(first)
        wait_for(self, started, "the first task never started")
        return pool, gate

    def test_try_submit_says_no_when_the_queue_is_full(self):
        pool, gate = self._blocked_pool()
        ran = AtomicCounter()
        self.assertTrue(pool.try_submit(ran.add))
        self.assertTrue(pool.try_submit(ran.add))
        self.assertFalse(pool.try_submit(ran.add), "queue should be full")
        gate.set()
        pool.shutdown()
        self.assertEqual(ran.value, 2, "the rejected task must not run")

    def test_execute_waits_for_room(self):
        pool, gate = self._blocked_pool()
        ran = AtomicCounter()
        pool.execute(ran.add)
        pool.execute(ran.add)
        b = Background(lambda: pool.execute(ran.add))
        self.assertFalse(b.wait(0.3), "execute() did not wait on a full queue")
        gate.set()
        self.assertTrue(b.wait(WAIT), "execute() was never woken up")
        b.result()
        pool.shutdown()
        self.assertEqual(ran.value, 3)

    def test_shutdown_rejects_an_execute_that_is_waiting(self):
        pool, gate = self._blocked_pool()
        ran = AtomicCounter()
        pool.execute(ran.add)
        pool.execute(ran.add)
        waiting = Background(lambda: pool.execute(ran.add))
        self.assertFalse(waiting.wait(0.3))
        stopper = Background(pool.shutdown)
        self.assertTrue(waiting.wait(WAIT),
                        "execute() kept waiting although the pool was shut down")
        with self.assertRaises(RuntimeError):
            waiting.result()
        gate.set()
        self.assertTrue(stopper.wait(WAIT))
        stopper.result()
        self.assertEqual(ran.value, 2, "the task that was rejected must not run")


class ManyThreadsTest(TimeLimitedTestCase):
    def test_many_threads_submit_at_once(self):
        ran = AtomicCounter()
        pool = ThreadPool(4)
        with interleaved():
            run_together(4, lambda t: [pool.execute(ran.add) for _ in range(150)],
                         timeout=60)
        pool.shutdown()
        self.assertEqual(ran.value, 600)

    def test_many_threads_submit_into_a_small_queue(self):
        ran = AtomicCounter()
        pool = ThreadPool(2, queue_capacity=3)
        with interleaved():
            run_together(4, lambda t: [pool.execute(ran.add) for _ in range(100)],
                         timeout=60)
        pool.shutdown()
        self.assertEqual(ran.value, 400)

    def test_submits_racing_with_shutdown_are_run_or_rejected_never_lost(self):
        for _ in range(10):
            ran = AtomicCounter()
            accepted = AtomicCounter()
            pool = ThreadPool(3, queue_capacity=4)

            def body(t):
                if t == 0:
                    time.sleep(0.005)
                    pool.shutdown()
                    return
                for _ in range(200):
                    try:
                        pool.execute(ran.add)
                    except RuntimeError:
                        return
                    accepted.add()

            run_together(4, body, timeout=60)
            pool.shutdown()
            self.assertEqual(ran.value, accepted.value,
                             "a task that execute() accepted never ran")


if __name__ == "__main__":
    unittest.main()
