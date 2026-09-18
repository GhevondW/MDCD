#include "day02/atomic_record.h"
#include <gtest/gtest.h>

TEST(AtomicRecord, GetOnUnknownKeyIsNullopt) {
    RecordStore store;
    EXPECT_FALSE(store.get("ghost").has_value());
}

TEST(AtomicRecord, ValidSetCommitsBothFields) {
    RecordStore store;
    EXPECT_TRUE(store.set("acct1", 10, 5));
    auto r = store.get("acct1");
    ASSERT_TRUE(r.has_value());
    EXPECT_EQ(r->a, 10);
    EXPECT_EQ(r->b, 5);
}

TEST(AtomicRecord, InvalidSetIsRejected) {
    RecordStore store;
    EXPECT_FALSE(store.set("acct2", -5, -5));
    EXPECT_FALSE(store.get("acct2").has_value());
}

TEST(AtomicRecord, RejectedSetLeavesExistingRecordUntouched) {
    // The core atomicity check: a rejected update on an EXISTING key must
    // not partially apply -- field a must not change even though it would
    // have been "valid" applied in isolation.
    RecordStore store;
    store.set("acct1", 10, 5);
    EXPECT_FALSE(store.set("acct1", -20, 3));  // -20 + 3 < 0 -> rejected
    auto r = store.get("acct1");
    ASSERT_TRUE(r.has_value());
    EXPECT_EQ(r->a, 10);   // unchanged, not -20
    EXPECT_EQ(r->b, 5);    // unchanged
}

TEST(AtomicRecord, ZeroSumIsValid) {
    RecordStore store;
    EXPECT_TRUE(store.set("a", 0, 0));
    auto r = store.get("a");
    ASSERT_TRUE(r.has_value());
    EXPECT_EQ(r->a, 0);
    EXPECT_EQ(r->b, 0);
}

TEST(AtomicRecord, NegativeFieldWithCompensatingPositiveIsValid) {
    // The rule is a + b >= 0, NOT "both fields individually non-negative" --
    // a single negative field with a large enough positive one is fine.
    RecordStore store;
    EXPECT_TRUE(store.set("acct3", -5, 10));
    auto r = store.get("acct3");
    ASSERT_TRUE(r.has_value());
    EXPECT_EQ(r->a, -5);
    EXPECT_EQ(r->b, 10);
}

TEST(AtomicRecord, ValidSetOverwritesPreviousValidRecord) {
    RecordStore store;
    store.set("acct1", 10, 5);
    EXPECT_TRUE(store.set("acct1", 1, 2));  // second valid set replaces it
    auto r = store.get("acct1");
    ASSERT_TRUE(r.has_value());
    EXPECT_EQ(r->a, 1);
    EXPECT_EQ(r->b, 2);
}

TEST(AtomicRecord, KeysAreIndependent) {
    RecordStore store;
    store.set("acct1", 10, 5);
    store.set("acct2", 1, 1);
    EXPECT_FALSE(store.set("acct2", -100, 0));  // reject on acct2 only
    auto r1 = store.get("acct1");
    auto r2 = store.get("acct2");
    ASSERT_TRUE(r1.has_value());
    ASSERT_TRUE(r2.has_value());
    EXPECT_EQ(r1->a, 10);
    EXPECT_EQ(r1->b, 5);
    EXPECT_EQ(r2->a, 1);   // acct1's rejection-adjacent op didn't touch acct2
    EXPECT_EQ(r2->b, 1);
}
