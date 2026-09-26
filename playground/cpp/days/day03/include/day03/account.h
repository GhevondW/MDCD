#pragma once
// Day 3 -- Critical sections: check-then-act
//
// A bank account shared by many threads -- think of one card used by
// several people at the same moment. The balance must never go below
// zero.
//
//   Account(initial)   start with `initial`. initial must be >= 0; if
//                      not, throw std::invalid_argument.
//   deposit(amount)    add `amount`. amount must be > 0; if not, throw
//                      std::invalid_argument and change nothing.
//   withdraw(amount)   amount must be > 0; if not, throw
//                      std::invalid_argument and change nothing.
//                      If balance >= amount: subtract it, return true.
//                      Otherwise return false and change nothing.
//   balance()          the current balance.
//
// All four may be called by many threads at the same time.
//
//     if (balance >= amount) balance -= amount;
//
// is a check and an act -- another thread can slip in between them.

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

private:
    // TODO: choose your own representation.
};
