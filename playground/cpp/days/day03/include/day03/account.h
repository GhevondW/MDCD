#pragma once
// Day 3 -- Critical sections: check-then-act, and two locks at once
//
// A bank account shared by many threads -- think of one card used by
// several people at the same moment. The balance must never go below
// zero.
//
//   Account(initial)      start with `initial`. initial must be >= 0; if
//                         not, throw std::invalid_argument.
//   deposit(amount)       add `amount`. amount must be > 0; if not, throw
//                         std::invalid_argument and change nothing.
//   withdraw(amount)      amount must be > 0; if not, throw
//                         std::invalid_argument and change nothing.
//                         If balance >= amount: subtract it, return true.
//                         Otherwise return false and change nothing.
//   balance()             the current balance.
//   transfer(to, amount)  move `amount` from this account to `to`, as ONE
//                         step: no thread may ever see the money gone
//                         from one account but not yet in the other.
//                         amount must be > 0 and `to` must be a different
//                         account; if not, throw std::invalid_argument.
//                         If this balance < amount, return false and
//                         change nothing.
//   total(a, b)           a's balance plus b's balance, read as ONE step:
//                         every transfer between a and b happens either
//                         completely before it or completely after it.
//                         a and b must be different accounts; if not,
//                         throw std::invalid_argument.
//
// Every method may be called by many threads at the same time.
//
//     if (balance >= amount) balance -= amount;
//
// is a check and an act -- another thread can slip in between them.
//
// transfer and total need two accounts locked at once. Two threads that
// lock the same two accounts in opposite orders wait for each other
// forever -- the slides' "Two transfers, two locks". Hint: if every
// thread locks any two accounts in the same order, that cannot happen.
// And a method that already holds a lock must not call another method
// that takes the same lock: that is locking twice.

#include <stdexcept>

class Account {
public:
    explicit Account(long long initial) {
        (void)initial;
        // TODO
    }

    void deposit(long long amount) {
        (void)amount;
        // TODO
    }

    bool withdraw(long long amount) {
        (void)amount;
        // TODO
        return false;
    }

    long long balance() const {
        // TODO
        return 0;
    }

    bool transfer(Account& to, long long amount) {
        (void)to;
        (void)amount;
        // TODO
        return false;
    }

    static long long total(const Account& a, const Account& b) {
        (void)a;
        (void)b;
        // TODO
        return 0;
    }

private:
    // TODO: choose your own representation.
};
