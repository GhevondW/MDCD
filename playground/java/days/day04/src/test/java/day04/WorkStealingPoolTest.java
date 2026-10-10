package day04;

import day04.WorkStealingPool.Task;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.Timeout;

import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.concurrent.CompletionException;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.Future;
import java.util.concurrent.RejectedExecutionException;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;
import java.util.concurrent.atomic.AtomicInteger;

import static day04.RunTogether.inBackground;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertInstanceOf;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

@Timeout(value = 60, threadMode = Timeout.ThreadMode.SEPARATE_THREAD)
class WorkStealingPoolTest {

    private static final long WAIT_SECONDS = 5;

    private static boolean await(CountDownLatch latch) throws InterruptedException {
        return latch.await(WAIT_SECONDS, TimeUnit.SECONDS);
    }

    private static <T> T get(Future<T> f) throws Exception {
        try {
            return f.get(WAIT_SECONDS, TimeUnit.SECONDS);
        } catch (TimeoutException e) {
            throw new AssertionError("did not finish within " + WAIT_SECONDS + " s -- a deadlock?");
        }
    }

    /** Runs `body` as a task of the pool (so join() happens on a worker) and returns its result. */
    private static <T> T runInPool(WorkStealingPool pool, java.util.concurrent.Callable<T> body) throws Exception {
        return pool.submit(body).join();
    }

    // ---- recursive fork/join ----

    private static int fib(WorkStealingPool pool, int n) throws Exception {
        if (n < 2) return n;
        Task<Integer> left = pool.submit(() -> fib(pool, n - 1));   // spawn...
        int right = fib(pool, n - 2);                                // ...compute the other half here...
        return left.join() + right;                                  // ...and join (helping!)
    }

    private static int fibPlain(int n) {
        return n < 2 ? n : fibPlain(n - 1) + fibPlain(n - 2);
    }

    private void checkFib(int workers) throws Exception {
        WorkStealingPool pool = new WorkStealingPool(workers);
        try {
            assertEquals(fibPlain(21), runInPool(pool, () -> fib(pool, 21)));
        } finally {
            pool.shutdown();
        }
    }

    @Test
    void fibWithOneWorkerDoesNotDeadlock() throws Exception {
        // One worker, and every task waits for a task it spawned. With a
        // plain pool this is the lecture's deadlock; here join() must run
        // the awaited task itself.
        checkFib(1);
    }

    @Test
    void fibWithTwoWorkers() throws Exception {
        checkFib(2);
    }

    @Test
    void fibWithFourWorkers() throws Exception {
        checkFib(4);
    }

    private static long sum(WorkStealingPool pool, int[] a, int from, int to) throws Exception {
        if (to - from <= 1000) {
            long s = 0;
            for (int i = from; i < to; i++) s += a[i];
            return s;
        }
        int mid = (from + to) >>> 1;
        Task<Long> left = pool.submit(() -> sum(pool, a, from, mid));
        long right = sum(pool, a, mid, to);
        return left.join() + right;
    }

    @Test
    void divideAndConquerSum() throws Exception {
        int[] a = new int[1_000_000];
        for (int i = 0; i < a.length; i++) a[i] = i % 1000;
        long expected = 0;
        for (int x : a) expected += x;
        for (int workers : new int[] {1, 3}) {
            WorkStealingPool pool = new WorkStealingPool(workers);
            try {
                assertEquals(expected, (long) runInPool(pool, () -> sum(pool, a, 0, a.length)),
                        workers + " worker(s)");
            } finally {
                pool.shutdown();
            }
        }
    }

    @Test
    void joinFromOutsideThePoolWaits() throws Exception {
        WorkStealingPool pool = new WorkStealingPool(2);
        try {
            Task<String> t = pool.submit(() -> {
                Thread.sleep(100);
                return "done";
            });
            assertEquals("done", t.join());
            assertTrue(t.isDone());
        } finally {
            pool.shutdown();
        }
    }

    @Test
    void aFailingTaskSurfacesInJoinAndTheWorkerSurvives() throws Exception {
        WorkStealingPool pool = new WorkStealingPool(1);
        try {
            Task<Integer> bad = pool.submit(() -> {
                throw new IllegalStateException("boom");
            });
            CompletionException e = assertThrows(CompletionException.class, bad::join);
            assertInstanceOf(IllegalStateException.class, e.getCause());

            pool.execute(() -> {
                throw new RuntimeException("a raw task failing on purpose");
            });
            assertEquals(9, pool.submit(() -> 9).join(), "the only worker died");
        } finally {
            pool.shutdown();
        }
    }

    // ---- the deque discipline ----

