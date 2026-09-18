#include "day02/bounded_stack.h"
#include <gtest/gtest.h>

TEST(BoundedStack, StartsEmpty) {
    BoundedStack s(2);
    EXPECT_TRUE(s.empty());
    EXPECT_EQ(s.size(), 0u);
    EXPECT_FALSE(s.full());
}

TEST(BoundedStack, PushUpToCapacity) {
    BoundedStack s(2);
    s.push(5);
    s.push(7);
    EXPECT_EQ(s.size(), 2u);
    EXPECT_TRUE(s.full());
}

TEST(BoundedStack, PushBeyondCapacityThrows) {
    BoundedStack s(1);
    s.push(1);
    EXPECT_THROW(s.push(2), std::runtime_error);
    EXPECT_EQ(s.size(), 1u);  // rejected push must not touch state
}

TEST(BoundedStack, PopReturnsLifoOrder) {
    BoundedStack s(3);
    s.push(1);
    s.push(2);
    s.push(3);
    EXPECT_EQ(s.pop(), 3);
    EXPECT_EQ(s.pop(), 2);
    EXPECT_EQ(s.size(), 1u);
}

TEST(BoundedStack, PopEmptyThrows) {
    BoundedStack s(1);
    EXPECT_THROW(s.pop(), std::runtime_error);
}

TEST(BoundedStack, TopDoesNotRemove) {
    BoundedStack s(2);
    s.push(9);
    EXPECT_EQ(s.top(), 9);
    EXPECT_EQ(s.size(), 1u);
}

TEST(BoundedStack, TopEmptyThrows) {
    BoundedStack s(1);
    EXPECT_THROW(s.top(), std::runtime_error);
}

TEST(BoundedStack, ZeroCapacityIsAlwaysFull) {
    BoundedStack s(0);
    EXPECT_TRUE(s.full());
    EXPECT_THROW(s.push(1), std::runtime_error);
}

TEST(BoundedStack, PopFreesCapacityForAnotherPush) {
    BoundedStack s(1);
    s.push(1);
    ASSERT_TRUE(s.full());
    s.pop();
    EXPECT_FALSE(s.full());
    EXPECT_TRUE(s.empty());
    s.push(2);  // must succeed -- capacity was freed, not stuck full forever
    EXPECT_EQ(s.top(), 2);
    EXPECT_EQ(s.size(), 1u);
}

TEST(BoundedStack, TopReflectsMostRecentPush) {
    BoundedStack s(3);
    s.push(1);
    EXPECT_EQ(s.top(), 1);
    s.push(2);
    EXPECT_EQ(s.top(), 2);
    s.push(3);
    EXPECT_EQ(s.top(), 3);
}

TEST(BoundedStack, EmptyAfterPoppingEverything) {
    BoundedStack s(2);
    s.push(1);
    s.push(2);
    s.pop();
    s.pop();
    EXPECT_TRUE(s.empty());
    EXPECT_EQ(s.size(), 0u);
    EXPECT_THROW(s.pop(), std::runtime_error);
}
