#include "day03/concurrent_booker.h"
#include "run_together.h"
#include <gtest/gtest.h>

#include <algorithm>
#include <atomic>
#include <random>
#include <vector>

namespace {

// No two bookings in a sorted list may overlap.
bool noOverlaps(const std::vector<Interval>& sorted) {
    for (std::size_t i = 1; i < sorted.size(); ++i) {
        if (sorted[i].start < sorted[i - 1].end) return false;
    }
    return true;
}

bool sortedByStart(const std::vector<Interval>& v) {
    return std::is_sorted(v.begin(), v.end(),
                          [](const Interval& a, const Interval& b) { return a.start < b.start; });
}

}  // namespace

// ---- single thread: the rules ----

TEST(ConcurrentBooker, BookRejectsOverlapAndAcceptsAdjacent) {
    ConcurrentBooker b;
    EXPECT_TRUE(b.book(10, 20));
    EXPECT_TRUE(b.book(20, 30));   // adjacent: [10,20) and [20,30) don't overlap
    EXPECT_FALSE(b.book(15, 25));
    EXPECT_FALSE(b.book(0, 11));
    EXPECT_EQ(b.count(), 2u);
}

TEST(ConcurrentBooker, EmptyOrReversedIntervalThrows) {
    ConcurrentBooker b;
    EXPECT_THROW(b.book(5, 5), std::invalid_argument);
    EXPECT_THROW(b.book(6, 5), std::invalid_argument);
    EXPECT_EQ(b.count(), 0u);
}

TEST(ConcurrentBooker, CancelNeedsAnExactMatch) {
    ConcurrentBooker b;
    ASSERT_TRUE(b.book(10, 20));
    EXPECT_FALSE(b.cancel(10, 15));
    EXPECT_TRUE(b.cancel(10, 20));
    EXPECT_FALSE(b.cancel(10, 20));
    EXPECT_EQ(b.count(), 0u);
}

TEST(ConcurrentBooker, BookingsAreSortedByStart) {
    ConcurrentBooker b;
    b.book(50, 60);
    b.book(0, 10);
    b.book(20, 30);
    EXPECT_EQ(b.bookings(), (std::vector<Interval>{{0, 10}, {20, 30}, {50, 60}}));
}

TEST(ConcurrentBooker, BookFirstFreeOnEmptyCalendarStartsAtFrom) {
    ConcurrentBooker b;
    EXPECT_EQ(b.bookFirstFree(7, 5), 7);
    EXPECT_EQ(b.bookings(), (std::vector<Interval>{{7, 12}}));
}

TEST(ConcurrentBooker, BookFirstFreeSkipsBusyRangesAndSmallGaps) {
    ConcurrentBooker b;
    b.book(0, 10);
    b.book(12, 20);
    EXPECT_EQ(b.bookFirstFree(0, 5), 20);  // the gap [10,12) is too small
    EXPECT_EQ(b.bookFirstFree(0, 2), 10);  // ...but fits a booking of 2
    EXPECT_EQ(b.count(), 4u);
}

TEST(ConcurrentBooker, BookFirstFreeWhenFromIsInsideABooking) {
    ConcurrentBooker b;
    b.book(10, 20);
    EXPECT_EQ(b.bookFirstFree(15, 5), 20);
}

TEST(ConcurrentBooker, BookFirstFreeRejectsNonPositiveDuration) {
    ConcurrentBooker b;
    EXPECT_THROW(b.bookFirstFree(0, 0), std::invalid_argument);
    EXPECT_THROW(b.bookFirstFree(0, -1), std::invalid_argument);
    EXPECT_EQ(b.count(), 0u);
}

TEST(ConcurrentBooker, BookAllBooksEveryInterval) {
    ConcurrentBooker b;
    EXPECT_TRUE(b.bookAll({{30, 40}, {0, 10}, {10, 20}}));
    EXPECT_EQ(b.bookings(), (std::vector<Interval>{{0, 10}, {10, 20}, {30, 40}}));
}

