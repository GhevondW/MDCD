package day04;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.Timeout;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.Optional;
import java.util.concurrent.Future;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;
import java.util.concurrent.atomic.AtomicBoolean;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.concurrent.locks.LockSupport;
import java.util.function.IntConsumer;

import static day04.RunTogether.inBackground;
import static day04.RunTogether.runTogether;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

@Timeout(value = 60, threadMode = Timeout.ThreadMode.SEPARATE_THREAD)
class BoundedBlockingQueueTest {

    /** True if the call is still running after `millis` ms (it is blocked). */
    private static boolean stillBlocked(Future<?> call, long millis) throws Exception {
        try {
            call.get(millis, TimeUnit.MILLISECONDS);
            return false;
        } catch (TimeoutException e) {
            return true;
        }
    }

    /** The call's result; fails the test if it does not return in 5 s. */
    private static <T> T resultOf(Future<T> call, String what) throws Exception {
        try {
            return call.get(5, TimeUnit.SECONDS);
        } catch (TimeoutException e) {
            throw new AssertionError(what + " did not return within 5 s");
        }
    }

    /** Pops until pop() says "nothing" -- which is only allowed once the queue is closed. */
    private static void drain(BoundedBlockingQueue<Integer> q, IntConsumer each) throws InterruptedException {
        for (Optional<Integer> v = q.pop(); v.isPresent(); v = q.pop()) each.accept(v.get());
        assertTrue(q.isClosed(), "pop() returned nothing although the queue is still open");
    }

    // ---- the single-threaded contract ----

    @Test
    void itemsComeOutInTheOrderTheyWentIn() throws Exception {
        BoundedBlockingQueue<Integer> q = new BoundedBlockingQueue<>(3);
        assertEquals(3, q.capacity());
        assertEquals(0, q.size());
        assertTrue(q.push(1));
        assertTrue(q.push(2));
        assertTrue(q.push(3));
        assertEquals(3, q.size());
        assertEquals(Optional.of(1), q.pop());
        assertTrue(q.push(4));
        assertEquals(Optional.of(2), q.pop());
        assertEquals(Optional.of(3), q.pop());
        assertEquals(Optional.of(4), q.pop());
        assertEquals(0, q.size());
        assertFalse(q.isClosed());
    }

    @Test
    void capacityMustBePositive() {
        assertThrows(IllegalArgumentException.class, () -> new BoundedBlockingQueue<Integer>(0));
        assertThrows(IllegalArgumentException.class, () -> new BoundedBlockingQueue<Integer>(-1));
    }

    @Test
    void tryPushAndTryPopNeverWait() {
        BoundedBlockingQueue<String> q = new BoundedBlockingQueue<>(2);
        assertEquals(Optional.empty(), q.tryPop(), "empty queue");
        assertTrue(q.tryPush("a"));
        assertTrue(q.tryPush("b"));
        assertFalse(q.tryPush("c"), "full queue");
        assertEquals(2, q.size());
        assertEquals(Optional.of("a"), q.tryPop());
        assertTrue(q.tryPush("c"));
        assertEquals(Optional.of("b"), q.tryPop());
        assertEquals(Optional.of("c"), q.tryPop());
        assertEquals(Optional.empty(), q.tryPop());
    }

    // ---- blocking ----

    @Test
    void pushWaitsWhileTheQueueIsFull() throws Exception {
        BoundedBlockingQueue<Integer> q = new BoundedBlockingQueue<>(1);
        assertTrue(q.push(1));
        Future<Boolean> second = inBackground(() -> q.push(2));
        assertTrue(stillBlocked(second, 300), "push on a full queue must wait");
        assertEquals(1, q.size());

        assertEquals(Optional.of(1), q.pop());          // makes room
        assertTrue(resultOf(second, "the waiting push"), "the waiting push should succeed");
        assertEquals(Optional.of(2), q.pop());
    }

