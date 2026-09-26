#include "day03/once_cache.h"
#include "run_together.h"
#include <gtest/gtest.h>

#include <atomic>
#include <chrono>
#include <future>
#include <stdexcept>
#include <string>
#include <thread>
#include <vector>

using namespace std::chrono_literals;

TEST(OnceCache, MissComputesAndReturnsTheValue) {
    OnceCache<int> c;
    int calls = 0;
    EXPECT_EQ(c.get("a", [&] { ++calls; return 42; }), 42);
    EXPECT_EQ(calls, 1);
    EXPECT_EQ(c.size(), 1u);
}

TEST(OnceCache, HitReturnsRememberedValueWithoutComputing) {
    OnceCache<int> c;
    c.get("a", [] { return 42; });
    bool called = false;
    EXPECT_EQ(c.get("a", [&] { called = true; return 7; }), 42);
    EXPECT_FALSE(called);
}

TEST(OnceCache, DifferentKeysGetTheirOwnValues) {
    OnceCache<std::string> c;
    EXPECT_EQ(c.get("a", [] { return std::string("apple"); }), "apple");
    EXPECT_EQ(c.get("b", [] { return std::string("banana"); }), "banana");
    EXPECT_EQ(c.get("a", [] { return std::string("avocado"); }), "apple");
    EXPECT_EQ(c.size(), 2u);
}

TEST(OnceCache, ThrowingComputeRemembersNothing) {
    OnceCache<int> c;
    EXPECT_THROW(c.get("a", []() -> int { throw std::runtime_error("boom"); }),
                 std::runtime_error);
    EXPECT_EQ(c.size(), 0u);
    int calls = 0;
    EXPECT_EQ(c.get("a", [&] { ++calls; return 5; }), 5);
    EXPECT_EQ(calls, 1) << "after a failure, the next get() must compute again";
    EXPECT_EQ(c.size(), 1u);
}

TEST(OnceCache, ConcurrentGetsOfOneKeyComputeOnce) {
    // 16 threads miss the same key at the same moment; compute() is slow,
    // so a check-then-act gap is wide open.
    constexpr int kThreads = 16;
    OnceCache<int> c;
    std::atomic<int> calls{0};
    std::vector<int> results(kThreads);
    runTogether(kThreads, [&](int t) {
        results[t] = c.get("config", [&] {
            calls.fetch_add(1);
            std::this_thread::sleep_for(50ms);
            return 99;
        });
    });

    EXPECT_EQ(calls.load(), 1) << "compute() ran more than once for one key";
    for (int r : results) EXPECT_EQ(r, 99);
    EXPECT_EQ(c.size(), 1u);
}

TEST(OnceCache, ConcurrentGetsOfManyKeysComputeEachOnce) {
    constexpr int kThreads = 8;
    constexpr int kKeys = 200;
    OnceCache<int> c;
    std::vector<std::atomic<int>> calls(kKeys);
    runTogether(kThreads, [&](int t) {
        for (int i = 0; i < kKeys; ++i) {
            int k = (i + t * 25) % kKeys;  // each thread starts at a different key
            int v = c.get(std::to_string(k), [&, k] {
                calls[k].fetch_add(1);
                std::this_thread::sleep_for(1ms);
                return k * 10;
            });
            ASSERT_EQ(v, k * 10);
        }
    });

    for (int k = 0; k < kKeys; ++k) {
        EXPECT_EQ(calls[k].load(), 1) << "key " << k << " was computed more than once";
    }
    EXPECT_EQ(c.size(), static_cast<std::size_t>(kKeys));
}

TEST(OnceCache, DifferentKeysDoNotWaitForEachOther) {
    // compute("a") and compute("b") each wait until the other one has
    // started. That only works if both can run at the same time; if a
    // lock is held while compute() runs, one of them gives up after 2 s.
    OnceCache<std::string> c;
    std::promise<void> aStarted, bStarted;
    std::shared_future<void> aStartedF = aStarted.get_future().share();
    std::shared_future<void> bStartedF = bStarted.get_future().share();
    std::string a, b;
    std::thread ta([&] {
        a = c.get("a", [&] {
            aStarted.set_value();
            return bStartedF.wait_for(2s) == std::future_status::ready
                       ? std::string("ok") : std::string("timed out");
        });
    });
    std::thread tb([&] {
        b = c.get("b", [&] {
            bStarted.set_value();
            return aStartedF.wait_for(2s) == std::future_status::ready
                       ? std::string("ok") : std::string("timed out");
        });
    });
    ta.join();
    tb.join();

    EXPECT_EQ(a, "ok") << "compute(\"a\") waited for compute(\"b\")";
    EXPECT_EQ(b, "ok") << "compute(\"b\") waited for compute(\"a\")";
}

TEST(OnceCache, SizeCountsOnlyFinishedValues) {
    OnceCache<int> c;
    std::promise<void> started, release;
    std::future<void> startedF = started.get_future();
    std::shared_future<void> releaseF = release.get_future().share();
    std::thread t([&] {
        c.get("slow", [&] {
            started.set_value();
            releaseF.wait();
            return 1;
        });
    });
    bool didStart = startedF.wait_for(5s) == std::future_status::ready;
    // Ask size() on another thread, so a size() that waits for the running
    // compute() fails this test instead of hanging it.
    auto sizeF = std::async(std::launch::async, [&] { return c.size(); });
    bool sizeAnswered = sizeF.wait_for(2s) == std::future_status::ready;
    release.set_value();
    t.join();
    std::size_t sizeWhileRunning = sizeF.get();

    ASSERT_TRUE(didStart) << "compute() was never called";
    EXPECT_TRUE(sizeAnswered) << "size() waited for a running compute()";
    EXPECT_EQ(sizeWhileRunning, 0u) << "a computation still running counted as remembered";
    EXPECT_EQ(c.size(), 1u);
}
