#include "day02/idempotent_payments.h"
#include <gtest/gtest.h>

TEST(IdempotentPayments, FirstChargeIsNew) {
    PaymentProcessor p;
    auto r = p.charge("orderA", 100);
    EXPECT_TRUE(r.wasNew);
    EXPECT_EQ(r.amount, 100);
    EXPECT_EQ(r.runningTotal, 100);
    EXPECT_EQ(p.totalCharged(), 100);
}

TEST(IdempotentPayments, RepeatKeyIsNotCharged) {
    PaymentProcessor p;
    p.charge("orderA", 100);
    auto r = p.charge("orderA", 100);
    EXPECT_FALSE(r.wasNew);
    EXPECT_EQ(r.amount, 100);
    EXPECT_EQ(r.runningTotal, 100);   // unchanged -- not double-charged
    EXPECT_EQ(p.totalCharged(), 100);
}

TEST(IdempotentPayments, DifferentKeysAccumulate) {
    PaymentProcessor p;
    p.charge("orderA", 100);
    p.charge("orderB", 250);
    EXPECT_EQ(p.totalCharged(), 350);
}

TEST(IdempotentPayments, RepeatDoesNotCareAboutNewAmountArgument) {
    // A caller retrying a request might resend the same amount; a repeat
    // must still report the ORIGINAL amount, never a re-derived one.
    PaymentProcessor p;
    p.charge("x", 5);
    auto r = p.charge("x", 5);
    EXPECT_EQ(r.amount, 5);
    EXPECT_EQ(p.totalCharged(), 5);
}

TEST(IdempotentPayments, ManyRepeatsStillChargeOnce) {
    PaymentProcessor p;
    p.charge("x", 5);
    p.charge("x", 5);
    p.charge("x", 5);
    p.charge("y", 10);
    EXPECT_EQ(p.totalCharged(), 15);
}

TEST(IdempotentPayments, RepeatIsIdempotentEvenWithOtherChargesBetween) {
    PaymentProcessor p;
    p.charge("A", 100);
    p.charge("B", 50);
    auto r = p.charge("A", 100);   // repeat of A, not adjacent to the first
    EXPECT_FALSE(r.wasNew);
    EXPECT_EQ(r.amount, 100);
    EXPECT_EQ(r.runningTotal, 150);  // reflects total so far, not re-added
    p.charge("C", 25);
    EXPECT_EQ(p.totalCharged(), 175);  // 100 + 50 + 25, A never double-counted
}

TEST(IdempotentPayments, ZeroAmountChargeIsStillTrackedByKey) {
    PaymentProcessor p;
    auto first = p.charge("free-trial", 0);
    EXPECT_TRUE(first.wasNew);
    EXPECT_EQ(first.amount, 0);
    auto repeat = p.charge("free-trial", 0);
    EXPECT_FALSE(repeat.wasNew);
    EXPECT_EQ(p.totalCharged(), 0);
}
