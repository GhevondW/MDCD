#pragma once
// Day 2 -- API Guarantees: Idempotency
//
// Every charge request comes with a key. The same key must never be
// charged twice.
//
// charge(key, amount):
//   - If this key was never seen before: charge it. Add `amount` to the
//     total and return {amount, runningTotal, wasNew=true}.
//   - If this key was already charged: do NOT charge again. Return the
//     amount of the FIRST charge and the current running total, with
//     wasNew=false. totalCharged() must stay the same.
//
// totalCharged() returns the sum of all real (first-time) charges.

#include <string>

class PaymentProcessor {
public:
    struct Result {
        long long amount;
        long long runningTotal;
        bool wasNew;
    };

    Result charge(const std::string& key, long long amount) {
        (void)key;
        (void)amount;
        // TODO
        return Result{0, 0, false};
    }

    long long totalCharged() const {
        // TODO
        return 0;
    }

private:
    // TODO: choose your own representation.
};