TEST(ConcurrentBooker, BookAllIsAllOrNothing) {
    ConcurrentBooker b;
    ASSERT_TRUE(b.book(10, 20));
    EXPECT_FALSE(b.bookAll({{0, 5}, {15, 25}, {30, 40}}));
    EXPECT_EQ(b.count(), 1u) << "a failed bookAll must book nothing";
    EXPECT_TRUE(b.book(0, 5));
    EXPECT_TRUE(b.book(30, 40));
}

TEST(ConcurrentBooker, BookAllRejectsOverlapInsideTheList) {
    ConcurrentBooker b;
    EXPECT_FALSE(b.bookAll({{0, 10}, {5, 15}}));
    EXPECT_EQ(b.count(), 0u);
}

TEST(ConcurrentBooker, BookAllOfEmptyListSucceeds) {
    ConcurrentBooker b;
    EXPECT_TRUE(b.bookAll({}));
    EXPECT_EQ(b.count(), 0u);
}

TEST(ConcurrentBooker, BookAllWithInvalidIntervalThrowsAndBooksNothing) {
    ConcurrentBooker b;
    EXPECT_THROW(b.bookAll({{0, 5}, {7, 7}}), std::invalid_argument);
    EXPECT_EQ(b.count(), 0u);
}

// ---- many threads ----

TEST(ConcurrentBooker, ConcurrentBooksOfOverlappingSlotsHaveOneWinner) {
    // Thread t books [t, t+16). Every pair of these overlaps, so in each
    // round exactly one thread may win.
    constexpr int kThreads = 16;
    for (int round = 0; round < 100; ++round) {
        ConcurrentBooker b;
        std::atomic<int> wins{0};
        runTogether(kThreads, [&](int t) {
            if (b.book(t, t + 16)) wins.fetch_add(1);
        });
        ASSERT_EQ(wins.load(), 1) << "round " << round;
        ASSERT_EQ(b.count(), 1u) << "round " << round;
    }
}

TEST(ConcurrentBooker, ConcurrentBookFirstFreeTilesTheCalendar) {
    // 8 threads each take the first free 10-unit slot 200 times. Nobody
    // cancels, so the slots must fill [0, 16000) with no gaps and no
    // overlaps: 0, 10, 20, ...
    constexpr int kThreads = 8;
    constexpr int kPerThread = 200;
    ConcurrentBooker b;
    std::vector<std::vector<long long>> starts(kThreads);
    runTogether(kThreads, [&](int t) {
        for (int i = 0; i < kPerThread; ++i) starts[t].push_back(b.bookFirstFree(0, 10));
    });

    std::vector<long long> all;
    for (const auto& s : starts) all.insert(all.end(), s.begin(), s.end());
    std::sort(all.begin(), all.end());
    ASSERT_EQ(all.size(), static_cast<std::size_t>(kThreads * kPerThread));
    for (std::size_t i = 0; i < all.size(); ++i) {
        ASSERT_EQ(all[i], static_cast<long long>(i) * 10)
            << "two threads got the same slot, or a slot was skipped";
    }
    auto booked = b.bookings();
    EXPECT_EQ(booked.size(), all.size());
    EXPECT_TRUE(noOverlaps(booked));
}

TEST(ConcurrentBooker, ConcurrentBookAllHasOneWinnerAndNoPartialBookings) {
    // Every thread's group contains the shared slot [50,60), plus two
    // slots of its own. Exactly one group can win each round -- and the
    // losers must leave none of their own slots behind.
    constexpr int kThreads = 8;
    for (int round = 0; round < 100; ++round) {
        ConcurrentBooker b;
        std::vector<char> won(kThreads, 0);
        runTogether(kThreads, [&](int t) {
            long long mine = 100 + 20 * t;
            won[t] = b.bookAll({{mine, mine + 10}, {50, 60}, {1000 + mine, 1010 + mine}});
        });

        int winners = static_cast<int>(std::count(won.begin(), won.end(), 1));
        ASSERT_EQ(winners, 1) << "round " << round;
        int w = static_cast<int>(std::find(won.begin(), won.end(), 1) - won.begin());
        long long mine = 100 + 20 * w;
        ASSERT_EQ(b.bookings(),
                  (std::vector<Interval>{{50, 60}, {mine, mine + 10}, {1000 + mine, 1010 + mine}}))
            << "round " << round << ": a losing bookAll left a partial booking";
    }
}

