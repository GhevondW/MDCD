#include "day03/event_bus.h"
#include "run_together.h"
#include <gtest/gtest.h>

#include <algorithm>
#include <atomic>
#include <chrono>
#include <future>
#include <set>
#include <string>
#include <thread>
#include <vector>

using namespace std::chrono_literals;

TEST(EventBus, PublishCallsSubscribersInOrder) {
    EventBus bus;
    std::vector<std::string> log;
    bus.subscribe([&](const std::string& e) { log.push_back("A:" + e); });
    bus.subscribe([&](const std::string& e) { log.push_back("B:" + e); });
    bus.subscribe([&](const std::string& e) { log.push_back("C:" + e); });
    EXPECT_EQ(bus.publish("x"), 3u);
    EXPECT_EQ(log, (std::vector<std::string>{"A:x", "B:x", "C:x"}));
}

TEST(EventBus, PublishWithNoSubscribersCallsNobody) {
    EventBus bus;
    EXPECT_EQ(bus.publish("x"), 0u);
    EXPECT_EQ(bus.subscriberCount(), 0u);
}

TEST(EventBus, IdsArePositiveAndUnique) {
    EventBus bus;
    std::set<long long> ids;
    for (int i = 0; i < 100; ++i) {
        long long id = bus.subscribe([](const std::string&) {});
        EXPECT_GT(id, 0);
        ids.insert(id);
    }
    EXPECT_EQ(ids.size(), 100u);
    EXPECT_EQ(bus.subscriberCount(), 100u);
}

TEST(EventBus, UnsubscribeStopsDelivery) {
    EventBus bus;
    std::vector<std::string> log;
    bus.subscribe([&](const std::string&) { log.push_back("A"); });
    long long b = bus.subscribe([&](const std::string&) { log.push_back("B"); });
    bus.subscribe([&](const std::string&) { log.push_back("C"); });

    EXPECT_TRUE(bus.unsubscribe(b));
    EXPECT_FALSE(bus.unsubscribe(b)) << "already removed";
    EXPECT_FALSE(bus.unsubscribe(12345)) << "never existed";
    EXPECT_EQ(bus.subscriberCount(), 2u);
    EXPECT_EQ(bus.publish("x"), 2u);
    EXPECT_EQ(log, (std::vector<std::string>{"A", "C"}));
}

TEST(EventBus, OrderIsKeptAfterUnsubscribeAndSubscribe) {
    EventBus bus;
    std::vector<std::string> log;
    bus.subscribe([&](const std::string&) { log.push_back("A"); });
    long long b = bus.subscribe([&](const std::string&) { log.push_back("B"); });
    bus.subscribe([&](const std::string&) { log.push_back("C"); });
    bus.unsubscribe(b);
    bus.subscribe([&](const std::string&) { log.push_back("D"); });
    bus.publish("x");
    EXPECT_EQ(log, (std::vector<std::string>{"A", "C", "D"}));
}

TEST(EventBus, SubscribeInsideACallbackCountsFromTheNextPublish) {
    EventBus bus;
    std::vector<std::string> log;
    bool added = false;
    bus.subscribe([&](const std::string&) {
        log.push_back("A");
        if (!added) {
            added = true;
            bus.subscribe([&](const std::string&) { log.push_back("B"); });
        }
    });

    EXPECT_EQ(bus.publish("first"), 1u) << "B joined mid-publish: not called by it";
    EXPECT_EQ(log, (std::vector<std::string>{"A"}));
    log.clear();
    EXPECT_EQ(bus.publish("second"), 2u);
    EXPECT_EQ(log, (std::vector<std::string>{"A", "B"}));
}

TEST(EventBus, UnsubscribeInsideACallbackCountsFromTheNextPublish) {
    EventBus bus;
    std::vector<std::string> log;
    long long b = 0;
    bool removed = false;
    bus.subscribe([&](const std::string&) {
        log.push_back("A");
        if (!removed) {
            removed = true;
            EXPECT_TRUE(bus.unsubscribe(b));
        }
    });
    b = bus.subscribe([&](const std::string&) { log.push_back("B"); });

    EXPECT_EQ(bus.publish("first"), 2u) << "B left mid-publish: still called by it";
    EXPECT_EQ(log, (std::vector<std::string>{"A", "B"}));
    log.clear();
    EXPECT_EQ(bus.publish("second"), 1u);
    EXPECT_EQ(log, (std::vector<std::string>{"A"}));
}

