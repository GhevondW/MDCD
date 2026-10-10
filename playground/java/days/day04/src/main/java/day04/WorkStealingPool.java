package day04;

import java.util.concurrent.Callable;
import java.util.concurrent.CompletionException;
import java.util.function.BooleanSupplier;

/**
 * Day 4 (challenge) -- A work-stealing pool.
 *
 * One shared queue is one lock that every worker fights for. And a worker
 * that waits for a task it spawned sleeps -- and takes a whole worker
 * with it (the lecture's "deadlock with no locks at all"). This pool fixes
 * both:
 *
 *   - Every worker has its OWN deque of tasks. A task spawned by code
 *     running on a worker goes to the back ("bottom") of that worker's
 *     deque. The owner takes its own work from the back too -- the
 *     NEWEST task first (LIFO): it is the one most likely still "warm".
 *   - Tasks submitted from outside the pool (any other thread) go to one
 *     shared injection queue, in first-in-first-out order.
 *   - A worker with nothing in its own deque takes a task from the
 *     injection queue, or STEALS from another worker's deque: it takes
 *     the OLDEST task (from the front, FIFO) -- the owner is busy at the
 *     other end, and the oldest task is usually the biggest.
 *   - join() on a task from a worker does not sleep: the worker keeps
 *     running other tasks (its own, then stolen ones) until the awaited
 *     task is done. So recursive fork/join code works even with ONE
 *     worker.
 *
 * You write: execute, helpUntil, the worker loop, shutdown, steals().
 * The task handle (Task) and submit() are provided. Do not use
 * ForkJoinPool, ExecutorService or any java.util.concurrent queue/deque;
 * plain synchronized or ReentrantLock deques are fine (a lock-free deque
 * is a research topic, not a requirement).
 *
 * Constructor:
 *   - WorkStealingPool(workers): starts `workers` daemon threads.
 *     Throws IllegalArgumentException if workers < 1.
 *
 * execute(task):
 *   - Called on a worker thread OF THIS POOL: push the task onto that
 *     worker's own deque (bottom). Never wait, never reject (see
 *     shutdown).
 *   - Called on any other thread: put it on the injection queue. Throws
 *     RejectedExecutionException if shutdown() has been called.
 *   - Either way, an idle worker must notice the new task and take it
 *     (workers that have nothing to do should sleep, not spin -- but they
 *     must wake up when work appears anywhere, including in another
 *     worker's deque).
 *   - A task that throws must not kill its worker.
 *
 * helpUntil(done):
 *   - Called on a worker thread of this pool: run tasks -- own deque
 *     (newest first), then the injection queue, then steal -- until
 *     done.getAsBoolean() is true. If there is nothing to run right now
 *     but done is still false (the awaited task is running on another
 *     worker), wait a moment (Thread.yield(), a short park/sleep) and
 *     look again. Return as soon as done is true; do not run one more
 *     task after that.
 *   - Called on any other thread: return at once. (Task.join waits for
 *     those threads by itself.)
 *
 * shutdown():
 *   - From now on, execute() on a thread that is NOT one of the pool's
 *     workers throws RejectedExecutionException.
 *   - Every task accepted so far runs to completion -- including tasks
 *     that running tasks spawn WHILE the pool is shutting down (those
 *     are accepted, since they come from a worker). Only when no
 *     accepted task is left, running or waiting anywhere, do the workers
 *     stop. shutdown() returns after they have all ended.
 *   - Calling it again does nothing. Not from inside a task.
 *
 * steals(): how many tasks, in total, were taken from ANOTHER worker's
 * deque by a worker (the injection queue does not count).
 */
public class WorkStealingPool {
    // TODO: choose your own representation.

    public WorkStealingPool(int workers) {
        // TODO: validate, create the deques, start the workers.
    }

    public void execute(Runnable task) {
        // TODO
    }

    public void helpUntil(BooleanSupplier done) {
        // TODO
    }

    public void shutdown() throws InterruptedException {
        // TODO
    }

    public long steals() {
        // TODO
        return 0;
    }

    /** Provided: spawn a task with a result; join() it later. */
    public <T> Task<T> submit(Callable<T> callable) {
        Task<T> task = new Task<>(this, callable);
        execute(task::run);
        return task;
    }

    /** Provided: the handle of a spawned task. */
    public static final class Task<T> {
        private static final long GIVE_UP_MILLIS = 5_000;

        private final WorkStealingPool pool;
        private final Callable<T> body;
        private T result;
        private Throwable failure;
        private volatile boolean done;

        Task(WorkStealingPool pool, Callable<T> body) {
            this.pool = pool;
            this.body = body;
        }

        void run() {
            T value = null;
            Throwable error = null;
            try {
                value = body.call();
            } catch (Throwable e) {
                error = e;
            }
            synchronized (this) {
                result = value;
                failure = error;
                done = true;
                notifyAll();
            }
        }

        public boolean isDone() {
            return done;
        }

        /**
         * The task's result. On a pool worker this helps run other tasks
         * until this one is done; elsewhere it waits. If the task threw,
         * join() throws CompletionException with that as the cause.
         */
        public T join() throws InterruptedException {
            pool.helpUntil(this::isDone);
            synchronized (this) {
                long deadline = System.nanoTime() + GIVE_UP_MILLIS * 1_000_000;
                while (!done) {
                    long left = (deadline - System.nanoTime()) / 1_000_000;
                    if (left <= 0) {
                        throw new IllegalStateException(
                                "the task was still not done after " + GIVE_UP_MILLIS
                                        + " ms -- was it ever run?");
                    }
                    wait(left);
                }
                if (failure != null) throw new CompletionException(failure);
                return result;
            }
        }
    }
}
