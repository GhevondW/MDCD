#pragma once
// Test helper: run body(0) ... body(n-1) on n threads that all start at
// (nearly) the same moment, so their calls overlap as much as possible.

#include <atomic>
#include <functional>
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
