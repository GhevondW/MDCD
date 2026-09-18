#include "day02/expiring_payments.h"
#include <gtest/gtest.h>

using Status = ExpiringPaymentProcessor::Status;

TEST(ExpiringPayments, FirstChargeIsNew) {
    ExpiringPaymentProcessor p(10);
    auto r = p.charge("k", 100, 0);
    EXPECT_EQ(r.status, Status::New);
    EXPECT_EQ(r.amount, 100);
    EXPECT_EQ(r.runningTotal, 100);
    EXPECT_EQ(p.totalCharged(), 100);
}

TEST(ExpiringPayments, SameAmountWithinTtlIsReplay) {
    ExpiringPaymentProcessor p(10);
    p.charge("k", 100, 0);
    auto r = p.charge("k", 100, 5);
    EXPECT_EQ(r.status, Status::Replay);
    EXPECT_EQ(r.amount, 100);
    EXPECT_EQ(r.runningTotal, 100);
    EXPECT_EQ(p.totalCharged(), 100);
}

TEST(ExpiringPayments, DifferentAmountWithinTtlIsConflict) {
    ExpiringPaymentProcessor p(10);
    p.charge("k", 100, 0);
    auto r = p.charge("k", 250, 5);
    EXPECT_EQ(r.status, Status::Conflict);
    EXPECT_EQ(r.amount, 100);            // the RECORDED amount, not 250
    EXPECT_EQ(r.runningTotal, 100);
    EXPECT_EQ(p.totalCharged(), 100);    // nothing was charged
}

TEST(ExpiringPayments, ExpiredKeyIsChargedAgain) {
    ExpiringPaymentProcessor p(10);
    p.charge("k", 100, 0);
    auto r = p.charge("k", 100, 20);
    EXPECT_EQ(r.status, Status::New);
    EXPECT_EQ(p.totalCharged(), 200);
}

TEST(ExpiringPayments, ExpiryBoundaryIsInclusive) {
    // Recorded at 0 with ttl 10 -> expired exactly at now == 10.
    ExpiringPaymentProcessor p(10);
    p.charge("k", 100, 0);
    EXPECT_EQ(p.charge("k", 100, 9).status, Status::Replay);
    EXPECT_EQ(p.charge("k", 100, 10).status, Status::New);
    EXPECT_EQ(p.totalCharged(), 200);
}

TEST(ExpiringPayments, ReplayDoesNotExtendExpiry) {
    ExpiringPaymentProcessor p(10);
    p.charge("k", 100, 0);
    EXPECT_EQ(p.charge("k", 100, 9).status, Status::Replay);
    // If the replay at 9 had refreshed the record, the key would still be
    // active at 10. It must not be.
    EXPECT_EQ(p.charge("k", 100, 10).status, Status::New);
}

TEST(ExpiringPayments, ConflictChangesNothing) {
    ExpiringPaymentProcessor p(10);
    p.charge("k", 100, 0);
    p.charge("k", 250, 5);                          // conflict
    auto r = p.charge("k", 100, 9);                 // original amount again
    EXPECT_EQ(r.status, Status::Replay);            // record still 100
    EXPECT_EQ(r.amount, 100);
    EXPECT_EQ(p.totalCharged(), 100);
    EXPECT_EQ(p.charge("k", 100, 10).status, Status::New);  // expiry not shifted
}

TEST(ExpiringPayments, AfterExpiryDifferentAmountIsNewNotConflict) {
    ExpiringPaymentProcessor p(10);
    p.charge("k", 100, 0);
    auto r = p.charge("k", 250, 15);
    EXPECT_EQ(r.status, Status::New);
    EXPECT_EQ(r.amount, 250);
    EXPECT_EQ(p.totalCharged(), 350);
    // The record was replaced: 250 is now the recorded amount.
    EXPECT_EQ(p.charge("k", 250, 16).status, Status::Replay);
    EXPECT_EQ(p.charge("k", 100, 16).status, Status::Conflict);
}

TEST(ExpiringPayments, KeysAreIndependent) {
    ExpiringPaymentProcessor p(10);
    p.charge("a", 100, 0);
    p.charge("b", 50, 8);
    EXPECT_EQ(p.totalCharged(), 150);
    EXPECT_EQ(p.charge("a", 100, 12).status, Status::New);     // a expired
    EXPECT_EQ(p.charge("b", 50, 12).status, Status::Replay);   // b still active
}

TEST(ExpiringPayments, ActiveKeysCountsOnlyUnexpired) {
    ExpiringPaymentProcessor p(10);
    p.charge("a", 1, 0);   // expires at 10
    p.charge("b", 2, 5);   // expires at 15
    p.charge("c", 3, 8);   // expires at 18
    EXPECT_EQ(p.activeKeys(0), 3u);
    EXPECT_EQ(p.activeKeys(12), 2u);
    EXPECT_EQ(p.activeKeys(15), 1u);
    EXPECT_EQ(p.activeKeys(18), 0u);
}

TEST(ExpiringPayments, ZeroTtlMeansEveryChargeIsNew) {
    ExpiringPaymentProcessor p(0);
    EXPECT_EQ(p.charge("k", 100, 5).status, Status::New);
    EXPECT_EQ(p.charge("k", 100, 5).status, Status::New);
    EXPECT_EQ(p.totalCharged(), 200);
    EXPECT_EQ(p.activeKeys(5), 0u);
}

TEST(ExpiringPayments, ReplayReportsCurrentTotalNotTotalAtFirstCharge) {
    ExpiringPaymentProcessor p(100);
    p.charge("a", 100, 0);
    p.charge("b", 50, 1);
    auto r = p.charge("a", 100, 2);
    EXPECT_EQ(r.status, Status::Replay);
    EXPECT_EQ(r.runningTotal, 150);
}

TEST(ExpiringPayments, TotalIsSumOfNewChargesOnly) {
    ExpiringPaymentProcessor p(10);
    p.charge("a", 100, 0);   // new
    p.charge("a", 100, 1);   // replay
    p.charge("a", 999, 2);   // conflict
    p.charge("b", 50, 3);    // new
    p.charge("a", 25, 12);   // expired -> new
    EXPECT_EQ(p.totalCharged(), 175);
}
