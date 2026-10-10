#include "day04/bounded_blocking_queue.h"
#include "run_together.h"
#include <gtest/gtest.h>

#include <algorithm>
#include <atomic>
#include <chrono>
#include <memory>
#include <optional>
#include <stdexcept>
#include <thread>
#include <vector>

using namespace std::chrono_literals;
using Clock = std::chrono::steady_clock;

namespace {
long long millis(Clock::duration d) {
    return std::chrono::duration_cast<std::chrono::milliseconds>(d).count();
}
}  // namespace

// ---- the single-threaded contract ----

TEST(BoundedBlockingQueue, ZeroCapacityIsRejected) {
    EXPECT_THROW({ BoundedBlockingQueue<int> q(0); }, std::invalid_argument);
}

TEST(BoundedBlockingQueue, ItemsComeOutInTheOrderTheyWentIn) {
    BoundedBlockingQueue<int> q(3);
    EXPECT_EQ(q.capacity(), 3u);
    EXPECT_EQ(q.size(), 0u);
    EXPECT_TRUE(q.push(1));
    EXPECT_TRUE(q.push(2));
    EXPECT_TRUE(q.push(3));
    EXPECT_EQ(q.size(), 3u);
    EXPECT_EQ(q.pop(), std::optional<int>(1));
    EXPECT_EQ(q.pop(), std::optional<int>(2));
    EXPECT_EQ(q.pop(), std::optional<int>(3));
    EXPECT_EQ(q.size(), 0u);
}

TEST(BoundedBlockingQueue, OrderSurvivesManyRoundsThroughASmallQueue) {
    BoundedBlockingQueue<int> q(3);
    int next = 0, expected = 0;
    for (int round = 0; round < 200; ++round) {
        for (int i = 0; i < 2; ++i) ASSERT_TRUE(q.push(next++));
        ASSERT_EQ(q.pop(), std::optional<int>(expected++));
        ASSERT_EQ(q.pop(), std::optional<int>(expected++));
    }
}

TEST(BoundedBlockingQueue, TryPushAndTryPopNeverWait) {
    BoundedBlockingQueue<int> q(2);
    EXPECT_EQ(q.tryPop(), std::nullopt) << "empty queue";
    EXPECT_TRUE(q.tryPush(1));
    EXPECT_TRUE(q.tryPush(2));
    EXPECT_FALSE(q.tryPush(3)) << "full queue";
    EXPECT_EQ(q.size(), 2u);
    EXPECT_EQ(q.tryPop(), std::optional<int>(1));
    EXPECT_EQ(q.tryPop(), std::optional<int>(2));
    EXPECT_EQ(q.tryPop(), std::nullopt);
}

TEST(BoundedBlockingQueue, WorksWithMoveOnlyValues) {
    BoundedBlockingQueue<std::unique_ptr<int>> q(2);
    EXPECT_TRUE(q.push(std::make_unique<int>(5)));
    EXPECT_TRUE(q.tryPush(std::make_unique<int>(6)));
    auto a = q.pop();
    ASSERT_TRUE(a.has_value());
    ASSERT_NE(*a, nullptr);
    EXPECT_EQ(**a, 5);
    auto b = q.popFor(100ms);
    ASSERT_TRUE(b.has_value());
    ASSERT_NE(*b, nullptr);
    EXPECT_EQ(**b, 6);
}

// ---- blocking ----

TEST(BoundedBlockingQueue, PushWaitsWhileTheQueueIsFull) {
    BoundedBlockingQueue<int> q(2);
    ASSERT_TRUE(q.push(1));
    ASSERT_TRUE(q.push(2));

    Background<bool> producer([&] { return q.push(3); });
    EXPECT_FALSE(producer.waitFor(150ms)) << "push() returned although the queue was full";
    EXPECT_EQ(q.size(), 2u);

    EXPECT_EQ(q.pop(), std::optional<int>(1));
    EXPECT_TRUE(producer.waitFor(3s)) << "push() was not woken up when room appeared";
    EXPECT_TRUE(producer.result());
    EXPECT_EQ(q.pop(), std::optional<int>(2));
    EXPECT_EQ(q.pop(), std::optional<int>(3));
}

TEST(BoundedBlockingQueue, PopWaitsWhileTheQueueIsEmpty) {
    BoundedBlockingQueue<int> q(2);

    Background<std::optional<int>> consumer([&] { return q.pop(); });
    EXPECT_FALSE(consumer.waitFor(150ms)) << "pop() returned although the queue was empty";

    EXPECT_TRUE(q.push(7));
    EXPECT_TRUE(consumer.waitFor(3s)) << "pop() was not woken up when an item arrived";
    EXPECT_EQ(consumer.result(), std::optional<int>(7));
}

