#include "day02/txn_kv.h"
#include <gtest/gtest.h>

TEST(TxnKV, SetThenGet) {
    TxnKV kv;
    kv.set("a", "1");
    ASSERT_TRUE(kv.get("a").has_value());
    EXPECT_EQ(*kv.get("a"), "1");
}

TEST(TxnKV, GetMissingKeyReturnsNothing) {
    TxnKV kv;
    EXPECT_FALSE(kv.get("nope").has_value());
}

TEST(TxnKV, OverwriteReplacesValue) {
    TxnKV kv;
    kv.set("a", "1");
    kv.set("a", "2");
    EXPECT_EQ(*kv.get("a"), "2");
}

TEST(TxnKV, EraseReturnsTrueOnlyIfKeyWasVisible) {
    TxnKV kv;
    kv.set("a", "1");
    EXPECT_TRUE(kv.erase("a"));
    EXPECT_FALSE(kv.get("a").has_value());
    EXPECT_FALSE(kv.erase("a"));
    EXPECT_FALSE(kv.erase("never-set"));
}

TEST(TxnKV, CommitWithoutOpenTransactionFails) {
    TxnKV kv;
    EXPECT_FALSE(kv.commit());
}

TEST(TxnKV, RollbackWithoutOpenTransactionFails) {
    TxnKV kv;
    EXPECT_FALSE(kv.rollback());
}

TEST(TxnKV, SetInsideTransactionIsVisibleBeforeCommit) {
    TxnKV kv;
    kv.begin();
    kv.set("a", "1");
    ASSERT_TRUE(kv.get("a").has_value());
    EXPECT_EQ(*kv.get("a"), "1");
}

TEST(TxnKV, RollbackDiscardsEveryChangeOfTheTransaction) {
    TxnKV kv;
    kv.set("a", "old");
    kv.begin();
    kv.set("a", "new");
    kv.set("b", "1");
    EXPECT_TRUE(kv.rollback());
    EXPECT_EQ(*kv.get("a"), "old");
    EXPECT_FALSE(kv.get("b").has_value());
}

TEST(TxnKV, CommitMakesChangesPermanent) {
    TxnKV kv;
    kv.begin();
    kv.set("a", "1");
    EXPECT_TRUE(kv.commit());
    ASSERT_TRUE(kv.get("a").has_value());
    EXPECT_EQ(*kv.get("a"), "1");
    EXPECT_FALSE(kv.rollback());  // nothing left open to roll back
    EXPECT_EQ(*kv.get("a"), "1");
}

TEST(TxnKV, EraseInsideTransactionHidesOuterValue) {
    TxnKV kv;
    kv.set("a", "base");
    kv.begin();
    EXPECT_TRUE(kv.erase("a"));
    EXPECT_FALSE(kv.get("a").has_value());  // gone as seen from inside
}

TEST(TxnKV, RollbackRestoresValueErasedInsideTransaction) {
    TxnKV kv;
    kv.set("a", "base");
    kv.begin();
    kv.erase("a");
    kv.rollback();
    ASSERT_TRUE(kv.get("a").has_value());
    EXPECT_EQ(*kv.get("a"), "base");
}

TEST(TxnKV, CommittedEraseDeletesFromTheStore) {
    TxnKV kv;
    kv.set("a", "base");
    kv.begin();
    kv.erase("a");
    EXPECT_TRUE(kv.commit());
    EXPECT_FALSE(kv.get("a").has_value());
}

TEST(TxnKV, NestedInnerRollbackKeepsOuterChanges) {
    TxnKV kv;
    kv.begin();
    kv.set("outer", "1");
    kv.begin();
    kv.set("inner", "2");
    EXPECT_TRUE(kv.rollback());   // discards only the inner transaction
    EXPECT_FALSE(kv.get("inner").has_value());
    ASSERT_TRUE(kv.get("outer").has_value());
    EXPECT_EQ(*kv.get("outer"), "1");
}

TEST(TxnKV, InnerCommitMergesIntoParentNotIntoStore) {
    TxnKV kv;
    kv.begin();
    kv.begin();
    kv.set("a", "1");
    EXPECT_TRUE(kv.commit());     // inner commit: `a` now belongs to the OUTER txn
    EXPECT_EQ(*kv.get("a"), "1");
    EXPECT_TRUE(kv.rollback());   // outer rollback must take `a` down with it
    EXPECT_FALSE(kv.get("a").has_value());
}

TEST(TxnKV, InnerCommittedEraseStaysInsideParent) {
    TxnKV kv;
    kv.set("a", "base");
    kv.begin();
    kv.begin();
    kv.erase("a");
    EXPECT_TRUE(kv.commit());     // erase merged into the outer txn
    EXPECT_FALSE(kv.get("a").has_value());
    EXPECT_TRUE(kv.rollback());   // outer rollback cancels the erase too
    ASSERT_TRUE(kv.get("a").has_value());
    EXPECT_EQ(*kv.get("a"), "base");
}

TEST(TxnKV, SetAfterEraseInSameTransactionResurrectsKey) {
    TxnKV kv;
    kv.set("a", "base");
    kv.begin();
    kv.erase("a");
    kv.set("a", "reborn");
    EXPECT_EQ(*kv.get("a"), "reborn");
    kv.commit();
    EXPECT_EQ(*kv.get("a"), "reborn");
}

TEST(TxnKV, ValueIsVisibleThroughMultipleUntouchedLayers) {
    TxnKV kv;
    kv.set("a", "base");
    kv.begin();
    kv.begin();
    kv.begin();
    ASSERT_TRUE(kv.get("a").has_value());
    EXPECT_EQ(*kv.get("a"), "base");
}
