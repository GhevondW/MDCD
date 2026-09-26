package day03;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.Timeout;

import java.util.concurrent.CountDownLatch;
import java.util.concurrent.Future;
import java.util.concurrent.TimeoutException;
import java.util.concurrent.atomic.AtomicBoolean;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.concurrent.atomic.AtomicIntegerArray;

import static day03.RunTogether.inBackground;
import static day03.RunTogether.runTogether;
import static java.util.concurrent.TimeUnit.SECONDS;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertSame;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

@Timeout(value = 60, threadMode = Timeout.ThreadMode.SEPARATE_THREAD)
class OnceCacheTest {

    @Test
    void missComputesAndReturnsTheValue() {
        OnceCache<Integer> c = new OnceCache<>();
        AtomicInteger calls = new AtomicInteger();
        assertEquals(42, c.get("a", () -> {
            calls.incrementAndGet();
            return 42;
        }));
        assertEquals(1, calls.get());
        assertEquals(1, c.size());
    }

    @Test
    void hitReturnsRememberedValueWithoutComputing() {
        OnceCache<Integer> c = new OnceCache<>();
        c.get("a", () -> 42);
        AtomicBoolean called = new AtomicBoolean(false);
        assertEquals(42, c.get("a", () -> {
            called.set(true);
            return 7;
        }));
        assertFalse(called.get());
    }

    @Test
    void differentKeysGetTheirOwnValues() {
        OnceCache<String> c = new OnceCache<>();
        assertEquals("apple", c.get("a", () -> "apple"));
        assertEquals("banana", c.get("b", () -> "banana"));
        assertEquals("apple", c.get("a", () -> "avocado"));
        assertEquals(2, c.size());
    }

    @Test
    void throwingComputeRemembersNothing() {
        OnceCache<Integer> c = new OnceCache<>();
        RuntimeException boom = new IllegalStateException("boom");
        RuntimeException thrown = assertThrows(RuntimeException.class,
                () -> c.get("a", () -> { throw boom; }));
        assertSame(boom, thrown, "compute's exception must reach the caller unchanged");
        assertEquals(0, c.size());
        AtomicInteger calls = new AtomicInteger();
        assertEquals(5, c.get("a", () -> {
            calls.incrementAndGet();
            return 5;
        }));
        assertEquals(1, calls.get(), "after a failure, the next get() must compute again");
        assertEquals(1, c.size());
    }

    @Test
    void concurrentGetsOfOneKeyComputeOnce() throws Exception {
        // 16 threads miss the same key at the same moment; compute is slow,
        // so a check-then-act gap is wide open.
        int threads = 16;
        OnceCache<Integer> c = new OnceCache<>();
        AtomicInteger calls = new AtomicInteger();
        Integer[] results = new Integer[threads];
        runTogether(threads, t -> {
            results[t] = c.get("config", () -> {
                calls.incrementAndGet();
                sleepMillis(50);
                return 99;
            });
        });

        assertEquals(1, calls.get(), "compute ran more than once for one key");
        for (Integer r : results) assertEquals(99, r);
        assertEquals(1, c.size());
    }

    @Test
    void concurrentGetsOfManyKeysComputeEachOnce() throws Exception {
        int threads = 8;
        int keys = 200;
        OnceCache<Integer> c = new OnceCache<>();
        AtomicIntegerArray calls = new AtomicIntegerArray(keys);
        runTogether(threads, t -> {
            for (int i = 0; i < keys; i++) {
                int k = (i + t * 25) % keys; // each thread starts at a different key
                Integer v = c.get(String.valueOf(k), () -> {
                    calls.incrementAndGet(k);
                    sleepMillis(1);
                    return k * 10;
                });
                assertEquals(k * 10, v, () -> "value for key " + k);
            }
        });

        for (int k = 0; k < keys; k++) {
            int key = k;
            assertEquals(1, calls.get(k), () -> "compute calls for key " + key + " (must be exactly 1)");
        }
        assertEquals(keys, c.size());
    }

    @Test
    void differentKeysDoNotWaitForEachOther() throws Exception {
        // compute for "a" and compute for "b" each wait until the other one
        // has started. That only works if both can run at the same time; if a
        // lock is held while compute runs, one of them gives up after 2 s.
        assertBothComputesRunAtOnce("a", "b");
    }

    @Test
    void keysWithTheSameHashCodeDoNotWaitForEachOther() throws Exception {
        // The same check with two keys that have the same hashCode(), so every
        // hash map puts them in the same bucket. ConcurrentHashMap's
        // computeIfAbsent locks that bucket while its function runs: calling
        // compute inside it makes these two keys wait for each other.
        assertEquals("Aa".hashCode(), "BB".hashCode());
        assertBothComputesRunAtOnce("Aa", "BB");
    }

    @Test
    void sizeCountsOnlyFinishedValues() throws Exception {
        OnceCache<Integer> c = new OnceCache<>();
        CountDownLatch started = new CountDownLatch(1);
        CountDownLatch release = new CountDownLatch(1);
        Future<Integer> slow = inBackground(() -> c.get("slow", () -> {
            started.countDown();
            awaitQuietly(release, 30);
            return 1;
        }));
        boolean didStart = started.await(5, SECONDS);
        // Ask size() on another thread, so a size() that waits for the running
        // compute fails this test instead of hanging it.
        Future<Integer> sizeWhileRunning = inBackground(c::size);
        boolean sizeAnswered;
        try {
            sizeWhileRunning.get(2, SECONDS);
            sizeAnswered = true;
        } catch (TimeoutException e) {
            sizeAnswered = false;
        }
        release.countDown();
        slow.get();

        assertTrue(didStart, "compute was never called");
        assertTrue(sizeAnswered, "size() waited for a running compute");
        assertEquals(0, sizeWhileRunning.get(), "a computation still running counted as remembered");
        assertEquals(1, c.size());
    }

    // Thread 0 gets key1 and thread 1 gets key2. Each compute says it has
    // started, then waits (up to 2 s) for the other compute to start too.
    private static void assertBothComputesRunAtOnce(String key1, String key2) throws Exception {
        OnceCache<String> c = new OnceCache<>();
        String[] keys = {key1, key2};
        CountDownLatch[] started = {new CountDownLatch(1), new CountDownLatch(1)};
        String[] result = new String[2];
        runTogether(2, t -> {
            result[t] = c.get(keys[t], () -> {
                started[t].countDown();
                return awaitQuietly(started[1 - t], 2) ? "ok" : "timed out";
            });
        });

        assertEquals("ok", result[0], "compute for \"" + key1 + "\" waited for compute for \"" + key2 + "\"");
        assertEquals("ok", result[1], "compute for \"" + key2 + "\" waited for compute for \"" + key1 + "\"");
    }

    private static boolean awaitQuietly(CountDownLatch latch, long seconds) {
        try {
            return latch.await(seconds, SECONDS);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            return false;
        }
    }

    private static void sleepMillis(long millis) {
        try {
            Thread.sleep(millis);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }
    }
}