TEST(EventBus, CallbackCanUnsubscribeItself) {
    EventBus bus;
    int aCalls = 0;
    long long a = 0;
    a = bus.subscribe([&](const std::string&) {
        ++aCalls;
        bus.unsubscribe(a);
    });
    bus.subscribe([](const std::string&) {});

    EXPECT_EQ(bus.publish("first"), 2u);
    EXPECT_EQ(bus.publish("second"), 1u);
    EXPECT_EQ(aCalls, 1);
}

TEST(EventBus, CallbackCanPublish) {
    EventBus bus;
    std::vector<std::string> bLog;
    bus.subscribe([&](const std::string& e) {
        if (e == "outer") bus.publish("inner");
    });
    bus.subscribe([&](const std::string& e) { bLog.push_back(e); });

    EXPECT_EQ(bus.publish("outer"), 2u);
    // A handles "outer" by publishing "inner", which reaches B first.
    EXPECT_EQ(bLog, (std::vector<std::string>{"inner", "outer"}));
}

TEST(EventBus, OtherThreadsCanUseTheBusWhileACallbackRuns) {
    // The callback waits (up to 2 s) for another thread to subscribe. If
    // publish() holds a lock while calling back, that thread is stuck.
    EventBus bus;
    std::string outcome;
    std::thread helper;
    bus.subscribe([&](const std::string&) {
        std::promise<void> done;
        std::future<void> doneF = done.get_future();
        helper = std::thread([&bus, p = std::move(done)]() mutable {
            bus.subscribe([](const std::string&) {});
            p.set_value();
        });
        outcome = doneF.wait_for(2s) == std::future_status::ready ? "ok" : "blocked";
    });

    bus.publish("x");
    if (helper.joinable()) helper.join();
    EXPECT_EQ(outcome, "ok") << "another thread could not subscribe while a callback ran";
    EXPECT_EQ(bus.subscriberCount(), 2u);
}

TEST(EventBus, ConcurrentSubscribesGetUniqueIds) {
    constexpr int kThreads = 8;
    constexpr int kPerThread = 1000;
    EventBus bus;
    std::vector<std::vector<long long>> ids(kThreads);
    runTogether(kThreads, [&](int t) {
        for (int i = 0; i < kPerThread; ++i) {
            ids[t].push_back(bus.subscribe([](const std::string&) {}));
        }
    });

    std::set<long long> all;
    for (const auto& v : ids) all.insert(v.begin(), v.end());
    EXPECT_EQ(all.size(), static_cast<std::size_t>(kThreads * kPerThread));
    EXPECT_EQ(bus.subscriberCount(), static_cast<std::size_t>(kThreads * kPerThread));
}

TEST(EventBus, ConcurrentPublishAndChurn) {
    // 3 subscribers stay for the whole test and must hear every event
    // exactly once, while other threads keep subscribing and
    // unsubscribing around them.
    constexpr int kPublishers = 4;
    constexpr int kChurners = 4;
    constexpr int kPerThread = 2000;
    EventBus bus;
    std::atomic<int> heard[3] = {{0}, {0}, {0}};
    for (auto& h : heard) {
        bus.subscribe([&h](const std::string&) { h.fetch_add(1); });
    }
    std::atomic<bool> badCount{false};
    std::atomic<bool> badUnsubscribe{false};

    runTogether(kPublishers + kChurners, [&](int t) {
        for (int i = 0; i < kPerThread; ++i) {
            if (t < kPublishers) {
                if (bus.publish("tick") < 3) badCount = true;
            } else {
                long long id = bus.subscribe([](const std::string&) {});
                if (!bus.unsubscribe(id)) badUnsubscribe = true;
            }
        }
    });

    for (auto& h : heard) EXPECT_EQ(h.load(), kPublishers * kPerThread);
    EXPECT_FALSE(badCount) << "a publish missed a permanent subscriber";
    EXPECT_FALSE(badUnsubscribe) << "unsubscribe of a fresh id returned false";
    EXPECT_EQ(bus.subscriberCount(), 3u);
}
