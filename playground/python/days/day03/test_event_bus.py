import unittest

from ._concurrency import (AtomicCounter, Background, TimeLimitedTestCase,
                           interleaved, run_in_thread, run_together)
from .event_bus import EventBus


def publish(bus, event):
    """bus.publish(event), run on another thread. The callbacks in these
    tests call back into the bus; if that deadlocks, the test fails after
    2 s instead of hanging."""
    return run_in_thread(
        lambda: bus.publish(event), timeout=2,
        msg="publish() never returned -- a callback that used the bus deadlocked")


class EventBusTest(TimeLimitedTestCase):
    def test_publish_calls_subscribers_in_order(self):
        bus = EventBus()
        log = []
        bus.subscribe(lambda e: log.append("A:" + e))
        bus.subscribe(lambda e: log.append("B:" + e))
        bus.subscribe(lambda e: log.append("C:" + e))
        self.assertEqual(bus.publish("x"), 3)
        self.assertEqual(log, ["A:x", "B:x", "C:x"])

    def test_publish_with_no_subscribers_calls_nobody(self):
        bus = EventBus()
        self.assertEqual(bus.publish("x"), 0)
        self.assertEqual(bus.subscriber_count(), 0)

    def test_ids_are_positive_and_unique(self):
        bus = EventBus()
        ids = set()
        for _ in range(100):
            sub_id = bus.subscribe(lambda e: None)
            self.assertGreater(sub_id, 0)
            ids.add(sub_id)
        self.assertEqual(len(ids), 100)
        self.assertEqual(bus.subscriber_count(), 100)

    def test_unsubscribe_stops_delivery(self):
        bus = EventBus()
        log = []
        bus.subscribe(lambda e: log.append("A"))
        b = bus.subscribe(lambda e: log.append("B"))
        bus.subscribe(lambda e: log.append("C"))

        self.assertTrue(bus.unsubscribe(b))
        self.assertFalse(bus.unsubscribe(b), "already removed")
        self.assertFalse(bus.unsubscribe(12345), "never existed")
        self.assertEqual(bus.subscriber_count(), 2)
        self.assertEqual(bus.publish("x"), 2)
        self.assertEqual(log, ["A", "C"])

    def test_order_is_kept_after_unsubscribe_and_subscribe(self):
        bus = EventBus()
        log = []
        bus.subscribe(lambda e: log.append("A"))
        b = bus.subscribe(lambda e: log.append("B"))
        bus.subscribe(lambda e: log.append("C"))
        bus.unsubscribe(b)
        bus.subscribe(lambda e: log.append("D"))
        bus.publish("x")
        self.assertEqual(log, ["A", "C", "D"])

    def test_subscribe_inside_a_callback_counts_from_the_next_publish(self):
        bus = EventBus()
        log = []
        added = False

        def a(event):
            nonlocal added
            log.append("A")
            if not added:
                added = True
                bus.subscribe(lambda e: log.append("B"))

        bus.subscribe(a)

        self.assertEqual(publish(bus, "first"), 1, "B joined mid-publish: not called by it")
        self.assertEqual(log, ["A"])
        log.clear()
        self.assertEqual(publish(bus, "second"), 2)
        self.assertEqual(log, ["A", "B"])

    def test_unsubscribe_inside_a_callback_counts_from_the_next_publish(self):
        bus = EventBus()
        log = []
        b = 0
        removed = None  # what unsubscribe(b) returned inside the callback

        def a(event):
            nonlocal removed
            log.append("A")
            if removed is None:
                removed = bus.unsubscribe(b)

        bus.subscribe(a)
        b = bus.subscribe(lambda e: log.append("B"))

        self.assertEqual(publish(bus, "first"), 2, "B left mid-publish: still called by it")
        self.assertEqual(log, ["A", "B"])
        self.assertTrue(removed, "unsubscribe() inside a callback did not remove B")
        log.clear()
        self.assertEqual(publish(bus, "second"), 1)
        self.assertEqual(log, ["A"])

    def test_callback_can_unsubscribe_itself(self):
        bus = EventBus()
        a_calls = 0

        def a_leaves(event):
            nonlocal a_calls
            a_calls += 1
            bus.unsubscribe(a)

        a = bus.subscribe(a_leaves)
        bus.subscribe(lambda e: None)

        self.assertEqual(publish(bus, "first"), 2)
        self.assertEqual(publish(bus, "second"), 1)
        self.assertEqual(a_calls, 1)

    def test_callback_can_publish(self):
        bus = EventBus()
        b_log = []

        def a(event):
            if event == "outer":
                bus.publish("inner")

        bus.subscribe(a)
        bus.subscribe(lambda e: b_log.append(e))

        self.assertEqual(publish(bus, "outer"), 2)
        # A handles "outer" by publishing "inner", which reaches B first.
        self.assertEqual(b_log, ["inner", "outer"])

    def test_other_threads_can_use_the_bus_while_a_callback_runs(self):
        # The callback waits (up to 2 s) for another thread to subscribe. If
        # publish() holds a lock while calling back, that thread is stuck.
        bus = EventBus()
        outcome = None
        helper = None

        def wait_for_helper(event):
            nonlocal outcome, helper
            helper = Background(lambda: bus.subscribe(lambda e: None))
            outcome = "ok" if helper.wait(2) else "blocked"

        bus.subscribe(wait_for_helper)
        with interleaved():
            bus.publish("x")
            if helper is not None:
                helper.wait(5)

        self.assertEqual(outcome, "ok", "another thread could not subscribe while a callback ran")
        self.assertEqual(bus.subscriber_count(), 2)

    def test_concurrent_subscribes_get_unique_ids(self):
        threads = 8
        per_thread = 250
        bus = EventBus()
        ids = [[] for _ in range(threads)]

        def subscribe_many(t):
            for _ in range(per_thread):
                ids[t].append(bus.subscribe(lambda e: None))

        with interleaved():
            run_together(threads, subscribe_many)

        unique = {sub_id for mine in ids for sub_id in mine}
        self.assertEqual(len(unique), threads * per_thread, "two subscribers got the same id")
        self.assertEqual(bus.subscriber_count(), threads * per_thread)

    def test_concurrent_publish_and_churn(self):
        # 3 subscribers stay for the whole test and must hear every event
        # exactly once, while other threads keep subscribing and
        # unsubscribing around them.
        publishers = 4
        churners = 4
        per_thread = 500
        bus = EventBus()
        heard = [AtomicCounter() for _ in range(3)]
        for counter in heard:
            bus.subscribe(lambda e, counter=counter: counter.add())
        bad_count = False
        bad_unsubscribe = False

        def publish_or_churn(t):
            nonlocal bad_count, bad_unsubscribe
            for _ in range(per_thread):
                if t < publishers:
                    if bus.publish("tick") < 3:
                        bad_count = True
                else:
                    sub_id = bus.subscribe(lambda e: None)
                    if not bus.unsubscribe(sub_id):
                        bad_unsubscribe = True

        with interleaved():
            run_together(publishers + churners, publish_or_churn)

        for counter in heard:
            self.assertEqual(counter.value, publishers * per_thread)
        self.assertFalse(bad_count, "a publish missed a permanent subscriber")
        self.assertFalse(bad_unsubscribe, "unsubscribe of a fresh id returned False")
        self.assertEqual(bus.subscriber_count(), 3)


if __name__ == "__main__":
    unittest.main()