    @Test
    void aWorkerRunsItsOwnNewestTaskFirst() throws Exception {
        // One worker. A task spawns A, B, C and finishes without joining.
        // The worker then takes from its own deque, newest first: C, B, A.
        WorkStealingPool pool = new WorkStealingPool(1);
        List<String> order = Collections.synchronizedList(new ArrayList<>());
        CountDownLatch all = new CountDownLatch(3);
        try {
            pool.execute(() -> {
                for (String name : new String[] {"A", "B", "C"}) {
                    pool.submit(() -> {
                        order.add(name);
                        all.countDown();
                        return null;
                    });
                }
            });
            assertTrue(await(all));
        } finally {
            pool.shutdown();
        }
        assertEquals(List.of("C", "B", "A"), order, "the owner takes its newest task first");
        assertEquals(0, pool.steals(), "a lone worker has nobody to steal from");
    }

    @Test
    void aThiefTakesTheOldestTaskFirst() throws Exception {
        // Two workers. A task spawns A, B, C and then blocks (it does NOT
        // help) until all three have run. The other worker is idle: it must
        // steal them, oldest first: A, B, C.
        WorkStealingPool pool = new WorkStealingPool(2);
        List<String> order = Collections.synchronizedList(new ArrayList<>());
        CountDownLatch all = new CountDownLatch(3);
        CountDownLatch parentDone = new CountDownLatch(1);
        try {
            pool.execute(() -> {
                for (String name : new String[] {"A", "B", "C"}) {
                    pool.submit(() -> {
                        order.add(name);
                        all.countDown();
                        return null;
                    });
                }
                try {
                    all.await(WAIT_SECONDS, TimeUnit.SECONDS);   // blocks this worker
                } catch (InterruptedException ignored) {
                }
                parentDone.countDown();
            });
            assertTrue(await(parentDone), "the idle worker never took the spawned tasks");
        } finally {
            pool.shutdown();
        }
        assertEquals(List.of("A", "B", "C"), order, "a thief takes the oldest task first");
        assertEquals(3, pool.steals());
    }

    @Test
    void anIdleWorkerStealsFromABusyOne() throws Exception {
        // One task spawns 160 children of 5 ms each and joins them. Run one
        // after the other that is 800 ms; four workers can do it in about
        // 200 ms -- but only if the three idle ones steal.
        int children = 160;
        WorkStealingPool pool = new WorkStealingPool(4);
        long start = System.nanoTime();
        try {
            int done = runInPool(pool, () -> {
                List<Task<Integer>> tasks = new ArrayList<>();
                for (int i = 0; i < children; i++) {
                    tasks.add(pool.submit(() -> {
                        Thread.sleep(5);
                        return 1;
                    }));
                }
                int n = 0;
                for (Task<Integer> t : tasks) n += t.join();
                return n;
            });
            assertEquals(children, done);
        } finally {
            pool.shutdown();
        }
        long ms = TimeUnit.NANOSECONDS.toMillis(System.nanoTime() - start);
        assertTrue(pool.steals() > 0, "no task was ever stolen");
        assertTrue(ms < 560, "took " + ms + " ms; 4 workers should need about 200 ms (800 ms if run one at a time)");
    }

    // ---- shutdown ----

    @Test
    void shutdownRunsEverythingTheTasksSpawnEvenWhileShuttingDown() throws Exception {
        WorkStealingPool pool = new WorkStealingPool(2);
        AtomicInteger ran = new AtomicInteger();
        CountDownLatch started = new CountDownLatch(1);
        CountDownLatch go = new CountDownLatch(1);
        pool.execute(() -> {
            started.countDown();
            try {
                go.await(WAIT_SECONDS, TimeUnit.SECONDS);
            } catch (InterruptedException ignored) {
            }
            // shutdown() has been called by now. These come from a worker,
            // so they are accepted -- and so are their children.
            for (int i = 0; i < 50; i++) {
                pool.execute(() -> {
                    ran.incrementAndGet();
                    pool.execute(ran::incrementAndGet);
                });
            }
            ran.incrementAndGet();
        });
        assertTrue(await(started));
        Future<Void> stopping = inBackground(() -> {
            pool.shutdown();
            return null;
        });
        Thread.sleep(200);                                  // let shutdown() begin
        go.countDown();
        get(stopping);
        assertEquals(1 + 50 + 50, ran.get(), "shutdown() returned before all the spawned work was done");
    }

    @Test
    void afterShutdownOutsideCallersAreRejected() throws Exception {
        WorkStealingPool pool = new WorkStealingPool(2);
        assertEquals(3, pool.submit(() -> 3).join());
        pool.shutdown();
        assertThrows(RejectedExecutionException.class, () -> pool.execute(() -> {}));
        assertThrows(RejectedExecutionException.class, () -> pool.submit(() -> 1));
        pool.shutdown();                                    // again: nothing happens
    }

    @Test
    void manyOutsideSubmittersAtOnce() throws Exception {
        WorkStealingPool pool = new WorkStealingPool(4);
        int submitters = 6;
        int each = 1000;
        AtomicInteger ran = new AtomicInteger();
        try {
            RunTogether.runTogether(submitters, t -> {
                for (int i = 0; i < each; i++) pool.execute(ran::incrementAndGet);
            });
        } finally {
            pool.shutdown();
        }
        assertEquals(submitters * each, ran.get());
    }
}
