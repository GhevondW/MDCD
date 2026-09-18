#pragma once
// Day 2 -- Challenge: Idempotency with Expiry & Conflict
//
// Like the idempotent payment processor, but keys are only remembered for
// a limited time. There is no real clock: every call receives the current
// time `now` as a plain number, and the numbers never go down from call
// to call. A key recorded at time T is remembered until time T + ttl;
// from now == T + ttl on, it is expired (forgotten).
//
// charge(key, amount, now) -- three cases:
//   - The key is unknown, or its record has expired: this is a NEW
//     charge. Add `amount` to the total and remember (key, amount) until
//     now + ttl.
//   - The key is still remembered and `amount` equals the recorded
//     amount: this is a REPLAY (the caller sent the same request twice).
//     Return the recorded amount and the current total. Change NOTHING:
//     not the total, not the recorded amount, not the expiry time.
//   - The key is still remembered but `amount` is different: this is a
//     CONFLICT (the same key was used for a different request -- a
//     mistake). Return the recorded amount and the current total.
//     Change nothing.
//
// totalCharged()   sum of all NEW charges.
// activeKeys(now)  how many recorded keys are not yet expired at `now`.

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

    explicit ExpiringPaymentProcessor(long long ttl) {
        (void)ttl;
        // TODO
    }

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
    // TODO: choose your own representation.
};
