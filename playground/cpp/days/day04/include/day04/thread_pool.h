#pragma once
// Day 4 -- The Thread Pool
//
// Keep a few threads busy with many tasks: a fixed set of worker threads
// takes tasks from a queue and runs them (the slides' "thread pool in 20
// lines"). You write the workers, the queue and the shutdown. submit() is
// already written: it wraps your callable so that its result -- or its
// exception -- arrives in a std::future, and hands it to execute().
//
//   ThreadPool(workers, queueCapacity)
//       `workers` threads start at once; workers must be at least 1,
//       otherwise throw std::invalid_argument. `queueCapacity` is how many
//       tasks may WAIT in the queue (tasks that are running do not count);
//       0 means the queue has no limit.
//
//   execute(task)      hand the task to the pool; one of the workers will
//                      run it. If the queue is full, wait until there is
//                      room (backpressure: a full queue slows down whoever
//                      fills it). If the pool is shut down -- already, or
//                      while this call waits -- throw std::runtime_error.
//   trySubmit(task)    like execute, but never waits: returns true if the
//                      task was accepted, false if the queue is full or
//                      the pool is shut down. Never throws.
//   shutdown()         stop accepting tasks, run every task that was
//                      already accepted -- queued ones too -- and join the
//                      worker threads. When it returns, all of them have
//                      finished. Calling it again does nothing. The
//                      destructor calls it. (Do not call it from inside a
//                      task: a worker cannot join itself.)
//   workerCount()      the number of worker threads.
//   submit(f)          provided. Returns a std::future for f()'s result.
//
// The rules that make it hard:
//
//   1. A task that throws must not kill its worker -- an exception that
//      escapes a std::thread ends the whole program. Catch it and carry on.
//      (submit's tasks never throw: the exception goes into the future.)
//   2. Workers sleep while there is nothing to do. No busy waiting.
//   3. Tasks run in parallel: with 4 workers, 4 slow tasks take about as
//      long as one. And there are never more threads than `workers`.
//   4. shutdown() must wake the idle workers and the submitters that are
//      waiting for room, and must not lose a single accepted task.
//   5. execute() and trySubmit() are called by many threads at once.
//
// A trap, not a task: if a task waits for a future of the SAME pool, it
// can use up every worker and deadlock (the slides' "deadlock with no
// locks at all"). Challenge 3, the work-stealing pool, is the way out.
//
// You may use your queue from the first task, or write the queue inside
// the pool -- either is fine. Don't use anything that is already a thread
// pool.

#include <condition_variable>
#include <cstddef>
#include <functional>
#include <future>
#include <memory>
#include <mutex>
#include <stdexcept>
#include <thread>
#include <type_traits>
#include <utility>
#include <vector>

class ThreadPool {
public:
    explicit ThreadPool(std::size_t workers, std::size_t queueCapacity = 0) {
        (void)workers;
        (void)queueCapacity;
        // TODO
    }

    ~ThreadPool() { shutdown(); }

    ThreadPool(const ThreadPool&) = delete;
    ThreadPool& operator=(const ThreadPool&) = delete;

    void execute(std::function<void()> task) {
        (void)task;
        // TODO
    }

    bool trySubmit(std::function<void()> task) {
        (void)task;
        // TODO
        return false;
    }

    void shutdown() {
        // TODO
    }

    std::size_t workerCount() const {
        // TODO
        return 0;
    }

    // Provided -- nothing to do here. A packaged_task is a function plus a
    // promise: running it fills the future with the value or the exception.
    template <typename F>
    auto submit(F f) -> std::future<std::invoke_result_t<F&>> {
        using R = std::invoke_result_t<F&>;
        auto task = std::make_shared<std::packaged_task<R()>>(std::move(f));
        std::future<R> result = task->get_future();
        execute([task] { (*task)(); });
        return result;
    }

private:
    // TODO: choose your own representation.
};
