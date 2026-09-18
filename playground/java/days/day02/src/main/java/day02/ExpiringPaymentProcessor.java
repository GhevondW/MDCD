package day02;

/**
 * Day 2 -- Challenge: Idempotency with Expiry & Conflict.
 *
 * Like the idempotent payment processor, but keys are only remembered for
 * a limited time. There is no real clock: every call receives the current
 * time `now` as a plain number, and the numbers never go down from call
 * to call. A key recorded at time T is remembered until time T + ttl;
 * from now == T + ttl on, it is expired (forgotten).
 *
 * charge(key, amount, now) -- three cases:
 *   - The key is unknown, or its record has expired: this is a NEW
 *     charge. Add `amount` to the total and remember (key, amount) until
 *     now + ttl.
 *   - The key is still remembered and `amount` equals the recorded
 *     amount: this is a REPLAY (the caller sent the same request twice).
 *     Return the recorded amount and the current total. Change NOTHING:
 *     not the total, not the recorded amount, not the expiry time.
 *   - The key is still remembered but `amount` is different: this is a
 *     CONFLICT (the same key was used for a different request -- a
 *     mistake). Return the recorded amount and the current total.
 *     Change nothing.
 *
 * totalCharged()   sum of all NEW charges.
 * activeKeys(now)  how many recorded keys are not yet expired at `now`.
 */
public class ExpiringPaymentProcessor {
    public enum Status { NEW, REPLAY, CONFLICT }

    // amount -- NEW: charged; REPLAY/CONFLICT: recorded
    public record Outcome(Status status, long amount, long runningTotal) {}

    // TODO: choose your own representation.

    public ExpiringPaymentProcessor(long ttl) {
        // TODO
    }

    public Outcome charge(String key, long amount, long now) {
        // TODO
        return new Outcome(Status.NEW, 0, 0);
    }

    public long totalCharged() {
        // TODO
        return 0;
    }

    public int activeKeys(long now) {
        // TODO
        return 0;
    }
}