TEST(BoundedBlockingQueue, BackpressureHoldsTheProducerBack) {
    // One producer, no consumer yet: it must stop after `capacity` items,
    // not run ahead.
    BoundedBlockingQueue<int> q(4);
    Background<bool> producer([&] {
        for (int i = 0; i < 1000; ++i) {
            if (!q.push(i)) return false;
        }
        return true;
    });
    EXPECT_FALSE(producer.waitFor(200ms)) << "the producer ran past the capacity";
    EXPECT_EQ(q.size(), 4u);

    for (int i = 0; i < 1000; ++i) {
        std::optional<int> v = q.pop();
        ASSERT_TRUE(v.has_value());
        ASSERT_EQ(*v, i);
    }
    EXPECT_TRUE(producer.result());
}

// ---- timed calls ----

TEST(BoundedBlockingQueue, PopForGivesUpAfterTheTimeout) {
    BoundedBlockingQueue<int> q(2);
    Background<std::optional<int>> consumer([&] { return q.popFor(100ms); });
    const auto start = Clock::now();
    std::optional<int> got = consumer.result(10);
    const auto waited = millis(Clock::now() - start);
    EXPECT_EQ(got, std::nullopt);
    EXPECT_GE(waited, 80) << "popFor(100ms) gave up too early";
    EXPECT_LT(waited, 3000) << "popFor(100ms) waited far too long";
}

TEST(BoundedBlockingQueue, PushForGivesUpAfterTheTimeout) {
    BoundedBlockingQueue<int> q(1);
    ASSERT_TRUE(q.push(1));
    Background<bool> producer([&] { return q.pushFor(2, 100ms); });
    const auto start = Clock::now();
    bool pushed = producer.result(10);
    const auto waited = millis(Clock::now() - start);
    EXPECT_FALSE(pushed);
    EXPECT_GE(waited, 80) << "pushFor(100ms) gave up too early";
    EXPECT_LT(waited, 3000) << "pushFor(100ms) waited far too long";
    EXPECT_EQ(q.size(), 1u);
    EXPECT_EQ(q.pop(), std::optional<int>(1)) << "a timed-out push must add nothing";
    EXPECT_EQ(q.tryPop(), std::nullopt);
}

TEST(BoundedBlockingQueue, PopForReturnsAsSoonAsAnItemArrives) {
    BoundedBlockingQueue<int> q(2);
    Background<std::optional<int>> consumer([&] { return q.popFor(20s); });
    std::this_thread::sleep_for(50ms);
    EXPECT_TRUE(q.push(9));
    EXPECT_TRUE(consumer.waitFor(3s)) << "popFor kept waiting although an item had arrived";
    EXPECT_EQ(consumer.result(), std::optional<int>(9));
}

TEST(BoundedBlockingQueue, PushForReturnsAsSoonAsRoomAppears) {
    BoundedBlockingQueue<int> q(1);
    ASSERT_TRUE(q.push(1));
    Background<bool> producer([&] { return q.pushFor(2, 20s); });
    std::this_thread::sleep_for(50ms);
    EXPECT_EQ(q.pop(), std::optional<int>(1));
    EXPECT_TRUE(producer.waitFor(3s)) << "pushFor kept waiting although room had appeared";
    EXPECT_TRUE(producer.result());
    EXPECT_EQ(q.pop(), std::optional<int>(2));
}

// ---- close ----

TEST(BoundedBlockingQueue, CloseWakesEverySleepingConsumer) {
    BoundedBlockingQueue<int> q(2);
    std::vector<std::unique_ptr<Background<std::optional<int>>>> consumers;
    for (int i = 0; i < 4; ++i) {
        consumers.push_back(std::make_unique<Background<std::optional<int>>>([&] { return q.pop(); }));
    }
    std::this_thread::sleep_for(100ms);
    for (auto& c : consumers) EXPECT_FALSE(c->waitFor(0ms)) << "pop() returned on an empty, open queue";
    q.close();
    for (auto& c : consumers) {
        EXPECT_TRUE(c->waitFor(3s)) << "a consumer is still asleep after close()";
        EXPECT_EQ(c->result(), std::nullopt);
    }
}

TEST(BoundedBlockingQueue, CloseWakesEverySleepingProducer) {
    BoundedBlockingQueue<int> q(1);
    ASSERT_TRUE(q.push(100));
    std::vector<std::unique_ptr<Background<bool>>> producers;
    for (int i = 0; i < 4; ++i) {
        producers.push_back(std::make_unique<Background<bool>>([&, i] { return q.push(i); }));
    }
    std::this_thread::sleep_for(100ms);
    q.close();
    for (auto& p : producers) {
        EXPECT_TRUE(p->waitFor(3s)) << "a producer is still asleep after close()";
        EXPECT_FALSE(p->result()) << "push() on a closed queue must return false";
    }
    EXPECT_EQ(q.size(), 1u) << "nothing may be added to a closed queue";
    EXPECT_EQ(q.pop(), std::optional<int>(100));
    EXPECT_EQ(q.pop(), std::nullopt);
}

TEST(BoundedBlockingQueue, CloseWakesTimedWaitersToo) {
    BoundedBlockingQueue<int> empty(2);
    BoundedBlockingQueue<int> full(1);
    ASSERT_TRUE(full.push(1));
    Background<std::optional<int>> consumer([&] { return empty.popFor(30s); });
    Background<bool> producer([&] { return full.pushFor(2, 30s); });
    std::this_thread::sleep_for(100ms);
    empty.close();
    full.close();
    EXPECT_TRUE(consumer.waitFor(3s)) << "popFor is still asleep after close()";
    EXPECT_TRUE(producer.waitFor(3s)) << "pushFor is still asleep after close()";
    EXPECT_EQ(consumer.result(), std::nullopt);
    EXPECT_FALSE(producer.result());
}