    @Test
    void popWaitsWhileTheQueueIsEmpty() throws Exception {
        BoundedBlockingQueue<Integer> q = new BoundedBlockingQueue<>(2);
        Future<Optional<Integer>> taken = inBackground(q::pop);
        assertTrue(stillBlocked(taken, 300), "pop on an empty queue must wait");

        assertTrue(q.push(7));
        assertEquals(Optional.of(7), resultOf(taken, "the waiting pop"));
    }

    @Test
    void everyWaitingPopperGetsExactlyOneItem() throws Exception {
        // Four consumers asleep; four items arrive, one at a time.
        BoundedBlockingQueue<Integer> q = new BoundedBlockingQueue<>(8);
        List<Future<Optional<Integer>>> takers = new ArrayList<>();
        for (int i = 0; i < 4; i++) takers.add(inBackground(q::pop));
        Thread.sleep(200);
        for (int v = 1; v <= 4; v++) {
            assertTrue(q.push(v));
            Thread.sleep(20);
        }
        List<Integer> got = new ArrayList<>();
        for (Future<Optional<Integer>> t : takers) {
            got.add(resultOf(t, "a waiting pop").orElseThrow());
        }
        Collections.sort(got);
        assertEquals(List.of(1, 2, 3, 4), got);
    }

    // ---- close ----

    @Test
    void closeWakesEverySleepingConsumer() throws Exception {
        BoundedBlockingQueue<Integer> q = new BoundedBlockingQueue<>(2);
        List<Future<Optional<Integer>>> takers = new ArrayList<>();
        for (int i = 0; i < 4; i++) takers.add(inBackground(q::pop));
        Thread.sleep(200);
        q.close();
        assertTrue(q.isClosed());
        for (Future<Optional<Integer>> t : takers) {
            assertEquals(Optional.empty(), resultOf(t, "a pop that was asleep when close() ran"));
        }
    }

    @Test
    void closeWakesEverySleepingProducer() throws Exception {
        BoundedBlockingQueue<Integer> q = new BoundedBlockingQueue<>(1);
        assertTrue(q.push(0));
        List<Future<Boolean>> pushers = new ArrayList<>();
        for (int i = 1; i <= 4; i++) {
            int v = i;
            pushers.add(inBackground(() -> q.push(v)));
        }
        Thread.sleep(200);
        q.close();
        for (Future<Boolean> p : pushers) {
            assertFalse(resultOf(p, "a push that was asleep when close() ran"),
                    "a push that was waiting when the queue closed must return false");
        }
        assertEquals(1, q.size(), "nothing may be added after close()");
        assertEquals(Optional.of(0), q.pop());
        assertEquals(Optional.empty(), q.pop());
    }

    @Test
    void afterCloseTheItemsInsideCanStillBeDrained() throws Exception {
        BoundedBlockingQueue<Integer> q = new BoundedBlockingQueue<>(5);
        q.push(1);
        q.push(2);
        q.close();
        q.close();                                       // again: nothing happens
        assertTrue(q.isClosed());
        assertFalse(q.push(3));
        assertFalse(q.tryPush(3));
        assertFalse(q.pushFor(3, 10, TimeUnit.MILLISECONDS));
        assertEquals(2, q.size());
        assertEquals(Optional.of(1), q.pop());
        assertEquals(Optional.of(2), q.tryPop());
        assertEquals(Optional.empty(), q.pop(), "closed and empty: pop returns at once");
        assertEquals(Optional.empty(), q.tryPop());
        assertEquals(Optional.empty(), q.popFor(10, TimeUnit.SECONDS));
    }

    @Test
    void aConsumerThatWasAsleepGetsTheLastItemsThenEmpty() throws Exception {
        // Close while the queue has items and consumers are about to drain.
        BoundedBlockingQueue<Integer> q = new BoundedBlockingQueue<>(10);
        for (int i = 0; i < 10; i++) q.push(i);
        q.close();
        AtomicInteger taken = new AtomicInteger();
        runTogether(4, t -> {
            drain(q, v -> taken.incrementAndGet());
        });
        assertEquals(10, taken.get());
    }

