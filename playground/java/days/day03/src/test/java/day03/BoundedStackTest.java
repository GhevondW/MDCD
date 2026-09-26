package day03;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.Timeout;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.NoSuchElementException;
import java.util.OptionalInt;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.function.BooleanSupplier;

import static day03.RunTogether.runTogether;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

@Timeout(value = 60, threadMode = Timeout.ThreadMode.SEPARATE_THREAD)
class BoundedStackTest {

    // The single-threaded contract is the same as Day 2 -- a quick check that
    // it still holds.

    @Test
    void pushPopIsLifo() {
        BoundedStack s = new BoundedStack(3);
        s.push(1);
        s.push(2);
        s.push(3);
        assertTrue(s.isFull());
        assertEquals(3, s.top());
        assertEquals(3, s.pop());
        assertEquals(2, s.pop());
        assertEquals(1, s.size());
    }

    @Test
    void fullAndEmptyStillThrow() {
        BoundedStack s = new BoundedStack(1);
        assertTrue(s.isEmpty());
        assertThrows(NoSuchElementException.class, s::pop);
        assertThrows(NoSuchElementException.class, s::top);
        s.push(7);
        assertThrows(IllegalStateException.class, () -> s.push(8));
        assertEquals(1, s.size());
        assertEquals(7, s.top());
    }

    @Test
    void concurrentPushesStopExactlyAtCapacity() throws Exception {
        // 8 threads push 2,000 values into a stack that holds 1,000.
        // Exactly 1,000 pushes may succeed, and the stack must then hold
        // exactly the values whose push succeeded. The race only matters at
        // the moment the stack fills up, so play it 20 times.
        int threads = 8;
        int perThread = 250;
        for (int round = 0; round < 20; round++) {
            BoundedStack s = new BoundedStack(1000);
            List<List<Integer>> pushed = listPerThread(threads);
            runTogether(threads, t -> {
                for (int i = 0; i < perThread; i++) {
                    int v = t * 1000 + i;
                    try {
                        s.push(v);
                        pushed.get(t).add(v);
                    } catch (IllegalStateException full) {
                        // the stack was full: this push failed, as it should
                    }
                }
            });

            int r = round;
            List<Integer> expected = sortedUnion(pushed);
            assertEquals(1000, expected.size(), () -> "round " + r + ": pushes that succeeded");
            assertEquals(1000, s.size(), () -> "round " + r + ": size()");
            assertTrue(s.isFull(), () -> "round " + r + ": isFull()");

            List<Integer> held = new ArrayList<>();
            for (int i = 0; i < 1000 && !s.isEmpty(); i++) held.add(s.pop());
            assertTrue(s.isEmpty(), () -> "round " + r + ": isEmpty() after popping 1000 values");
            Collections.sort(held);
            assertEquals(expected, held, () -> "round " + r + ": the stack lost a value or made one up");
        }
    }

    @Test
    void concurrentPopsReturnEachValueOnce() throws Exception {
        int values = 20_000;
        int threads = 8;
        BoundedStack s = new BoundedStack(values);
        for (int v = 0; v < values; v++) s.push(v);

        // Pop until the stack says it is empty. (Stop early if more values come
        // out than ever went in -- the stack is making them up.)
        List<List<Integer>> popped = listPerThread(threads);
        AtomicInteger poppedCount = new AtomicInteger();
        runTogether(threads, t -> {
            while (poppedCount.get() <= values) {
                try {
                    popped.get(t).add(s.pop());
                    poppedCount.incrementAndGet();
                } catch (NoSuchElementException empty) {
                    return;
                }
            }
        });

        List<Integer> all = sortedUnion(popped);
        assertEquals(values, all.size(), "a value was popped twice, or lost");
        for (int v = 0; v < values; v++) {
            assertEquals(v, all.get(v), "a value was popped twice, or lost");
        }
        assertTrue(s.isEmpty());
    }

    @Test
    void concurrentPushAndPopConserveValues() throws Exception {
        // A small stack, 4 pushers and 4 poppers. Pushers retry while it is
        // full, poppers retry while it is empty. Every value pushed must come
        // out exactly once.
        int pushers = 4;
        int perPusher = 5000;
        int total = pushers * perPusher;
        BoundedStack s = new BoundedStack(16);
        AtomicInteger poppedCount = new AtomicInteger();
        List<List<Integer>> popped = listPerThread(pushers);
        long deadline = System.nanoTime() + TimeUnit.SECONDS.toNanos(20);
        BooleanSupplier pastDeadline = () -> System.nanoTime() - deadline > 0;

        runTogether(2 * pushers, t -> {
            if (t < pushers) {
                for (int i = 0; i < perPusher && !pastDeadline.getAsBoolean(); i++) {
                    while (true) {
                        try {
                            s.push(t * perPusher + i);
                            break;
                        } catch (IllegalStateException full) {
                            if (pastDeadline.getAsBoolean()) return;
                            Thread.yield();
                        }
                    }
                }
            } else {
                List<Integer> mine = popped.get(t - pushers);
                while (poppedCount.get() < total && !pastDeadline.getAsBoolean()) {
                    try {
                        mine.add(s.pop());
                        poppedCount.incrementAndGet();
                    } catch (NoSuchElementException empty) {
                        Thread.yield();
                    }
                }
            }
        });

        assertFalse(pastDeadline.getAsBoolean(), "gave up after 20 s -- values were lost");
        List<Integer> all = sortedUnion(popped);
        assertEquals(total, all.size(), "a value was popped twice, or lost");
        for (int v = 0; v < total; v++) {
            assertEquals(v, all.get(v), "a value was popped twice, or lost");
        }
        assertTrue(s.isEmpty());
    }

