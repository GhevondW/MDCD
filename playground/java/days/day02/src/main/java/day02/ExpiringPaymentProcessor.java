package day02;

/**
 * Day 2 -- Challenge: Idempotency with Expiry & Conflict.
 *
 * A payment processor whose idempotency keys EXPIRE. Time is a logical
 * clock: every call takes `now` (non-decreasing across calls); a key
 * recorded at time T with time-to-live `ttl` is expired once now >= T + ttl.
 *
 * charge(key, amount, now):
 *   - key unknown, or its record expired  -> charge: add `amount` to the
 *     total, record (key, amount) with expiry now + ttl, return NEW.
 *   - key active, same amount             -> return REPLAY with the
 *     originally recorded amount and the current total. Nothing changes --
 *     not the total, not the recorded amount, not the expiry.
 *   - key active, different amount        -> return CONFLICT with the
 *     recorded amount and the current total. Nothing changes.
 *
 * totalCharged()   sum of all charges that returned NEW.
 * activeKeys(now)  number of recorded keys not yet expired at `now`.
 */
public class ExpiringPaymentProcessor {
    public enum Status { NEW, REPLAY, CONFLICT }

    // amount -- NEW: charged; REPLAY/CONFLICT: recorded
    public record Outcome(Status status, long amount, long runningTotal) {}

    private final long ttl;
    // TODO: choose your own representation.

    public ExpiringPaymentProcessor(long ttl) {
        this.ttl = ttl;
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
