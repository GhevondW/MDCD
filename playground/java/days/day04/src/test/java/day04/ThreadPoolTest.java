package day04;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.Timeout;

import java.util.ArrayList;
import java.util.Collections;
import java.util.HashSet;
import java.util.List;
import java.util.Set;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.Future;
import java.util.concurrent.RejectedExecutionException;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;
import java.util.concurrent.atomic.AtomicInteger;

import static day04.RunTogether.inBackground;
import static day04.RunTogether.runTogether;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertInstanceOf;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

@Timeout(value = 60, threadMode = Timeout.ThreadMode.SEPARATE_THREAD)
class ThreadPoolTest {

    private static final long WAIT_SECONDS = 5;

    private static boolean await(CountDownLatch latch) throws InterruptedException {
        return latch.await(WAIT_SECONDS, TimeUnit.SECONDS);
    }

    private static <T> T get(Future<T> f) throws Exception {
        try {
            return f.get(WAIT_SECONDS, TimeUnit.SECONDS);
        } catch (TimeoutException e) {
            throw new AssertionError("the task did not finish within " + WAIT_SECONDS + " s -- was it ever run?");
        }
    }

    private static boolean stillBlocked(Future<?> call, long millis) throws Exception {
        try {
            call.get(millis, TimeUnit.MILLISECONDS);
            return false;
        } catch (TimeoutException e) {
            return true;
        }
    }

    // ---- basics ----

    @Test
    void submitReturnsTheResult() throws Exception {
        ThreadPool pool = new ThreadPool(2);
        try {
            Future<Integer> f = pool.submit(() -> 6 * 7);
            assertEquals(42, get(f));
        } finally {
            pool.shutdown();
        }
    }

    @Test
    void anExceptionTravelsThroughTheFuture() throws Exception {
        ThreadPool pool = new ThreadPool(2);
        try {
            Future<Integer> bad = pool.submit(() -> {
                throw new IllegalStateException("boom");
            });
            ExecutionException e = assertThrows(ExecutionException.class, () -> get(bad));
            assertInstanceOf(IllegalStateException.class, e.getCause());
            assertEquals(5, get(pool.submit(() -> 5)), "the pool still works");
        } finally {
            pool.shutdown();
        }
    }

    @Test
    void theArgumentsAreChecked() {
        assertThrows(IllegalArgumentException.class, () -> new ThreadPool(0));
        assertThrows(IllegalArgumentException.class, () -> new ThreadPool(-2, 5));
        assertThrows(IllegalArgumentException.class, () -> new ThreadPool(2, -1));
    }

    @Test
    void everyTaskRunsExactlyOnce() throws Exception {
        int tasks = 10_000;
        ThreadPool pool = new ThreadPool(4);
        AtomicInteger[] runs = new AtomicInteger[tasks];
        for (int i = 0; i < tasks; i++) runs[i] = new AtomicInteger();
        CountDownLatch all = new CountDownLatch(tasks);
        try {
            for (int i = 0; i < tasks; i++) {
                int id = i;
                pool.execute(() -> {
                    runs[id].incrementAndGet();
                    all.countDown();
                });
            }
            assertTrue(await(all), "not every task ran");
        } finally {
            pool.shutdown();
        }
        for (int i = 0; i < tasks; i++) assertEquals(1, runs[i].get(), "task " + i + " ran this many times");
    }

    @Test
    void manyThreadsSubmitAtTheSameTime() throws Exception {
        int submitters = 8;
        int perSubmitter = 2000;
        ThreadPool pool = new ThreadPool(4, 64);
        AtomicInteger ran = new AtomicInteger();
        try {
            runTogether(submitters, t -> {
                for (int i = 0; i < perSubmitter; i++) pool.execute(ran::incrementAndGet);
            });
        } finally {
            pool.shutdown();
        }
        assertEquals(submitters * perSubmitter, ran.get());
    }

    // ---- parallelism ----

    @Test
    void tasksRunAtTheSameTime() throws Exception {
        // Four tasks that each wait for the other three: they can only all
        // finish if four workers run them at the same moment.
        int workers = 4;
        ThreadPool pool = new ThreadPool(workers);
        CountDownLatch together = new CountDownLatch(workers);
        try {
            List<Future<Boolean>> fs = new ArrayList<>();
            for (int i = 0; i < workers; i++) {
                fs.add(pool.submit(() -> {
                    together.countDown();
                    return together.await(WAIT_SECONDS, TimeUnit.SECONDS);
                }));
            }
            for (Future<Boolean> f : fs) assertTrue(get(f), "the tasks did not run at the same time");
        } finally {
            pool.shutdown();
        }
    }

    @Test
    void neverMoreThreadsThanWorkers() throws Exception {
        int workers = 3;
        ThreadPool pool = new ThreadPool(workers);
        Set<Thread> threads = Collections.synchronizedSet(new HashSet<>());
        CountDownLatch all = new CountDownLatch(500);
        try {
            for (int i = 0; i < 500; i++) {
                pool.execute(() -> {
                    threads.add(Thread.currentThread());
                    all.countDown();
                });
            }
            assertTrue(await(all));
        } finally {
            pool.shutdown();
        }
        assertTrue(threads.size() <= workers, "tasks ran on " + threads.size() + " threads; the pool has " + workers);
        assertFalse(threads.contains(Thread.currentThread()), "a task ran on the caller's thread");
    }

