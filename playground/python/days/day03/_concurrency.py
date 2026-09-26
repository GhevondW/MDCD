"""Test helpers for Day 3.

CPython runs one thread at a time (the GIL) and switches between threads
only at certain points, so many races that break code in C++, Java or Go
stay hidden in Python -- rare, not impossible. Inside `interleaved()`,
before each line of the day03 solution files runs, the thread briefly
gives up the GIL, so another thread can run in between -- the way real
threads on a multi-core machine can.
"""

import os
import select
import sys
import threading
import time
import traceback
import unittest
from contextlib import contextmanager

_HERE = os.path.dirname(os.path.abspath(__file__))
_is_solution = {}


def _solution_file(filename):
    known = _is_solution.get(filename)
    if known is None:
        path = os.path.abspath(filename)
        name = os.path.basename(path)
        known = (os.path.dirname(path) == _HERE
                 and not name.startswith(("test_", "_")))
        _is_solution[filename] = known
    return known


# Any blocking call gives up the GIL while it runs. time.sleep(0) would
# do, but on Linux from Python 3.11 on it really sleeps (~70 microseconds
# a line), which makes these tests ~20x slower there. A select() on no
# files with a zero timeout returns at once. Windows' select() needs at
# least one socket -- but there, sleep(0) just gives up the time slice.
if sys.platform == "win32":
    def _give_up_gil():
        time.sleep(0)
else:
    def _give_up_gil():
        select.select([], [], [], 0)


def _each_line(frame, event, arg):
    if event == "line":
        _give_up_gil()  # another thread may run now
    return _each_line


def _on_call(frame, event, arg):
    if _solution_file(frame.f_code.co_filename):
        return _each_line
    return None


@contextmanager
def interleaved():
    """Let other threads run between any two lines of the solution files.

    Only threads started inside the `with` block are affected, so start
    them there.
    """
    threading.settrace(_on_call)
    sys.settrace(_on_call)
    try:
        yield
    finally:
        sys.settrace(None)
        threading.settrace(None)


# ---- running code on other threads, with a time limit ----------------

# Threads that are still running after this long are taken to be
# deadlocked. With a correct solution, no test needs even one second.
RUN_TIMEOUT = 10.0


def _stuck_at(thread):
    """Where `thread` is right now: its stack, most recent call last,
    showing only the day03 test and solution files."""
    frame = sys._current_frames().get(thread.ident)
    if frame is None:
        return "  (the thread has finished)\n"
    stack = traceback.extract_stack(frame)
    ours = [f for f in stack
            if os.path.dirname(os.path.abspath(f.filename)) == _HERE
            and os.path.basename(f.filename) != "_concurrency.py"]
    return "".join(traceback.format_list(ours or stack))


def _join_all(threads, errors, timeout, msg):
    deadline = time.monotonic() + timeout
    for thread in threads:
        thread.join(max(0.0, deadline - time.monotonic()))
    if errors:
        raise errors[0]
    stuck = [thread for thread in threads if thread.is_alive()]
    if stuck:
        lines = [msg] if msg else []
        lines.append(f"{len(stuck)} of {len(threads)} thread(s) still running "
                     f"after {timeout:g} s -- a deadlock? One of them is stuck here:")
        lines.append(_stuck_at(stuck[0]))
        raise AssertionError("\n".join(lines))


def run_together(n, body, timeout=RUN_TIMEOUT, msg=""):
    """Run body(0) ... body(n-1) on n threads that all start at the same
    moment, so their calls overlap as much as possible, and wait for all
    of them.

    If a thread raised, that exception is raised here. If a thread is
    still running after `timeout` seconds, the test fails (with `msg`)
    instead of hanging.
    """
    barrier = threading.Barrier(n)
    errors = []

    def worker(t):
        try:
            barrier.wait()
            body(t)
        except BaseException as e:  # hand it to the test's own thread
            errors.append(e)

    # daemon: a deadlocked thread must not keep Python from exiting.
    threads = [threading.Thread(target=worker, args=(t,), daemon=True)
               for t in range(n)]
    for thread in threads:
        thread.start()
    _join_all(threads, errors, timeout, msg)


def run_in_thread(fn, timeout=RUN_TIMEOUT, msg=""):
    """Call fn() on another thread and return what it returns (or raise
    what it raised). If it is still running after `timeout` seconds, the
    test fails (with `msg`) instead of hanging -- for calls that might
    deadlock.
    """
    result = []
    run_together(1, lambda _: result.append(fn()), timeout, msg)
    return result[0]


class Background:
    """fn() running on a daemon thread, started right away -- for tests
    that need to check on a call while it is still running."""

    def __init__(self, fn):
        self._done = threading.Event()
        self._result = None
        self._error = None

        def run():
            try:
                self._result = fn()
            except BaseException as e:
                self._error = e
            finally:
                self._done.set()

        threading.Thread(target=run, daemon=True).start()

    def wait(self, timeout):
        """Wait up to `timeout` seconds; return whether fn() has returned."""
        return self._done.wait(timeout)

    def result(self):
        """fn()'s return value, or raise what fn() raised. Call only after
        wait() returned True."""
        if self._error is not None:
            raise self._error
        return self._result


class AtomicCounter:
    """An int that many threads can add to safely -- the tests' own
    bookkeeping (std::atomic<int> in the C++ tests)."""

    def __init__(self):
        self._lock = threading.Lock()
        self._value = 0

    def add(self, n=1):
        with self._lock:
            self._value += n

    @property
    def value(self):
        with self._lock:
            return self._value


# ---- a time limit for every test --------------------------------------

TEST_TIME_LIMIT = 30  # seconds


def _stop_the_run(test_id, thread):
    out = sys.__stderr__
    if out is not None:
        out.write(f"\n\n{test_id} is still running after {TEST_TIME_LIMIT} s "
                  f"-- a deadlock?\nStopping the test run. The test is stuck "
                  f"here:\n{_stuck_at(thread)}\n")
        out.flush()
    os._exit(1)


class TimeLimitedTestCase(unittest.TestCase):
    """A TestCase whose tests may run for at most TEST_TIME_LIMIT seconds.

    Most deadlocks make a test fail with a message (see run_together). If
    one hangs a test anyway -- say, a thread waits for a lock it already
    holds -- then after TEST_TIME_LIMIT seconds this prints where the test
    is stuck and stops the test run, instead of hanging forever.
    """

    def setUp(self):
        super().setUp()
        watchdog = threading.Timer(TEST_TIME_LIMIT, _stop_the_run,
                                   args=(self.id(), threading.current_thread()))
        watchdog.daemon = True
        watchdog.start()
        self.addCleanup(watchdog.cancel)