TEST(BoundedBlockingQueue, ClosedQueueHandsOutWhatItHoldsThenStops) {
    BoundedBlockingQueue<int> q(3);
    ASSERT_TRUE(q.push(1));
    ASSERT_TRUE(q.push(2));
    ASSERT_TRUE(q.push(3));
    EXPECT_FALSE(q.closed());
    q.close();
    q.close();  // again: nothing happens
    EXPECT_TRUE(q.closed());

    // Full and closed: push must not wait for room that will never matter.
    Background<bool> late([&] { return q.push(4); });
    EXPECT_TRUE(late.waitFor(2s)) << "push() on a closed, full queue must return at once";
    EXPECT_FALSE(late.result());
    EXPECT_FALSE(q.tryPush(5));
    EXPECT_FALSE(q.pushFor(6, 10ms));

    EXPECT_EQ(q.pop(), std::optional<int>(1));
    EXPECT_EQ(q.tryPop(), std::optional<int>(2));
    EXPECT_EQ(q.popFor(100ms), std::optional<int>(3));

    // Closed and empty: every pop returns at once, even a long popFor.
    Background<std::optional<int>> a([&] { return q.pop(); });
    Background<std::optional<int>> b([&] { return q.popFor(20s); });
    EXPECT_TRUE(a.waitFor(2s)) << "pop() on a closed, empty queue must return at once";
    EXPECT_TRUE(b.waitFor(2s)) << "popFor() on a closed, empty queue must return at once";
    EXPECT_EQ(a.result(), std::nullopt);
    EXPECT_EQ(b.result(), std::nullopt);
}

// ---- many threads ----

namespace {

// 4 producers and 4 consumers share one queue. Every value must come out
// exactly once, in order for each producer, and the queue must never hold
// more than its capacity. A small capacity forces constant waiting.
void conserveValues(std::size_t capacity) {
    constexpr int kProducers = 4;
    constexpr int kConsumers = 4;
    constexpr int kPerProducer = 5000;
    constexpr int kTotal = kProducers * kPerProducer;
    BoundedBlockingQueue<int> q(capacity);
    std::atomic<int> producersLeft{kProducers};
    std::atomic<bool> overfull{false};
    std::atomic<bool> pushFailed{false};
    std::vector<std::vector<int>> got(kConsumers);

    runTogetherWithin(30, kProducers + kConsumers, [&](int t) {
        if (t < kProducers) {
            for (int i = 0; i < kPerProducer; ++i) {
                if (!q.push(t * kPerProducer + i)) {
                    pushFailed = true;
                    break;
                }
                if (q.size() > capacity) overfull = true;
            }
            if (producersLeft.fetch_sub(1) == 1) q.close();  // the last one out
        } else {
            auto& mine = got[t - kProducers];
            while (std::optional<int> v = q.pop()) mine.push_back(*v);
        }
    });

    EXPECT_FALSE(pushFailed.load()) << "push() returned false on a queue that was not closed";
    EXPECT_FALSE(overfull.load()) << "size() went above capacity";
    std::vector<int> all;
    for (const auto& mine : got) {
        std::vector<int> last(kProducers, -1);
        for (int v : mine) {
            int p = v / kPerProducer;
            ASSERT_GT(v, last[p]) << "items of one producer came out of order";
            last[p] = v;
        }
        all.insert(all.end(), mine.begin(), mine.end());
    }
    std::sort(all.begin(), all.end());
    ASSERT_EQ(all.size(), static_cast<std::size_t>(kTotal)) << "a value was lost or made up";
    for (int v = 0; v < kTotal; ++v) {
        ASSERT_EQ(all[v], v) << "a value was popped twice, or lost";
    }
}

}  // namespace

TEST(BoundedBlockingQueue, ManyProducersAndConsumersConserveValues) {
    conserveValues(8);
}

TEST(BoundedBlockingQueue, ManyProducersAndConsumersThroughASlotOfOne) {
    // Capacity 1: every push must wake a consumer and every pop a
    // producer. One lost notification and a thread sleeps forever.
    conserveValues(1);
}

TEST(BoundedBlockingQueue, ConcurrentTryPopTakesEachValueOnce) {
    constexpr int kValues = 20000;
    constexpr int kThreads = 8;
    BoundedBlockingQueue<int> q(kValues);
    for (int v = 0; v < kValues; ++v) ASSERT_TRUE(q.push(v));

    std::vector<std::vector<int>> popped(kThreads);
    std::atomic<int> poppedCount{0};
    runTogether(kThreads, [&](int t) {
        while (poppedCount.load() <= kValues) {  // more than kValues: made up
            std::optional<int> v = q.tryPop();
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
    EXPECT_EQ(q.size(), 0u);
}
