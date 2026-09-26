import threading
import time
import unittest

from ._concurrency import (AtomicCounter, Background, TimeLimitedTestCase,
                           interleaved, run_in_thread, run_together)
from .once_cache import OnceCache


class OnceCacheTest(TimeLimitedTestCase):
    def test_miss_computes_and_returns_the_value(self):
        c = OnceCache()
        calls = 0

        def compute():
            nonlocal calls
            calls += 1
            return 42

        self.assertEqual(c.get("a", compute), 42)
        self.assertEqual(calls, 1)
        self.assertEqual(c.size(), 1)

    def test_hit_returns_remembered_value_without_computing(self):
        c = OnceCache()
        c.get("a", lambda: 42)
        called = False

        def compute():
            nonlocal called
            called = True
            return 7

        self.assertEqual(c.get("a", compute), 42)
        self.assertFalse(called)

    def test_different_keys_get_their_own_values(self):
        c = OnceCache()
        self.assertEqual(c.get("a", lambda: "apple"), "apple")
        self.assertEqual(c.get("b", lambda: "banana"), "banana")
        self.assertEqual(c.get("a", lambda: "avocado"), "apple")
        self.assertEqual(c.size(), 2)

    def test_throwing_compute_remembers_nothing(self):
        c = OnceCache()
        boom = RuntimeError("boom")

        def failing():
            raise boom

        with self.assertRaises(RuntimeError) as caught:
            c.get("a", failing)
        self.assertIs(caught.exception, boom, "get() must raise the exception compute() raised")
        self.assertEqual(c.size(), 0)

        calls = 0

        def compute():
            nonlocal calls
            calls += 1
            return 5

        # On another thread, so a get() that waits for the failed
        # computation fails this test instead of hanging it.
        value = run_in_thread(
            lambda: c.get("a", compute), timeout=2,
            msg="get() after a failed compute() never returned -- is it "
                "waiting for the failed computation?")
        self.assertEqual(value, 5)
        self.assertEqual(calls, 1, "after a failure, the next get() must compute again")
        self.assertEqual(c.size(), 1)

    def test_concurrent_gets_of_one_key_compute_once(self):
        # 16 threads miss the same key at the same moment; compute() is slow,
        # so a check-then-act gap is wide open.
        threads = 16
        c = OnceCache()
        calls = AtomicCounter()
        results = [None] * threads

        def slow_compute():
            calls.add()
            time.sleep(0.05)
            return 99

        def get_config(t):
            results[t] = c.get("config", slow_compute)

        with interleaved():
            run_together(threads, get_config)

        self.assertEqual(calls.value, 1, "compute() ran more than once for one key")
        for r in results:
            self.assertEqual(r, 99)
        self.assertEqual(c.size(), 1)

    def test_concurrent_gets_of_many_keys_compute_each_once(self):
        threads = 8
        keys = 200
        c = OnceCache()
        calls = [AtomicCounter() for _ in range(keys)]

        def compute_for(k):
            def compute():
                calls[k].add()
                time.sleep(0.001)
                return k * 10
            return compute

        def get_every_key(t):
            for i in range(keys):
                k = (i + t * 25) % keys  # each thread starts at a different key
                self.assertEqual(c.get(str(k), compute_for(k)), k * 10)

        with interleaved():
            run_together(threads, get_every_key)

        for k in range(keys):
            self.assertEqual(calls[k].value, 1, f"key {k} was computed more than once")
        self.assertEqual(c.size(), keys)

    def test_different_keys_do_not_wait_for_each_other(self):
        # compute("a") and compute("b") each wait until the other one has
        # started. That only works if both can run at the same time; if a
        # lock is held while compute() runs, one of them gives up after 2 s.
        c = OnceCache()
        started = {"a": threading.Event(), "b": threading.Event()}
        got = {}

        def get_one(t):
            key, other = ("a", "b") if t == 0 else ("b", "a")

            def compute():
                started[key].set()
                return "ok" if started[other].wait(2) else "timed out"

            got[key] = c.get(key, compute)

        with interleaved():
            run_together(2, get_one)

        self.assertEqual(got["a"], "ok", 'compute("a") waited for compute("b")')
        self.assertEqual(got["b"], "ok", 'compute("b") waited for compute("a")')

    def test_size_counts_only_finished_values(self):
        c = OnceCache()
        started = threading.Event()
        release = threading.Event()

        def slow():
            started.set()
            release.wait(10)
            return 1

        with interleaved():
            getter = Background(lambda: c.get("slow", slow))
            try:
                # Wait (up to 5 s) until compute() is running -- or until
                # get() returned without ever calling it.
                deadline = time.monotonic() + 5
                while not started.is_set() and time.monotonic() < deadline:
                    if getter.wait(0.01):
                        break
                did_start = started.is_set()
                # Ask size() on another thread, so a size() that waits for the
                # running compute() fails this test instead of hanging it.
                asker = Background(c.size)
                size_answered = asker.wait(2)
            finally:
                release.set()
            getter.wait(5)
            asker.wait(5)

        self.assertTrue(did_start, "compute() was never called")
        self.assertTrue(size_answered, "size() waited for a running compute()")
        self.assertEqual(asker.result(), 0, "a computation still running counted as remembered")
        self.assertEqual(c.size(), 1)


if __name__ == "__main__":
    unittest.main()
