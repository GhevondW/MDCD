package day03;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.Timeout;

import java.util.ArrayList;
import java.util.HashSet;
import java.util.List;
import java.util.Set;
import java.util.concurrent.Future;
import java.util.concurrent.TimeoutException;
import java.util.concurrent.atomic.AtomicBoolean;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.concurrent.atomic.AtomicLong;
import java.util.concurrent.atomic.AtomicReference;

import static day03.RunTogether.inBackground;
import static day03.RunTogether.runTogether;
import static java.util.concurrent.TimeUnit.SECONDS;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

@Timeout(value = 60, threadMode = Timeout.ThreadMode.SEPARATE_THREAD)
class EventBusTest {

    @Test
    void publishCallsSubscribersInOrder() {
        EventBus bus = new EventBus();
        List<String> log = new ArrayList<>();
        bus.subscribe(e -> log.add("A:" + e));
        bus.subscribe(e -> log.add("B:" + e));
        bus.subscribe(e -> log.add("C:" + e));
        assertEquals(3, bus.publish("x"));
        assertEquals(List.of("A:x", "B:x", "C:x"), log);
    }

    @Test
    void publishWithNoSubscribersCallsNobody() {
        EventBus bus = new EventBus();
        assertEquals(0, bus.publish("x"));
        assertEquals(0, bus.subscriberCount());
    }

    @Test
    void idsArePositiveAndUnique() {
        EventBus bus = new EventBus();
        Set<Long> ids = new HashSet<>();
        for (int i = 0; i < 100; i++) {
            long id = bus.subscribe(e -> {});
            assertTrue(id > 0, "ids must be > 0, got " + id);
            ids.add(id);
        }
        assertEquals(100, ids.size(), "different ids among 100 subscribes");
        assertEquals(100, bus.subscriberCount());
    }

    @Test
    void unsubscribeStopsDelivery() {
        EventBus bus = new EventBus();
        List<String> log = new ArrayList<>();
        bus.subscribe(e -> log.add("A"));
        long b = bus.subscribe(e -> log.add("B"));
        bus.subscribe(e -> log.add("C"));

        assertTrue(bus.unsubscribe(b));
        assertFalse(bus.unsubscribe(b), "already removed");
        assertFalse(bus.unsubscribe(12345), "never existed");
        assertEquals(2, bus.subscriberCount());
        assertEquals(2, bus.publish("x"));
        assertEquals(List.of("A", "C"), log);
    }

    @Test
    void orderIsKeptAfterUnsubscribeAndSubscribe() {
        EventBus bus = new EventBus();
        List<String> log = new ArrayList<>();
        bus.subscribe(e -> log.add("A"));
        long b = bus.subscribe(e -> log.add("B"));
        bus.subscribe(e -> log.add("C"));
        bus.unsubscribe(b);
        bus.subscribe(e -> log.add("D"));
        bus.publish("x");
        assertEquals(List.of("A", "C", "D"), log);
    }

    @Test
    void subscribeInsideACallbackCountsFromTheNextPublish() {
        EventBus bus = new EventBus();
        List<String> log = new ArrayList<>();
        AtomicBoolean added = new AtomicBoolean(false);
        bus.subscribe(e -> {
            log.add("A");
            if (!added.getAndSet(true)) {
                bus.subscribe(e2 -> log.add("B"));
            }
        });

        assertEquals(1, bus.publish("first"), "B joined mid-publish: not called by it");
        assertEquals(List.of("A"), log);
        log.clear();
        assertEquals(2, bus.publish("second"));
        assertEquals(List.of("A", "B"), log);
    }

    @Test
    void unsubscribeInsideACallbackCountsFromTheNextPublish() {
        EventBus bus = new EventBus();
        List<String> log = new ArrayList<>();
        AtomicLong b = new AtomicLong();
        AtomicBoolean removed = new AtomicBoolean(false);
        AtomicBoolean unsubscribeReturned = new AtomicBoolean(false);
        bus.subscribe(e -> {
            log.add("A");
            if (!removed.getAndSet(true)) {
                unsubscribeReturned.set(bus.unsubscribe(b.get()));
            }
        });
        b.set(bus.subscribe(e -> log.add("B")));

        assertEquals(2, bus.publish("first"), "B left mid-publish: still called by it");
        assertEquals(List.of("A", "B"), log);
        assertTrue(unsubscribeReturned.get(), "unsubscribe(b) inside a callback returned false");
        log.clear();
        assertEquals(1, bus.publish("second"));
        assertEquals(List.of("A"), log);
    }

