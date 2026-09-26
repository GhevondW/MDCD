#pragma once
// Test helper: run body(0) ... body(n-1) on n threads that all start at
// (nearly) the same moment, so their calls overlap as much as possible.

#include <atomic>
#include <chrono>
#include <cstdio>
#include <cstdlib>
#include <functional>
#include <future>
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
                     "Two threads locking the same locks in opposite orders, "
                     "or one thread locking the same lock twice. "
                     "Stopping the test run.\n",
                     seconds);
        std::fflush(stderr);
        std::_Exit(1);
    }
    runner.join();
}
