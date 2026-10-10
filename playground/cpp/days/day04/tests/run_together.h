#pragma once
// Test helpers for Day 4.
//
//   runTogether / runTogetherWithin   run a body on n threads that start at
//                                     the same moment (as on Day 3)
//   Background<T>                     run one call on another thread, so
//                                     the test can check that it is still
//                                     waiting -- and what it returns later
//   Gate                              a door that tasks wait at until the
//                                     test opens it
//   waitUntil                         poll a condition, with a time limit

#include <atomic>
#include <chrono>
#include <condition_variable>
#include <cstdio>
#include <cstdlib>
#include <functional>
#include <future>
#include <mutex>
#include <thread>
#include <vector>

inline void runTogether(int n, const std::function<void(int)>& body) {
    std::atomic<int> ready{0};
    std::atomic<bool> go{false};
    std::vector<std::thread> threads;
    for (int t = 0; t < n; ++t) {
        threads.emplace_back([&, t] {
            ready.fetch_add(1);
            while (!go.load()) std::this_thread::yield();
            body(t);
        });
    }
    while (ready.load() < n) std::this_thread::yield();
    go.store(true);
    for (auto& th : threads) th.join();
}

// Like runTogether, for code that might deadlock: if the threads are
// still running after `seconds`, stop the whole test run with a message
// instead of hanging. (A deadlocked thread can't be stopped any other
// way.)
inline void runTogetherWithin(int seconds, int n, const std::function<void(int)>& body) {
    std::promise<void> done;
    std::future<void> doneF = done.get_future();
    std::thread runner([&] {
        runTogether(n, body);
        done.set_value();
    });
    if (doneF.wait_for(std::chrono::seconds(seconds)) != std::future_status::ready) {
        std::fprintf(stderr,
                     "\nThreads still running after %d s -- deadlocked? "
                     "A thread that sleeps and is never woken up (a lost "
                     "notification), or a wait that nobody ends. "
                     "Stopping the test run.\n",
                     seconds);
        std::fflush(stdout);
        std::fflush(stderr);
        std::_Exit(1);
    }
    runner.join();
}

// One call running on its own thread, started right away.
//
//   Background<bool> b([&] { return q.push(1); });
//   b.waitFor(100ms)   did the call return within 100 ms? (false: it is
//                      still waiting)
//   b.result()         the call's return value; throws what it threw. If
//                      the call is still running after `seconds`, stops
//                      the whole test run with a message.
//
// The destructor waits for the call to return (at most 10 s, then it
// stops the test run too), so a test must unblock its background calls
// before it ends.
template <typename T>
class Background {
public:
    explicit Background(std::function<T()> fn) : future_(promise_.get_future().share()) {
        thread_ = std::thread([this, fn = std::move(fn)] {
            try {
                promise_.set_value(fn());
            } catch (...) {
                promise_.set_exception(std::current_exception());
            }
        });
    }

    Background(const Background&) = delete;
    Background& operator=(const Background&) = delete;

    bool waitFor(std::chrono::milliseconds d) {
        return future_.wait_for(d) == std::future_status::ready;
    }

    T result(int seconds = 10) {
        if (!waitFor(std::chrono::seconds(seconds))) giveUp(seconds);
        return future_.get();
    }

    ~Background() {
        if (!waitFor(std::chrono::seconds(10))) giveUp(10);
        thread_.join();
    }

private:
    [[noreturn]] static void giveUp(int seconds) {
        std::fprintf(stderr,
                     "\nA call is still running after %d s -- deadlocked, or "
                     "waiting for a wake-up that never comes? "
                     "Stopping the test run.\n",
                     seconds);
        std::fflush(stdout);
        std::fflush(stderr);
        std::_Exit(1);
    }

    std::promise<T> promise_;
    std::shared_future<T> future_;  // shared: result() and the destructor both look at it
    std::thread thread_;
};

// A door. Tasks call wait(); the test calls open(). wait() gives up after
// 10 s so that a failing test never hangs on a closed gate.
class Gate {
public:
    void open() {
        {
            std::lock_guard<std::mutex> lk(m_);
            open_ = true;
        }
        cv_.notify_all();
    }

    bool wait() {
        std::unique_lock<std::mutex> lk(m_);
        return cv_.wait_for(lk, std::chrono::seconds(10), [&] { return open_; });
    }

private:
    std::mutex m_;
    std::condition_variable cv_;
    bool open_ = false;
};

// Polls pred() every millisecond until it is true or `timeout` has passed.
// Returns whether pred() became true.
inline bool waitUntil(const std::function<bool()>& pred, std::chrono::milliseconds timeout) {
    const auto deadline = std::chrono::steady_clock::now() + timeout;
    while (!pred()) {
        if (std::chrono::steady_clock::now() > deadline) return false;
        std::this_thread::sleep_for(std::chrono::milliseconds(1));
    }
    return true;
}
