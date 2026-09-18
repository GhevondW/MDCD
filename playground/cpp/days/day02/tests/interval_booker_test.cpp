#include "day02/interval_booker.h"
#include <gtest/gtest.h>

TEST(IntervalBooker, BookOnEmptyCalendarSucceeds) {
    IntervalBooker b;
    EXPECT_TRUE(b.book(10, 20));
    EXPECT_EQ(b.count(), 1u);
}

TEST(IntervalBooker, AdjacentIntervalsDoNotOverlap) {
    // Half-open semantics: [0,10) and [10,20) share only the boundary point.
    IntervalBooker b;
    EXPECT_TRUE(b.book(0, 10));
    EXPECT_TRUE(b.book(10, 20));
    EXPECT_EQ(b.count(), 2u);
}

TEST(IntervalBooker, PartialOverlapIsRejectedBothDirections) {
    IntervalBooker b;
    ASSERT_TRUE(b.book(10, 20));
    EXPECT_FALSE(b.book(15, 25));   // overlaps the tail
    EXPECT_FALSE(b.book(5, 15));    // overlaps the head
    EXPECT_EQ(b.count(), 1u);
}

TEST(IntervalBooker, ContainmentIsRejectedBothDirections) {
    IntervalBooker b;
    ASSERT_TRUE(b.book(10, 20));
    EXPECT_FALSE(b.book(12, 18));   // inside the existing booking
    EXPECT_FALSE(b.book(5, 25));    // swallows the existing booking
}

TEST(IntervalBooker, IdenticalIntervalIsRejected) {
    IntervalBooker b;
    ASSERT_TRUE(b.book(10, 20));
    EXPECT_FALSE(b.book(10, 20));
}

TEST(IntervalBooker, RejectedBookChangesNothing) {
    IntervalBooker b;
    ASSERT_TRUE(b.book(10, 20));
    EXPECT_FALSE(b.book(15, 25));
    EXPECT_EQ(b.count(), 1u);
    EXPECT_TRUE(b.book(20, 30));    // the slot the rejected book didn't take
}

TEST(IntervalBooker, EmptyOrReversedIntervalThrows) {
    IntervalBooker b;
    EXPECT_THROW(b.book(10, 10), std::invalid_argument);
    EXPECT_THROW(b.book(20, 10), std::invalid_argument);
    EXPECT_EQ(b.count(), 0u);
}

TEST(IntervalBooker, CancelExactMatchFreesTheSlot) {
    IntervalBooker b;
    ASSERT_TRUE(b.book(10, 20));
    EXPECT_TRUE(b.cancel(10, 20));
    EXPECT_EQ(b.count(), 0u);
    EXPECT_TRUE(b.book(15, 25));    // the freed range is bookable again
}

TEST(IntervalBooker, CancelRequiresExactMatch) {
    IntervalBooker b;
    ASSERT_TRUE(b.book(10, 20));
    EXPECT_FALSE(b.cancel(10, 15));   // sub-range: no
    EXPECT_FALSE(b.cancel(5, 25));    // super-range: no
    EXPECT_EQ(b.count(), 1u);
}

TEST(IntervalBooker, CancelUnknownIntervalReturnsFalse) {
    IntervalBooker b;
    EXPECT_FALSE(b.cancel(10, 20));
}

TEST(IntervalBooker, FirstFreeOnEmptyCalendarIsFrom) {
    IntervalBooker b;
    EXPECT_EQ(b.firstFree(7, 5), 7);
}

TEST(IntervalBooker, FirstFreeSkipsBusyRangesAndTooSmallGaps) {
    IntervalBooker b;
    b.book(10, 20);
    b.book(20, 30);
    // [5,15) collides; the only room left starts after the last booking.
    EXPECT_EQ(b.firstFree(5, 10), 30);
}

TEST(IntervalBooker, FirstFreeFindsGapOfExactSize) {
    IntervalBooker b;
    b.book(0, 10);
    b.book(15, 25);
    EXPECT_EQ(b.firstFree(0, 5), 10);   // the gap [10,15) fits exactly
}

TEST(IntervalBooker, FirstFreeWhenFromFallsInsideABooking) {
    IntervalBooker b;
    b.book(10, 20);
    EXPECT_EQ(b.firstFree(12, 5), 20);
}

TEST(IntervalBooker, FirstFreeIgnoresBookingsEntirelyBeforeFrom) {
    IntervalBooker b;
    b.book(0, 10);
    EXPECT_EQ(b.firstFree(50, 5), 50);
}

TEST(IntervalBooker, FirstFreeResultIsActuallyBookable) {
    IntervalBooker b;
    b.book(3, 9);
    b.book(12, 40);
    b.book(41, 50);
    long long s = b.firstFree(0, 4);
    EXPECT_TRUE(b.book(s, s + 4));
}

TEST(IntervalBooker, NonPositiveDurationThrows) {
    IntervalBooker b;
    EXPECT_THROW(b.firstFree(0, 0), std::invalid_argument);
    EXPECT_THROW(b.firstFree(0, -3), std::invalid_argument);
}
