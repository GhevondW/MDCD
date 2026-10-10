#pragma once
// Day 4 -- Challenge: The Work-Stealing Pool
//
// The slides ended with two problems of the plain thread pool:
//
//   - One shared queue is one lock: every submit and every take fights
//     for it, and the more workers, the more they wait.
//   - A task that waits for a child task of the same pool can use up all
//     the workers and deadlock.
//
// This pool fixes both, the way Java's ForkJoinPool, Rust's rayon and the
// .NET pool do:
//
//   - Every worker has its OWN deque of tasks. A task that a running task
//     creates (a "fork") goes onto the back of that worker's own deque.
//     The owner takes tasks from the back: newest first (LIFO), which is
//     fast and cache-friendly. A worker whose deque is empty STEALS: it
//     takes the OLDEST task (the front) of another worker's deque, which
//     is usually the biggest piece of work, and leaves the owner's hot
//     end alone.
//   - Tasks that come from threads that are not workers go into one
//     shared injection queue (FIFO). A worker looks for work in this
//     order: its own deque, then the injection queue, then the deques of
//     the other workers.
//   - While a worker waits for a result (Handle::get), it does not
//     sleep: it keeps running other tasks until the result is there. A
//     pool with ONE worker can therefore run recursive code like
//     fib(n) = fib(n-1) + fib(n-2) with a task per call.
//
// Already written: Handle<R> (the result of a task) and submit(f). What
// you write:
//
//   WorkStealingPool(workers)
//       starts `workers` threads (at least 1, otherwise throw
//       std::invalid_argument).
//   execute(task)
//       - called on one of this pool's own workers (a task forking a
//         task): push onto the back of that worker's deque. Always
//         accepted, even while the pool is shutting down.
//       - called from any other thread: if shutdown() has begun, throw
//         std::runtime_error. Otherwise put it into the injection queue.
//       A task is either accepted and will run, or rejected with an
//       exception -- never accepted and then dropped.
//       A task that throws must not kill its worker: catch and go on.
//   helpUntil(done)
//       returns when done() is true. On a worker: keep finding and running
//       tasks (same search order as above) until done(); when there is
//       nothing to find right now, yield or sleep for a moment -- the
//       task you wait for is running on another worker. On any other
//       thread: just wait until done() (polling with a short sleep is
//       fine; there is no waiting room to sleep on).
//   shutdown()
//       stop accepting tasks from other threads, wait until EVERY
//       accepted task -- including the tasks that running tasks fork in
//       the meantime -- has finished, then stop and join the workers.
//       Calling it again does nothing; the destructor calls it. (Do not
//       call it from inside a task.)
//   workerCount()   the number of worker threads.
//   steals()        how many tasks workers have taken from ANOTHER
//                   worker's deque. Tasks from the injection queue are
//                   not steals.
//
// The rules that make it hard:
//
//   1. Each deque has its own lock (a lock-free deque is a bonus, not a
//      requirement); the owner works at the back and thieves at the
//      front. Never hold two deque locks at once.
//   2. An idle worker must not burn CPU at full speed, but it must find
//      new work soon: sleep a moment between rounds, or wait on a
//      condition variable that execute() notifies.
//   3. "Is everything done?" is the hard part of shutdown. A task that is
//      running may still fork new tasks, so empty queues alone do not
//      mean the work is over. Count the tasks that were accepted and have
//      not finished yet.
//   4. How does a thread know it is a worker of this pool, and which one?
//      A thread_local variable set by the worker loop.
//
// Don't use anything that is already a thread pool or a concurrent deque.

#include <atomic>
#include <chrono>
#include <condition_variable>
#include <cstddef>
#include <deque>
#include <exception>
#include <functional>
#include <future>
#include <memory>
#include <mutex>
#include <stdexcept>
#include <thread>
#include <type_traits>
#include <utility>
#include <vector>

class WorkStealingPool {
public:
    // Provided. The result of a task: get() waits for it -- and while it
    // waits, helps the pool by running other tasks (helpUntil).
    template <typename R>
    class Handle {
    public:
        Handle(Handle&&) = default;
        Handle& operator=(Handle&&) = default;

        bool ready() const { return done_->load(std::memory_order_acquire); }

        // The value, or the exception the task threw. Once only.
        R get() {
            auto done = done_;
            pool_->helpUntil([done] { return done->load(std::memory_order_acquire); });
            return future_.get();
        }

    private:
        friend class WorkStealingPool;
        Handle(WorkStealingPool* pool, std::future<R> future,
               std::shared_ptr<std::atomic<bool>> done)
            : pool_(pool), future_(std::move(future)), done_(std::move(done)) {}

        WorkStealingPool* pool_;
        std::future<R> future_;
        std::shared_ptr<std::atomic<bool>> done_;
    };

    explicit WorkStealingPool(std::size_t workers) {
        (void)workers;
        // TODO
    }

    ~WorkStealingPool() { shutdown(); }

    WorkStealingPool(const WorkStealingPool&) = delete;
    WorkStealingPool& operator=(const WorkStealingPool&) = delete;

    void execute(std::function<void()> task) {
        (void)task;
        // TODO
    }

    void helpUntil(const std::function<bool()>& done) {
        (void)done;
        // TODO
    }

    void shutdown() {
        // TODO
    }

    std::size_t workerCount() const {
        // TODO
        return 0;
    }

    std::size_t steals() const {
        // TODO
        return 0;
    }

    // Provided. Same idea as ThreadPool::submit.
    template <typename F>
    auto submit(F f) -> Handle<std::invoke_result_t<F&>> {
        using R = std::invoke_result_t<F&>;
        auto task = std::make_shared<std::packaged_task<R()>>(std::move(f));
        auto done = std::make_shared<std::atomic<bool>>(false);
        Handle<R> handle(this, task->get_future(), done);
        execute([task, done] {
            (*task)();
            done->store(true, std::memory_order_release);
        });
        return handle;
    }

private:
    // TODO: choose your own representation.
};
