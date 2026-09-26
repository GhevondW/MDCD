package day03;

import java.util.concurrent.Callable;
import java.util.concurrent.Future;
import java.util.concurrent.FutureTask;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicInteger;

/**
 * Test helpers for starting threads.
 *
 * JUnit only notices an assertion that fails on the test's own thread. An
 * AssertionError thrown on any other thread is simply lost, and the test
 * would pass. So these helpers catch whatever a thread throws and hand it
 * back to the test's thread.
 *
 * All threads started here are daemon threads: if a broken solution
 * deadlocks one of them, it cannot keep the JVM alive after the tests.
 */
final class RunTogether {
    private RunTogether() {}

    // Well under each test class's 60 s limit, so a stuck thread is reported
    // with its own stack trace instead of a bare timeout.
    private static final long STUCK_AFTER_SECONDS = 45;

    /** The work of thread number t, for t = 0 .. n-1. */
    @FunctionalInterface
    interface Body {
        void run(int t) throws Exception;
    }

    /**
     * Runs body.run(0) ... body.run(n-1) on n threads that all start at
     * (nearly) the same moment, so their calls overlap as much as possible.
     * Waits for all of them to finish. If any of them threw, rethrows the
     * first one here (the others are attached as suppressed). A thread that
     * has not finished after 45 s fails the test with that thread's stack
     * trace.
     */
    static void runTogether(int n, Body body) throws Exception {
        AtomicInteger ready = new AtomicInteger();
        // Spin on a flag instead of blocking on a latch: waking blocked
        // threads one by one would spread their starts out, and single-call
        // races would get much harder to hit.
        AtomicInteger go = new AtomicInteger();
        Throwable[] thrown = new Throwable[n];
        Thread[] threads = new Thread[n];
        for (int t = 0; t < n; t++) {
            int id = t;
            threads[t] = new Thread(() -> {
                ready.incrementAndGet();
                while (go.get() == 0) Thread.yield();
                try {
                    body.run(id);
                } catch (Throwable e) {
                    thrown[id] = e;
                }
            }, "runTogether-" + t);
            threads[t].setDaemon(true);
            threads[t].start();
        }
        while (ready.get() < n) Thread.yield();
        go.set(1);

        long deadline = System.nanoTime() + TimeUnit.SECONDS.toNanos(STUCK_AFTER_SECONDS);
        for (Thread th : threads) {
            long msLeft = TimeUnit.NANOSECONDS.toMillis(deadline - System.nanoTime());
            if (msLeft > 0) th.join(msLeft);
        }

        Throwable first = null;
        for (Thread th : threads) {
            if (th.isAlive()) {
                // Report where the stuck thread is right now -- that is the
                // line of your code it cannot get past.
                first = new AssertionError(th.getName() + " is still running after "
                        + STUCK_AFTER_SECONDS + " s -- a deadlock, or an endless loop?"
                        + " The stack trace below is where that thread is stuck.");
                first.setStackTrace(th.getStackTrace());
                break;
            }
        }
        for (int t = 0; t < n; t++) {
            Throwable e = threads[t].isAlive() ? null : thrown[t];
            if (e == null) continue;
            if (first == null) {
                first = e;
            } else if (e != first) {
                first.addSuppressed(e);
            }
        }
        if (first instanceof Exception e) throw e;
        if (first instanceof Error e) throw e;
        if (first != null) throw new AssertionError(first);
    }

    /**
     * Starts task on a new daemon thread and returns right away. The Future
     * gives the task's result, or (wrapped in ExecutionException) what it
     * threw.
     */
    static <T> Future<T> inBackground(Callable<T> task) {
        FutureTask<T> future = new FutureTask<>(task);
        Thread th = new Thread(future, "inBackground");
        th.setDaemon(true);
        th.start();
        return future;
    }
}
