package day04;

import java.util.concurrent.Callable;
import java.util.concurrent.Future;
import java.util.concurrent.FutureTask;

/**
 * Day 4 -- A thread pool.
 *
 * A fixed number of worker threads run many tasks. Build it yourself from
 * threads, a lock and conditions (or from your own BoundedBlockingQueue):
 * do not use ExecutorService, Executors, ThreadPoolExecutor, ForkJoinPool
 * or any java.util.concurrent queue.
 *
 * Constructor:
 *   - ThreadPool(workers, queueCapacity): starts `workers` threads right
 *     away. The queue holds the tasks that are waiting (not the ones
 *     running) -- at most queueCapacity of them. queueCapacity 0 means
 *     "no limit". Throws IllegalArgumentException if workers < 1 or
 *     queueCapacity < 0.
 *   - ThreadPool(workers): the same with no limit.
 *   Make the workers daemon threads, so that a pool nobody shut down
 *   cannot keep the program alive.
 *
 * Giving the pool work:
 *   - execute(task): hands the task to the pool. If the queue is full,
 *     waits until there is room (backpressure). Throws
 *     RejectedExecutionException if the pool is shut down -- also when it
 *     is shut down while execute is waiting for room.
 *   - trySubmit(task): like execute, but never waits: returns false if
 *     the queue is full (and the task is not run). Throws
 *     RejectedExecutionException if the pool is shut down.
 *   - submit(callable): PROVIDED. Wraps the callable in a Future and
 *     calls execute. The Future gets the result -- or the exception the
 *     callable threw (future.get() then throws ExecutionException).
 *
 * Running tasks:
 *   - Each task runs on one of the workers, exactly once. Tasks start in
 *     the order they were accepted.
 *   - Up to `workers` tasks run at the same time.
 *   - A task that throws must not kill its worker: the worker catches
 *     it and takes the next task. (A raw Runnable passed to execute
 *     has no Future to carry the exception; drop it.)
 *
 * Shutdown:
 *   - shutdown(): stops accepting new tasks, lets the workers finish
 *     EVERY task that was already accepted -- running or waiting -- then
 *     stops them and waits until they have all ended. When shutdown()
 *     returns, no worker thread is alive. Calling it again does nothing.
 *     (It may not be called from inside one of the pool's own tasks.)
 */
public class ThreadPool {
    // TODO: choose your own representation.

    public ThreadPool(int workers, int queueCapacity) {
        // TODO: validate the arguments, start the workers.
    }

    public ThreadPool(int workers) {
        this(workers, 0);
    }

    public void execute(Runnable task) throws InterruptedException {
        // TODO
    }

    public boolean trySubmit(Runnable task) {
        // TODO
        return false;
    }

    public void shutdown() throws InterruptedException {
        // TODO
    }

    /** Provided: a task with a result. */
    public <T> Future<T> submit(Callable<T> callable) throws InterruptedException {
        FutureTask<T> future = new FutureTask<>(callable);
        execute(future);
        return future;
    }
}