    // ---- timeouts ----

    @Test
    void popForGivesUpAfterTheTimeout() throws Exception {
        BoundedBlockingQueue<Integer> q = new BoundedBlockingQueue<>(2);
        long start = System.nanoTime();
        assertEquals(Optional.empty(), q.popFor(100, TimeUnit.MILLISECONDS));
        long ms = TimeUnit.NANOSECONDS.toMillis(System.nanoTime() - start);
        assertTrue(ms >= 80, "returned after only " + ms + " ms: it must wait for the timeout");
        assertTrue(ms < 3000, "returned after " + ms + " ms: far too late for a 100 ms timeout");
    }

    @Test
    void pushForGivesUpAfterTheTimeout() throws Exception {
        BoundedBlockingQueue<Integer> q = new BoundedBlockingQueue<>(1);
        q.push(1);
        long start = System.nanoTime();
        assertFalse(q.pushFor(2, 100, TimeUnit.MILLISECONDS));
        long ms = TimeUnit.NANOSECONDS.toMillis(System.nanoTime() - start);
        assertTrue(ms >= 80, "returned after only " + ms + " ms: it must wait for the timeout");
        assertTrue(ms < 3000, "returned after " + ms + " ms: far too late for a 100 ms timeout");
        assertEquals(1, q.size());
    }

    @Test
    void popForReturnsAsSoonAsAnItemArrives() throws Exception {
        BoundedBlockingQueue<Integer> q = new BoundedBlockingQueue<>(2);
        Future<Optional<Integer>> taken = inBackground(() -> q.popFor(30, TimeUnit.SECONDS));
        assertTrue(stillBlocked(taken, 200));
        q.push(5);
        assertEquals(Optional.of(5), resultOf(taken, "popFor after an item arrived"));
    }

    @Test
    void pushForReturnsAsSoonAsThereIsRoom() throws Exception {
        BoundedBlockingQueue<Integer> q = new BoundedBlockingQueue<>(1);
        q.push(1);
        Future<Boolean> pushed = inBackground(() -> q.pushFor(2, 30, TimeUnit.SECONDS));
        assertTrue(stillBlocked(pushed, 200));
        assertEquals(Optional.of(1), q.pop());
        assertTrue(resultOf(pushed, "pushFor after room appeared"));
        assertEquals(Optional.of(2), q.pop());
    }

    @Test
    void timedCallsWakeWhenTheQueueCloses() throws Exception {
        BoundedBlockingQueue<Integer> full = new BoundedBlockingQueue<>(1);
        full.push(1);
        BoundedBlockingQueue<Integer> empty = new BoundedBlockingQueue<>(1);
        Future<Boolean> pushing = inBackground(() -> full.pushFor(2, 30, TimeUnit.SECONDS));
        Future<Optional<Integer>> popping = inBackground(() -> empty.popFor(30, TimeUnit.SECONDS));
        Thread.sleep(200);
        full.close();
        empty.close();
        assertFalse(resultOf(pushing, "pushFor after close()"));
        assertEquals(Optional.empty(), resultOf(popping, "popFor after close()"));
    }

    @Test
    void theTimeoutIsADeadlineNotARestartableWait() throws Exception {
        // Two consumers wait up to 1 s. At 0.5 s one item arrives: one of
        // them takes it. The other may be woken up -- it must find nothing,
        // go back to sleep, and still give up at 1 s, not at 1.5 s.
        BoundedBlockingQueue<Integer> q = new BoundedBlockingQueue<>(4);
        long start = System.nanoTime();
        List<Future<Optional<Integer>>> results = new ArrayList<>();
        for (int i = 0; i < 2; i++) {
            Future<Optional<Integer>> r = inBackground(() -> q.popFor(1000, TimeUnit.MILLISECONDS));
            results.add(r);
        }
        Thread.sleep(500);
        q.push(1);
        int got = 0;
        for (Future<Optional<Integer>> r : results) {
            if (resultOf(r, "popFor").isPresent()) got++;
        }
        long ms = TimeUnit.NANOSECONDS.toMillis(System.nanoTime() - start);
        assertEquals(1, got, "exactly one consumer gets the one item");
        assertTrue(ms >= 900, "both calls ended after " + ms + " ms; the empty-handed one must wait the full 1000 ms");
        assertTrue(ms < 1400, "both calls ended after " + ms
                + " ms: the empty-handed consumer restarted its timeout when it was woken up");
    }

