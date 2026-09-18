#pragma once
// Day 2 -- Challenge: Idempotency with Expiry & Conflict
//
// A payment processor whose idempotency keys EXPIRE. Time is a logical
// clock: every call takes `now` (non-decreasing across calls); a key
// recorded at time T with time-to-live `ttl` is expired once now >= T + ttl.
//
// charge(key, amount, now):
//   - key unknown, or its record expired  -> charge: add `amount` to the
//     total, record (key, amount) with expiry now + ttl, return New.
//   - key active, same amount             -> return Replay with the
//     originally recorded amount and the current total. Nothing changes --
//     not the total, not the recorded amount, not the expiry.
//   - key active, different amount        -> return Conflict with the
//     recorded amount and the current total. Nothing changes.
//
// totalCharged()   sum of all charges that returned New.
// activeKeys(now)  number of recorded keys not yet expired at `now`.

#include <cstddef>
#include <string>

class ExpiringPaymentProcessor {
public:
    enum class Status { New, Replay, Conflict };

    struct Outcome {
        Status status;
        long long amount;        // New: charged; Replay/Conflict: recorded
        long long runningTotal;
    };

    explicit ExpiringPaymentProcessor(long long ttl) : ttl_(ttl) {}

    Outcome charge(const std::string& key, long long amount, long long now) {
        (void)key;
        (void)amount;
        (void)now;
        // TODO
        return Outcome{Status::New, 0, 0};
    }

    long long totalCharged() const {
        // TODO
        return 0;
    }

    std::size_t activeKeys(long long now) const {
        (void)now;
        // TODO
        return 0;
    }

private:
    long long ttl_;
    // TODO: choose your own representation.
};