    // ---- errors ----

    @Test
    void aTaskThatThrowsDoesNotKillItsWorker() throws Exception {
        ThreadPool pool = new ThreadPool(1);
        try {
            for (int i = 0; i < 5; i++) {
                pool.execute(() -> {
                    throw new RuntimeException("a task failing on purpose");
                });
            }
            assertEquals("still alive", get(pool.submit(() -> "still alive")),
                    "the only worker died after a task threw");
        } finally {
            pool.shutdown();
        }
    }

    // ---- bounded queue: backpressure ----

    @Test
    void executeWaitsWhenTheQueueIsFullAndTrySubmitSaysNo() throws Exception {
        ThreadPool pool = new ThreadPool(1, 2);
        CountDownLatch started = new CountDownLatch(1);
        CountDownLatch release = new CountDownLatch(1);
        AtomicInteger ran = new AtomicInteger();
        try {
            pool.execute(() -> {                          // occupies the only worker
                started.countDown();
                try {
                    release.await();
                } catch (InterruptedException ignored) {
                }
                ran.incrementAndGet();
            });
            assertTrue(await(started));
            // The queue is empty now: two tasks fill it.
            assertTrue(pool.trySubmit(ran::incrementAndGet));
            assertTrue(pool.trySubmit(ran::incrementAndGet));
            assertFalse(pool.trySubmit(ran::incrementAndGet), "the queue is full: trySubmit must refuse");

            Future<Boolean> blocked = inBackground(() -> {
                pool.execute(ran::incrementAndGet);
                return true;
            });
            assertTrue(stillBlocked(blocked, 300), "execute must wait while the queue is full");

            release.countDown();
            assertTrue(get(blocked));
        } finally {
            release.countDown();
            pool.shutdown();
        }
        assertEquals(4, ran.get(), "the refused task must not run; all others must");
    }

    // ---- shutdown ----

    @Test
    void shutdownFinishesEveryAcceptedTaskAndEndsTheWorkers() throws Exception {
        ThreadPool pool = new ThreadPool(3, 10);
        AtomicInteger ran = new AtomicInteger();
        Set<Thread> threads = Collections.synchronizedSet(new HashSet<>());
        for (int i = 0; i < 40; i++) {
            pool.execute(() -> {
                threads.add(Thread.currentThread());
                try {
                    Thread.sleep(5);
                } catch (InterruptedException ignored) {
                }
                ran.incrementAndGet();
            });
        }
        pool.shutdown();
        assertEquals(40, ran.get(), "shutdown() returned before every accepted task had finished");
        for (Thread t : threads) assertFalse(t.isAlive(), "a worker thread is still alive after shutdown()");
        pool.shutdown();                                  // again: nothing happens
    }

    @Test
    void afterShutdownNothingIsAccepted() throws Exception {
        ThreadPool pool = new ThreadPool(2);
        pool.shutdown();
        assertThrows(RejectedExecutionException.class, () -> pool.execute(() -> {}));
        assertThrows(RejectedExecutionException.class, () -> pool.trySubmit(() -> {}));
        assertThrows(RejectedExecutionException.class, () -> pool.submit(() -> 1));
    }

    @Test
    void shutdownWakesAnExecuteThatWasWaitingForRoom() throws Exception {
        ThreadPool pool = new ThreadPool(1, 1);
        CountDownLatch started = new CountDownLatch(1);
        CountDownLatch release = new CountDownLatch(1);
        pool.execute(() -> {
            started.countDown();
            try {
                release.await();
            } catch (InterruptedException ignored) {
            }
        });
        assertTrue(await(started));
        assertTrue(pool.trySubmit(() -> {}));              // the queue is full now
        Future<Boolean> waiting = inBackground(() -> {
            try {
                pool.execute(() -> {});
                return false;
            } catch (RejectedExecutionException e) {
                return true;
            }
        });
        assertTrue(stillBlocked(waiting, 200));
        Future<Void> stopping = inBackground(() -> {
            pool.shutdown();
            return null;
        });
        assertTrue(get(waiting), "execute() must throw RejectedExecutionException when the pool shuts down under it");
        release.countDown();
        get(stopping);
    }

    @Test
    void aTaskMaySubmitMoreTasks() throws Exception {
        // Spawning without waiting is fine: the parent does not block.
        ThreadPool pool = new ThreadPool(2);
        CountDownLatch all = new CountDownLatch(101);
        try {
            pool.execute(() -> {
                for (int i = 0; i < 100; i++) {
                    try {
                        pool.execute(all::countDown);
                    } catch (InterruptedException e) {
                        return;
                    }
                }
                all.countDown();
            });
            assertTrue(await(all));
        } finally {
            pool.shutdown();
        }
    }
}