    // ---- many threads ----

    @Test
    void manyProducersAndManyConsumersLoseNothing() throws Exception {
        int producers = 4;
        int consumers = 4;
        int perProducer = 5000;
        int capacity = 16;
        BoundedBlockingQueue<Integer> q = new BoundedBlockingQueue<>(capacity);
        AtomicInteger producersLeft = new AtomicInteger(producers);
        AtomicBoolean overfull = new AtomicBoolean();
        List<List<Integer>> taken = new ArrayList<>();
        for (int i = 0; i < consumers; i++) taken.add(new ArrayList<>());

        runTogether(producers + consumers, t -> {
            if (t < producers) {
                for (int i = 0; i < perProducer; i++) {
                    assertTrue(q.push(t * perProducer + i), "push failed on an open queue");
                    if (q.size() > capacity) overfull.set(true);
                }
                if (producersLeft.decrementAndGet() == 0) q.close();
            } else {
                List<Integer> mine = taken.get(t - producers);
                drain(q, mine::add);
            }
        });

        assertFalse(overfull.get(), "size() went above the capacity");
        // Everything arrives exactly once...
        List<Integer> all = new ArrayList<>();
        for (List<Integer> c : taken) all.addAll(c);
        Collections.sort(all);
        assertEquals(producers * perProducer, all.size(), "an item was lost or returned twice");
        for (int v = 0; v < all.size(); v++) assertEquals(v, all.get(v), "an item was lost or returned twice");
        // ...and one producer's items reach any one consumer in order (FIFO).
        for (List<Integer> c : taken) {
            int[] last = new int[producers];
            java.util.Arrays.fill(last, -1);
            for (int v : c) {
                int p = v / perProducer;
                assertTrue(v > last[p], "items of one producer came out of order");
                last[p] = v;
            }
        }
        assertEquals(0, q.size());
    }

    @Test
    void pingPongThroughAQueueOfOneNeverGetsStuck() throws Exception {
        // Capacity 1 makes the producer and the consumer wake each other on
        // every single item. A missed wake-up shows up as a hang.
        int items = 20_000;
        BoundedBlockingQueue<Integer> q = new BoundedBlockingQueue<>(1);
        AtomicInteger sum = new AtomicInteger();
        runTogether(2, t -> {
            if (t == 0) {
                for (int i = 1; i <= items; i++) assertTrue(q.push(i));
                q.close();
            } else {
                drain(q, sum::addAndGet);
            }
        });
        assertEquals(items * (items + 1) / 2, sum.get());
    }

    @Test
    void severalProducersAndConsumersOnATinyQueueNeverGetStuck() throws Exception {
        int perProducer = 5000;
        BoundedBlockingQueue<Integer> q = new BoundedBlockingQueue<>(2);
        AtomicInteger count = new AtomicInteger();
        AtomicInteger producersLeft = new AtomicInteger(3);
        runTogether(6, t -> {
            if (t < 3) {
                for (int i = 0; i < perProducer; i++) assertTrue(q.push(i));
                if (producersLeft.decrementAndGet() == 0) q.close();
            } else {
                drain(q, v -> count.incrementAndGet());
            }
        });
        assertEquals(3 * perProducer, count.get());
    }

