#include "day03/account.h"
#include "run_together.h"
#include <gtest/gtest.h>

#include <atomic>

TEST(Account, StartsWithInitialBalance) {
    Account a(100);
    EXPECT_EQ(a.balance(), 100);
}

TEST(Account, NegativeInitialBalanceThrows) {
    EXPECT_THROW(Account(-1), std::invalid_argument);
}

TEST(Account, DepositAddsAndWithdrawSubtracts) {
    Account a(0);
    a.deposit(50);
    EXPECT_TRUE(a.withdraw(20));
    EXPECT_EQ(a.balance(), 30);
}

TEST(Account, WithdrawExactBalanceSucceeds) {
    Account a(40);
    EXPECT_TRUE(a.withdraw(40));
    EXPECT_EQ(a.balance(), 0);
}

TEST(Account, WithdrawMoreThanBalanceFailsAndChangesNothing) {
    Account a(30);
    EXPECT_FALSE(a.withdraw(31));
    EXPECT_EQ(a.balance(), 30);
}

TEST(Account, NonPositiveAmountsThrowAndChangeNothing) {
    Account a(10);
    EXPECT_THROW(a.deposit(0), std::invalid_argument);
    EXPECT_THROW(a.deposit(-5), std::invalid_argument);
    EXPECT_THROW(a.withdraw(0), std::invalid_argument);
    EXPECT_THROW(a.withdraw(-5), std::invalid_argument);
    EXPECT_EQ(a.balance(), 10);
}

TEST(Account, ConcurrentWithdrawalsNeverOverdraw) {
    // 8 threads try to take 40,000 in total from an account holding
    // 10,000. Exactly 10,000 withdrawals of 1 can succeed. The race only
    // matters at the moment the money runs out, so play it 20 times.
    constexpr int kThreads = 8;
    constexpr int kTriesPerThread = 5000;
    for (int round = 0; round < 20; ++round) {
        Account a(10000);
        std::atomic<long long> succeeded{0};
        runTogether(kThreads, [&](int) {
            for (int i = 0; i < kTriesPerThread; ++i) {
                if (a.withdraw(1)) succeeded.fetch_add(1);
            }
        });

        ASSERT_EQ(succeeded.load(), 10000) << "round " << round;
        ASSERT_EQ(a.balance(), 0) << "round " << round;
    }
}

TEST(Account, ConcurrentDepositsAreNeverLost) {
    constexpr int kThreads = 8;
    constexpr int kPerThread = 20000;
    Account a(0);
    runTogether(kThreads, [&](int) {
        for (int i = 0; i < kPerThread; ++i) a.deposit(1);
    });

    EXPECT_EQ(a.balance(), kThreads * kPerThread);
}

TEST(Account, ConcurrentDepositsAndWithdrawalsBalanceOut) {
    // 4 threads deposit 3 at a time, 4 threads withdraw 2 at a time, and
    // one watcher checks that the balance is never seen below zero.
    constexpr int kPerThread = 20000;
    Account a(0);
    std::atomic<long long> withdrawn{0};
    std::atomic<bool> sawNegative{false};
    std::atomic<int> workersLeft{8};
    runTogether(9, [&](int t) {
        if (t == 8) {  // the watcher
            while (workersLeft.load() > 0) {
                if (a.balance() < 0) sawNegative = true;
            }
            return;
        }
        for (int i = 0; i < kPerThread; ++i) {
            if (t < 4) {
                a.deposit(3);
            } else if (a.withdraw(2)) {
                withdrawn.fetch_add(2);
            }
        }
        workersLeft.fetch_sub(1);
    });

    EXPECT_FALSE(sawNegative);
    EXPECT_EQ(a.balance(), 4LL * kPerThread * 3 - withdrawn.load());
}
