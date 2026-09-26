#include "day03/bounded_stack.h"
#include "run_together.h"
#include <gtest/gtest.h>

#include <algorithm>
#include <atomic>
#include <chrono>
#include <optional>
#include <thread>
#include <vector>

// The single-threaded contract is the same as Day 2 -- a quick check that
// it still holds.

TEST(BoundedStack, PushPopIsLifo) {
    BoundedStack s(3);
    s.push(1);
    s.push(2);
    s.push(3);
    EXPECT_TRUE(s.full());
    EXPECT_EQ(s.top(), 3);
    EXPECT_EQ(s.pop(), 3);
    EXPECT_EQ(s.pop(), 2);
    EXPECT_EQ(s.size(), 1u);
}

TEST(BoundedStack, FullAndEmptyStillThrow) {
    BoundedStack s(1);
    EXPECT_TRUE(s.empty());
    EXPECT_THROW(s.pop(), std::runtime_error);
    EXPECT_THROW(s.top(), std::runtime_error);
    s.push(7);
    EXPECT_THROW(s.push(8), std::runtime_error);
    EXPECT_EQ(s.size(), 1u);
    EXPECT_EQ(s.top(), 7);
}

TEST(BoundedStack, ConcurrentPushesStopExactlyAtCapacity) {
    // 8 threads push 2,000 values into a stack that holds 1,000.
    // Exactly 1,000 pushes may succeed, and the stack must then hold
    // exactly the values whose push succeeded. The race only matters at
    // the moment the stack fills up, so play it 20 times.
    constexpr int kThreads = 8;
    constexpr int kPerThread = 250;
    for (int round = 0; round < 20; ++round) {
        BoundedStack s(1000);
        std::vector<std::vector<int>> pushed(kThreads);
        runTogether(kThreads, [&](int t) {
            for (int i = 0; i < kPerThread; ++i) {
                int v = t * 1000 + i;
                try {
                    s.push(v);
                    pushed[t].push_back(v);
                } catch (const std::runtime_error&) {
                }
            }
        });

        std::vector<int> expected;
        for (const auto& p : pushed) expected.insert(expected.end(), p.begin(), p.end());
        ASSERT_EQ(expected.size(), 1000u) << "round " << round << ": pushes that succeeded";
        ASSERT_EQ(s.size(), 1000u) << "round " << round;
        ASSERT_TRUE(s.full()) << "round " << round;

        std::vector<int> held;
        for (int i = 0; i < 1000 && !s.empty(); ++i) held.push_back(s.pop());
        ASSERT_TRUE(s.empty()) << "round " << round;
        std::sort(expected.begin(), expected.end());
        std::sort(held.begin(), held.end());
        ASSERT_EQ(held, expected) << "round " << round << ": the stack lost a value or made one up";
    }
}

TEST(BoundedStack, ConcurrentPopsReturnEachValueOnce) {
    constexpr int kValues = 20000;
    constexpr int kThreads = 8;
    BoundedStack s(kValues);
    for (int v = 0; v < kValues; ++v) s.push(v);

    // Pop until the stack says it is empty. (Stop early if more values come
    // out than ever went in -- the stack is making them up.)
    std::vector<std::vector<int>> popped(kThreads);
    std::atomic<int> poppedCount{0};
    runTogether(kThreads, [&](int t) {
        while (poppedCount.load() <= kValues) {
            try {
                popped[t].push_back(s.pop());
                poppedCount.fetch_add(1);
            } catch (const std::runtime_error&) {
                return;
            }
        }
    });

    std::vector<int> all;
    for (const auto& p : popped) all.insert(all.end(), p.begin(), p.end());
    std::sort(all.begin(), all.end());
    ASSERT_EQ(all.size(), static_cast<std::size_t>(kValues))
        << "a value was popped twice, or lost";
    for (int v = 0; v < kValues; ++v) {
        ASSERT_EQ(all[v], v) << "a value was popped twice, or lost";
    }
    EXPECT_TRUE(s.empty());
}

TEST(BoundedStack, ConcurrentPushAndPopConserveValues) {
    // A small stack, 4 pushers and 4 poppers. Pushers retry while it is
    // full, poppers retry while it is empty. Every value pushed must come
    // out exactly once.
    constexpr int kPushers = 4;
    constexpr int kPerPusher = 5000;
    constexpr int kTotal = kPushers * kPerPusher;
    BoundedStack s(16);
    std::atomic<int> poppedCount{0};
    std::vector<std::vector<int>> popped(kPushers);
    const auto deadline = std::chrono::steady_clock::now() + std::chrono::seconds(20);
    auto pastDeadline = [&] { return std::chrono::steady_clock::now() > deadline; };

    runTogether(2 * kPushers, [&](int t) {
        if (t < kPushers) {
            for (int i = 0; i < kPerPusher && !pastDeadline(); ++i) {
                for (;;) {
                    try {
                        s.push(t * kPerPusher + i);
                        break;
                    } catch (const std::runtime_error&) {
                        if (pastDeadline()) return;
                        std::this_thread::yield();
                    }
                }
            }
        } else {
            auto& mine = popped[t - kPushers];
            while (poppedCount.load() < kTotal && !pastDeadline()) {
                try {
                    mine.push_back(s.pop());
                    poppedCount.fetch_add(1);
                } catch (const std::runtime_error&) {
                    std::this_thread::yield();
                }
            }
        }
    });

    ASSERT_FALSE(pastDeadline()) << "gave up after 20 s -- values were lost";
    std::vector<int> all;
    for (const auto& p : popped) all.insert(all.end(), p.begin(), p.end());
    std::sort(all.begin(), all.end());
    ASSERT_EQ(all.size(), static_cast<std::size_t>(kTotal));
    for (int v = 0; v < kTotal; ++v) {
        ASSERT_EQ(all[v], v) << "a value was popped twice, or lost";
    }
    EXPECT_TRUE(s.empty());
}