    @Test
    void callbackCanUnsubscribeItself() {
        EventBus bus = new EventBus();
        AtomicInteger aCalls = new AtomicInteger();
        AtomicLong a = new AtomicLong();
        a.set(bus.subscribe(e -> {
            aCalls.incrementAndGet();
            bus.unsubscribe(a.get());
        }));
        bus.subscribe(e -> {});

        assertEquals(2, bus.publish("first"));
        assertEquals(1, bus.publish("second"));
        assertEquals(1, aCalls.get());
    }

    @Test
    void callbackCanPublish() {
        EventBus bus = new EventBus();
        List<String> bLog = new ArrayList<>();
        bus.subscribe(e -> {
            if (e.equals("outer")) bus.publish("inner");
        });
        bus.subscribe(bLog::add);

        assertEquals(2, bus.publish("outer"));
        // A handles "outer" by publishing "inner", which reaches B first.
        assertEquals(List.of("inner", "outer"), bLog);
    }

    @Test
    void otherThreadsCanUseTheBusWhileACallbackRuns() throws Exception {
        // The callback waits (up to 2 s) for another thread to subscribe. If
        // publish() holds a lock while calling back, that thread is stuck.
        EventBus bus = new EventBus();
        AtomicReference<String> outcome = new AtomicReference<>();
        AtomicReference<Future<Long>> helper = new AtomicReference<>();
        bus.subscribe(e -> {
            Future<Long> h = inBackground(() -> bus.subscribe(e2 -> {}));
            helper.set(h);
            try {
                h.get(2, SECONDS);
                outcome.set("ok");
            } catch (TimeoutException stuck) {
                outcome.set("blocked");
            } catch (Exception other) {
                outcome.set("helper failed: " + other);
            }
        });

        bus.publish("x");
        assertNotNull(helper.get(), "publish() did not call the callback");
        helper.get().get(); // wait for the helper thread to finish
        assertEquals("ok", outcome.get(), "another thread could not subscribe while a callback ran");
        assertEquals(2, bus.subscriberCount());
    }

    @Test
    void concurrentSubscribesGetUniqueIds() throws Exception {
        int threads = 8;
        int perThread = 1000;
        EventBus bus = new EventBus();
        long[][] ids = new long[threads][perThread];
        runTogether(threads, t -> {
            for (int i = 0; i < perThread; i++) {
                ids[t][i] = bus.subscribe(e -> {});
            }
        });

        Set<Long> all = new HashSet<>();
        for (long[] mine : ids) {
            for (long id : mine) all.add(id);
        }
        assertEquals(threads * perThread, all.size(), "different ids handed out");
        assertEquals(threads * perThread, bus.subscriberCount());
    }

    @Test
    void concurrentPublishAndChurn() throws Exception {
        // 3 subscribers stay for the whole test and must hear every event
        // exactly once, while other threads keep subscribing and
        // unsubscribing around them.
        int publishers = 4;
        int churners = 4;
        int perThread = 2000;
        EventBus bus = new EventBus();
        AtomicInteger[] heard = {new AtomicInteger(), new AtomicInteger(), new AtomicInteger()};
        for (AtomicInteger h : heard) {
            bus.subscribe(e -> h.incrementAndGet());
        }
        AtomicBoolean badCount = new AtomicBoolean(false);
        AtomicBoolean badUnsubscribe = new AtomicBoolean(false);

        runTogether(publishers + churners, t -> {
            for (int i = 0; i < perThread; i++) {
                if (t < publishers) {
                    if (bus.publish("tick") < 3) badCount.set(true);
                } else {
                    long id = bus.subscribe(e -> {});
                    if (!bus.unsubscribe(id)) badUnsubscribe.set(true);
                }
            }
        });

        for (AtomicInteger h : heard) {
            assertEquals(publishers * perThread, h.get(), "events heard by a permanent subscriber");
        }
        assertFalse(badCount.get(), "a publish missed a permanent subscriber");
        assertFalse(badUnsubscribe.get(), "unsubscribe of a fresh id returned false");
        assertEquals(3, bus.subscriberCount());
    }
}