    // ---- tryPush / tryPop: check and act as one step ----

    @Test
    void tryPushAndTryPopFollowTheRules() {
        BoundedStack s = new BoundedStack(2);
        assertEquals(OptionalInt.empty(), s.tryPop(), "empty stack");
        assertTrue(s.tryPush(1));
        assertTrue(s.tryPush(2));
        assertFalse(s.tryPush(3), "full stack");
        assertEquals(2, s.size());
        assertEquals(OptionalInt.of(2), s.tryPop());
        assertEquals(OptionalInt.of(1), s.tryPop());
        assertEquals(OptionalInt.empty(), s.tryPop());
        assertTrue(s.isEmpty());
    }

    @Test
    void concurrentTryPopTakesEachValueOnce() throws Exception {
        // 8 threads empty the stack with tryPop -- the one-call version of
        // "if (!s.isEmpty()) { v = s.top(); s.pop(); }".
        int values = 20000;
        int threads = 8;
        BoundedStack s = new BoundedStack(values);
        for (int v = 0; v < values; v++) s.push(v);

        List<List<Integer>> popped = listPerThread(threads);
        AtomicInteger poppedCount = new AtomicInteger();
        runTogether(threads, t -> {
            while (poppedCount.get() <= values) {  // more than values: made up
                OptionalInt v = s.tryPop();
                if (v.isEmpty()) return;
                popped.get(t).add(v.getAsInt());
                poppedCount.incrementAndGet();
            }
        });

        List<Integer> all = sortedUnion(popped);
        assertEquals(values, all.size(), "a value was popped twice, or lost");
        for (int v = 0; v < values; v++) {
            assertEquals(v, all.get(v), "a value was popped twice, or lost");
        }
        assertTrue(s.isEmpty());
    }

    @Test
    void concurrentTryPushStopsExactlyAtCapacity() throws Exception {
        // Like concurrentPushesStopExactlyAtCapacity, with tryPush.
        int threads = 8;
        int perThread = 250;
        for (int round = 0; round < 20; round++) {
            BoundedStack s = new BoundedStack(1000);
            List<List<Integer>> pushed = listPerThread(threads);
            runTogether(threads, t -> {
                for (int i = 0; i < perThread; i++) {
                    int v = t * 1000 + i;
                    if (s.tryPush(v)) pushed.get(t).add(v);
                }
            });

            int r = round;
            List<Integer> expected = sortedUnion(pushed);
            assertEquals(1000, expected.size(), () -> "round " + r + ": tryPush calls that returned true");
            assertEquals(1000, s.size(), () -> "round " + r);
            List<Integer> held = new ArrayList<>();
            for (int i = 0; i < 1000; i++) {
                OptionalInt v = s.tryPop();
                if (v.isEmpty()) break;
                held.add(v.getAsInt());
            }
            Collections.sort(held);
            assertEquals(expected, held, () -> "round " + r + ": the stack lost a value or made one up");
        }
    }

    @Test
    void concurrentTryPushAndTryPopConserveValues() throws Exception {
        // A small stack, 4 pushers and 4 poppers, all using the try- calls.
        int pushers = 4;
        int perPusher = 5000;
        int total = pushers * perPusher;
        BoundedStack s = new BoundedStack(16);
        AtomicInteger poppedCount = new AtomicInteger();
        List<List<Integer>> popped = listPerThread(pushers);
        long deadline = System.nanoTime() + TimeUnit.SECONDS.toNanos(10);
        BooleanSupplier pastDeadline = () -> System.nanoTime() - deadline > 0;

        runTogether(2 * pushers, t -> {
            if (t < pushers) {
                for (int i = 0; i < perPusher && !pastDeadline.getAsBoolean(); i++) {
                    while (!s.tryPush(t * perPusher + i)) {
                        if (pastDeadline.getAsBoolean()) return;
                        Thread.yield();
                    }
                }
            } else {
                List<Integer> mine = popped.get(t - pushers);
                while (poppedCount.get() < total && !pastDeadline.getAsBoolean()) {
                    OptionalInt v = s.tryPop();
                    if (v.isPresent()) {
                        mine.add(v.getAsInt());
                        poppedCount.incrementAndGet();
                    } else {
                        Thread.yield();
                    }
                }
            }
        });

        assertFalse(pastDeadline.getAsBoolean(), "gave up after 10 s -- values were lost");
        List<Integer> all = sortedUnion(popped);
        assertEquals(total, all.size(), "a value was popped twice, or lost");
        for (int v = 0; v < total; v++) {
            assertEquals(v, all.get(v), "a value was popped twice, or lost");
        }
        assertTrue(s.isEmpty());
    }

    // One list per thread, so the threads never share a list.
    private static List<List<Integer>> listPerThread(int threads) {
        List<List<Integer>> lists = new ArrayList<>();
        for (int t = 0; t < threads; t++) lists.add(new ArrayList<>());
        return lists;
    }

    private static List<Integer> sortedUnion(List<List<Integer>> lists) {
        List<Integer> all = new ArrayList<>();
        for (List<Integer> l : lists) all.addAll(l);
        Collections.sort(all);
        return all;
    }
}
