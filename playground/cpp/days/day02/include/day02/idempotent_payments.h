#pragma once
// Day 2 -- API Guarantees: Idempotency
//
// The SAME idempotency key must never be charged twice: a repeat of a key
// you've already seen must return the result of the FIRST charge -- same
// amount, unchanged running total -- not charge again.

#include <string>
#include <unordered_map>

class PaymentProcessor {
public:
    struct Result {
        long long amount;
        long long runningTotal;
        bool wasNew;
    };

    // TODO: if `key` was already charged, return the ORIGINAL amount and
    // the CURRENT running total, with wasNew=false -- and don't change
    // totalCharged(). Otherwise remember it, add it to the total, and
    // return it with wasNew=true.
    Result charge(const std::string& key, long long amount) {
        (void)key;
        (void)amount;
        return Result{0, 0, false};
    }

    long long totalCharged() const {
        // TODO
        return 0;
    }

private:
    std::unordered_map<std::string, long long> seen_;
    long long total_ = 0;
};