    @Test
    void manyConsumersCompeteForASlowlyFedQueue() throws Exception {
        // Eight consumers, one producer that adds one item at a time. A
        // consumer that was woken up for an item may find another consumer
        // has already taken it -- it must go back to sleep, not take
        // "nothing".
        int consumers = 8;
        int items = 10_000;
        BoundedBlockingQueue<Integer> q = new BoundedBlockingQueue<>(4);
        AtomicInteger sum = new AtomicInteger();
        AtomicInteger count = new AtomicInteger();
        runTogether(consumers + 1, t -> {
            if (t == 0) {
                for (int i = 1; i <= items; i++) {
                    assertTrue(q.push(i));
                    // Two items in quick succession, then a pause: the consumers
                    // fall asleep, and the second item can be taken by a consumer
                    // that is awake before the one that was woken up for it runs.
                    if (i % 2 == 0) LockSupport.parkNanos(30_000);
                }
                q.close();
            } else {
                drain(q, v -> {
                    sum.addAndGet(v);
                    count.incrementAndGet();
                });
            }
        });
        assertEquals(items, count.get());
        assertEquals(items * (items + 1) / 2, sum.get());
    }

    @Test
    void manyProducersCompeteForASlowlyDrainedQueue() throws Exception {
        // The mirror image: eight producers, one consumer that frees one
        // slot at a time.
        int producers = 8;
        int each = 2500;
        BoundedBlockingQueue<Integer> q = new BoundedBlockingQueue<>(2);
        AtomicInteger producersLeft = new AtomicInteger(producers);
        AtomicInteger count = new AtomicInteger();
        AtomicBoolean overfull = new AtomicBoolean();
        runTogether(producers + 1, t -> {
            if (t == 0) {
                drain(q, v -> {
                    count.incrementAndGet();
                    if (count.get() % 4 == 0) Thread.yield();
                });
            } else {
                for (int i = 0; i < each; i++) {
                    assertTrue(q.push(i));
                    if (q.size() > 2) overfull.set(true);
                }
                if (producersLeft.decrementAndGet() == 0) q.close();
            }
        });
        assertFalse(overfull.get(), "size() went above the capacity");
        assertEquals(producers * each, count.get());
    }

    @Test
    void aConsumerWokenForAnItemSomeoneElseTookGoesBackToSleep() throws Exception {
        // Two consumers sleep in pop(). A third thread does nothing but
        // tryPop() in a loop, so it usually grabs each new item before the
        // sleeper that was woken up for it gets to run. That sleeper must
        // notice the queue is empty again and wait some more -- not take
        // "nothing".
        int items = 3000;
        BoundedBlockingQueue<Integer> q = new BoundedBlockingQueue<>(8);
        AtomicInteger count = new AtomicInteger();
        runTogether(4, t -> {
            if (t == 0) {
                for (int i = 0; i < items; i++) {
                    assertTrue(q.push(i));
                    LockSupport.parkNanos(50_000);
                }
                q.close();
            } else if (t == 1) {
                while (!q.isClosed() || q.size() > 0) {
                    if (q.tryPop().isPresent()) count.incrementAndGet();
                }
            } else {
                drain(q, v -> count.incrementAndGet());
            }
        });
        assertEquals(items, count.get());
    }

    @Test
    void popForGoesBackToSleepWhenItWasWokenForNothing() throws Exception {
        // A consumer waits up to 1 s. Three times another thread pushes an
        // item and takes it right back, so the consumer is woken up -- and
        // finds the queue empty again. It must keep waiting for the rest of
        // its second. (If it does win an item, that is fine too.)
        BoundedBlockingQueue<Integer> q = new BoundedBlockingQueue<>(4);
        long start = System.nanoTime();
        Future<Optional<Integer>> consumer = inBackground(() -> q.popFor(1000, TimeUnit.MILLISECONDS));
        for (int i = 0; i < 3; i++) {
            Thread.sleep(150);
            q.push(i);
            q.tryPop();
        }
        Optional<Integer> got = resultOf(consumer, "popFor");
        long ms = TimeUnit.NANOSECONDS.toMillis(System.nanoTime() - start);
        if (got.isEmpty()) {
            assertTrue(ms >= 900, "popFor gave up after " + ms + " ms: it was woken up for nothing and stopped waiting");
        }
    }
}