// ---- tryPush / tryPop: check and act as one step ----

TEST(BoundedStack, TryPushAndTryPopFollowTheRules) {
    BoundedStack s(2);
    EXPECT_EQ(s.tryPop(), std::nullopt) << "empty stack";
    EXPECT_TRUE(s.tryPush(1));
    EXPECT_TRUE(s.tryPush(2));
    EXPECT_FALSE(s.tryPush(3)) << "full stack";
    EXPECT_EQ(s.size(), 2u);
    EXPECT_EQ(s.tryPop(), std::optional<int>(2));
    EXPECT_EQ(s.tryPop(), std::optional<int>(1));
    EXPECT_EQ(s.tryPop(), std::nullopt);
    EXPECT_TRUE(s.empty());
}

TEST(BoundedStack, ConcurrentTryPopTakesEachValueOnce) {
    // 8 threads empty the stack with tryPop -- the one-call version of
    // "if (!s.empty()) { v = s.top(); s.pop(); }".
    constexpr int kValues = 20000;
    constexpr int kThreads = 8;
    BoundedStack s(kValues);
    for (int v = 0; v < kValues; ++v) s.push(v);

    std::vector<std::vector<int>> popped(kThreads);
    std::atomic<int> poppedCount{0};
    runTogether(kThreads, [&](int t) {
        while (poppedCount.load() <= kValues) {  // more than kValues: made up
            std::optional<int> v = s.tryPop();
            if (!v) return;
            popped[t].push_back(*v);
            poppedCount.fetch_add(1);
        }
    });

    std::vector<int> all;
    for (const auto& p : popped) all.insert(all.end(), p.begin(), p.end());
    std::sort(all.begin(), all.end());
    ASSERT_EQ(all.size(), static_cast<std::size_t>(kValues))
        << "a value was popped twice, or lost";
    for (int v = 0; v < kValues; ++v) {
        ASSERT_EQ(all[v], v) << "a value was popped twice, or lost";
    }
    EXPECT_TRUE(s.empty());
}

TEST(BoundedStack, ConcurrentTryPushStopsExactlyAtCapacity) {
    // Like ConcurrentPushesStopExactlyAtCapacity, with tryPush.
    constexpr int kThreads = 8;
    constexpr int kPerThread = 250;
    for (int round = 0; round < 20; ++round) {
        BoundedStack s(1000);
        std::vector<std::vector<int>> pushed(kThreads);
        runTogether(kThreads, [&](int t) {
            for (int i = 0; i < kPerThread; ++i) {
                int v = t * 1000 + i;
                if (s.tryPush(v)) pushed[t].push_back(v);
            }
        });

        std::vector<int> expected;
        for (const auto& p : pushed) expected.insert(expected.end(), p.begin(), p.end());
        ASSERT_EQ(expected.size(), 1000u) << "round " << round << ": tryPush calls that returned true";
        ASSERT_EQ(s.size(), 1000u) << "round " << round;

        std::vector<int> held;
        for (int i = 0; i < 1000; ++i) {
            std::optional<int> v = s.tryPop();
            if (!v) break;
            held.push_back(*v);
        }
        std::sort(expected.begin(), expected.end());
        std::sort(held.begin(), held.end());
        ASSERT_EQ(held, expected) << "round " << round << ": the stack lost a value or made one up";
    }
}

TEST(BoundedStack, ConcurrentTryPushAndTryPopConserveValues) {
    // A small stack, 4 pushers and 4 poppers, all using the try- calls.
    constexpr int kPushers = 4;
    constexpr int kPerPusher = 5000;
    constexpr int kTotal = kPushers * kPerPusher;
    BoundedStack s(16);
    std::atomic<int> poppedCount{0};
    std::vector<std::vector<int>> popped(kPushers);
    const auto deadline = std::chrono::steady_clock::now() + std::chrono::seconds(10);
    auto pastDeadline = [&] { return std::chrono::steady_clock::now() > deadline; };

    runTogether(2 * kPushers, [&](int t) {
        if (t < kPushers) {
            for (int i = 0; i < kPerPusher && !pastDeadline(); ++i) {
                while (!s.tryPush(t * kPerPusher + i)) {
                    if (pastDeadline()) return;
                    std::this_thread::yield();
                }
            }
        } else {
            auto& mine = popped[t - kPushers];
            while (poppedCount.load() < kTotal && !pastDeadline()) {
                if (std::optional<int> v = s.tryPop()) {
                    mine.push_back(*v);
                    poppedCount.fetch_add(1);
                } else {
                    std::this_thread::yield();
                }
            }
        }
    });

    ASSERT_FALSE(pastDeadline()) << "gave up after 10 s -- values were lost";
    std::vector<int> all;
    for (const auto& p : popped) all.insert(all.end(), p.begin(), p.end());
    std::sort(all.begin(), all.end());
    ASSERT_EQ(all.size(), static_cast<std::size_t>(kTotal));
    for (int v = 0; v < kTotal; ++v) {
        ASSERT_EQ(all[v], v) << "a value was popped twice, or lost";
    }
    EXPECT_TRUE(s.empty());
}