TEST(ConcurrentBooker, AFailedBookAllIsNeverSeenHalfDone) {
    // [50,60) is taken, so bookAll({[0,10), [50,60)}) must always fail and
    // book nothing -- which leaves [0,10) free for the one thread that
    // books and cancels it over and over. A bookAll that books [0,10)
    // first and takes it back once [50,60) turns out to be taken is not
    // all-or-nothing: for a moment, that thread finds [0,10) taken.
    constexpr int kGroupThreads = 3;
    constexpr int kTries = 20000;
    ConcurrentBooker b;
    ASSERT_TRUE(b.book(50, 60));
    std::atomic<bool> groupWon{false};
    std::atomic<int> blocked{0};
    runTogether(kGroupThreads + 1, [&](int t) {
        for (int i = 0; i < kTries; ++i) {
            if (t < kGroupThreads) {
                if (b.bookAll({{0, 10}, {50, 60}})) groupWon = true;
            } else if (b.book(0, 10)) {
                b.cancel(0, 10);
            } else {
                blocked.fetch_add(1);
            }
        }
    });

    EXPECT_FALSE(groupWon) << "bookAll succeeded although [50,60) was taken";
    EXPECT_EQ(blocked.load(), 0)
        << "book(0, 10) failed: a bookAll that failed held [0,10) for a moment";
    EXPECT_EQ(b.bookings(), (std::vector<Interval>{{50, 60}}));
}

TEST(ConcurrentBooker, ConcurrentMixedOperationsKeepTheCalendarConsistent) {
    // Threads book, cancel and bookFirstFree at random. Afterwards the
    // calendar must hold exactly the bookings that were made and not
    // cancelled -- and none may overlap.
    constexpr int kThreads = 8;
    constexpr int kOps = 3000;
    ConcurrentBooker b;
    std::vector<std::vector<Interval>> kept(kThreads);
    runTogether(kThreads, [&](int t) {
        std::mt19937 rng(1234 + t);
        auto& mine = kept[t];
        for (int i = 0; i < kOps; ++i) {
            int op = static_cast<int>(rng() % 3);
            if (op == 0) {
                long long s = static_cast<long long>(rng() % 5000);
                long long len = 1 + static_cast<long long>(rng() % 20);
                if (b.book(s, s + len)) mine.push_back({s, s + len});
            } else if (op == 1) {
                long long len = 1 + static_cast<long long>(rng() % 20);
                long long s = b.bookFirstFree(static_cast<long long>(rng() % 5000), len);
                mine.push_back({s, s + len});
            } else if (!mine.empty()) {
                std::size_t k = rng() % mine.size();
                ASSERT_TRUE(b.cancel(mine[k].start, mine[k].end))
                    << "could not cancel a booking this thread made";
                mine.erase(mine.begin() + static_cast<std::ptrdiff_t>(k));
            }
        }
    });

    std::vector<Interval> expected;
    for (const auto& m : kept) expected.insert(expected.end(), m.begin(), m.end());
    std::sort(expected.begin(), expected.end(),
              [](const Interval& a, const Interval& c) { return a.start < c.start; });
    auto booked = b.bookings();
    EXPECT_TRUE(sortedByStart(booked));
    EXPECT_TRUE(noOverlaps(booked));
    EXPECT_EQ(booked, expected);
    EXPECT_EQ(b.count(), expected.size());
}
