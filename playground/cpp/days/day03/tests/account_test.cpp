#include "day03/account.h"
#include "run_together.h"
#include <gtest/gtest.h>

#include <atomic>
#include <memory>
#include <random>
#include <vector>

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

// ---- transfer and total: two accounts locked at once ----

TEST(Account, TransferMovesMoney) {
    Account a(100), b(50);
    EXPECT_TRUE(a.transfer(b, 30));
    EXPECT_EQ(a.balance(), 70);
    EXPECT_EQ(b.balance(), 80);
}

TEST(Account, TransferWithoutEnoughMoneyChangesNothing) {
    Account a(10), b(0);
    EXPECT_FALSE(a.transfer(b, 11));
    EXPECT_EQ(a.balance(), 10);
    EXPECT_EQ(b.balance(), 0);
}

TEST(Account, TransferRejectsBadArguments) {
    Account a(10), b(0);
    EXPECT_THROW(a.transfer(b, 0), std::invalid_argument);
    EXPECT_THROW(a.transfer(b, -5), std::invalid_argument);
    EXPECT_THROW(a.transfer(a, 1), std::invalid_argument) << "a transfer to the same account";
    EXPECT_EQ(a.balance(), 10);
    EXPECT_EQ(b.balance(), 0);
}

TEST(Account, TotalAddsBothBalances) {
    Account a(100), b(50);
    EXPECT_EQ(Account::total(a, b), 150);
    EXPECT_EQ(Account::total(b, a), 150);
    EXPECT_THROW(Account::total(a, a), std::invalid_argument);
}

TEST(Account, OppositeTransfersDoNotDeadlock) {
    // Half the threads move money from x to y, the other half from y to
    // x -- the slides' "Two transfers, two locks", 80,000 times.
    constexpr int kThreads = 4;
    constexpr int kPerThread = 20000;
    Account x(1000), y(1000);
    runTogetherWithin(20, kThreads, [&](int t) {
        for (int i = 0; i < kPerThread; ++i) {
            if (t % 2 == 0) {
                x.transfer(y, 1);
            } else {
                y.transfer(x, 1);
            }
        }
    });

    EXPECT_EQ(x.balance() + y.balance(), 2000) << "money appeared or disappeared";
    EXPECT_GE(x.balance(), 0);
    EXPECT_GE(y.balance(), 0);
}

TEST(Account, TotalNeverSeesAHalfDoneTransfer) {
    // Transfers move money back and forth between x and y while a watcher
    // keeps asking for the total. A transfer done as two steps (take
    // from x, then give to y) shows the watcher money in flight.
    constexpr int kMovers = 4;
    constexpr int kPerThread = 20000;
    Account x(1000), y(1000);
    std::atomic<int> moversLeft{kMovers};
    std::atomic<long long> wrongTotal{2000};
    runTogetherWithin(20, kMovers + 1, [&](int t) {
        if (t == kMovers) {  // the watcher
            while (moversLeft.load() > 0) {
                long long seen = Account::total(x, y);
                if (seen != 2000) wrongTotal = seen;
            }
            return;
        }
        for (int i = 0; i < kPerThread; ++i) {
            if (t % 2 == 0) {
                x.transfer(y, 7);
            } else {
                y.transfer(x, 7);
            }
        }
        moversLeft.fetch_sub(1);
    });

    EXPECT_EQ(wrongTotal.load(), 2000) << "total() saw money that was in flight";
    EXPECT_EQ(Account::total(x, y), 2000);
}

TEST(Account, ConcurrentTransfersAmongManyAccountsKeepTheMoney) {
    // 8 threads move random amounts between random pairs of 6 accounts,
    // in both directions. No deadlock, no money made or lost, no account
    // below zero.
    constexpr int kAccounts = 6;
    constexpr int kThreads = 8;
    constexpr int kPerThread = 20000;
    std::vector<std::unique_ptr<Account>> accounts;
    for (int i = 0; i < kAccounts; ++i) accounts.push_back(std::make_unique<Account>(1000));
    runTogetherWithin(20, kThreads, [&](int t) {
        std::mt19937 rng(99 + t);
        for (int i = 0; i < kPerThread; ++i) {
            int from = static_cast<int>(rng() % kAccounts);
            int to = static_cast<int>(rng() % (kAccounts - 1));
            if (to >= from) ++to;  // any account but `from`
            accounts[from]->transfer(*accounts[to], 1 + static_cast<long long>(rng() % 50));
        }
    });

    long long sum = 0;
    for (const auto& a : accounts) {
        EXPECT_GE(a->balance(), 0);
        sum += a->balance();
    }
    EXPECT_EQ(sum, 1000LL * kAccounts) << "money appeared or disappeared";
}
