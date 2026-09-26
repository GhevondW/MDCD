#include "day03/ticket_dispenser.h"
#include "run_together.h"
#include <gtest/gtest.h>

#include <algorithm>
#include <atomic>
#include <vector>

TEST(TicketDispenser, FirstTicketIsOne) {
    TicketDispenser d;
    EXPECT_EQ(d.next(), 1);
}

TEST(TicketDispenser, TicketsCountUpByOne) {
    TicketDispenser d;
    EXPECT_EQ(d.next(), 1);
    EXPECT_EQ(d.next(), 2);
    EXPECT_EQ(d.next(), 3);
}

TEST(TicketDispenser, IssuedCountsHandedOutTickets) {
    TicketDispenser d;
    EXPECT_EQ(d.issued(), 0);
    d.next();
    d.next();
    EXPECT_EQ(d.issued(), 2);
}

TEST(TicketDispenser, ConcurrentTicketsAreUniqueAndGapless) {
    constexpr int kThreads = 8;
    constexpr int kPerThread = 20000;
    TicketDispenser d;
    std::vector<std::vector<long long>> got(kThreads);
    runTogether(kThreads, [&](int t) {
        for (int i = 0; i < kPerThread; ++i) got[t].push_back(d.next());
    });

    std::vector<long long> all;
    for (const auto& g : got) all.insert(all.end(), g.begin(), g.end());
    std::sort(all.begin(), all.end());
    for (std::size_t i = 0; i < all.size(); ++i) {
        ASSERT_EQ(all[i], static_cast<long long>(i) + 1)
            << "a ticket was handed out twice, or a number was skipped";
    }
    EXPECT_EQ(d.issued(), kThreads * kPerThread);
}

TEST(TicketDispenser, EachThreadSeesItsTicketsIncrease) {
    constexpr int kThreads = 8;
    constexpr int kPerThread = 20000;
    TicketDispenser d;
    std::vector<std::vector<long long>> got(kThreads);
    runTogether(kThreads, [&](int t) {
        for (int i = 0; i < kPerThread; ++i) got[t].push_back(d.next());
    });

    for (const auto& g : got) {
        for (std::size_t i = 1; i < g.size(); ++i) {
            ASSERT_LT(g[i - 1], g[i]) << "a later ticket had a smaller number";
        }
    }
}

TEST(TicketDispenser, IssuedNeverGoesBackwardsWhileTicketsAreTaken) {
    constexpr int kTakers = 4;
    constexpr int kPerThread = 20000;
    TicketDispenser d;
    std::atomic<bool> wentBackwards{false};
    runTogether(kTakers + 1, [&](int t) {
        if (t == kTakers) {  // the watcher
            long long last = 0;
            for (int i = 0; i < 20000; ++i) {
                long long now = d.issued();
                if (now < last) wentBackwards = true;
                last = now;
            }
            return;
        }
        for (int i = 0; i < kPerThread; ++i) d.next();
    });

    EXPECT_FALSE(wentBackwards);
    EXPECT_EQ(d.issued(), kTakers * kPerThread);
}
