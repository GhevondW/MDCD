package day02;

/**
 * Day 2 -- API Guarantees: Idempotency.
 *
 * Every charge request comes with a key. The same key must never be
 * charged twice.
 *
 * charge(key, amount):
 *   - If this key was never seen before: charge it. Add `amount` to the
 *     total and return new ChargeResult(amount, runningTotal, true).
 *   - If this key was already charged: do NOT charge again. Return the
 *     amount of the FIRST charge and the current running total, with
 *     wasNew=false. totalCharged() must stay the same.
 *
 * totalCharged() returns the sum of all real (first-time) charges.
 */
public class PaymentProcessor {
    public record ChargeResult(long amount, long runningTotal, boolean wasNew) {}

    // TODO: choose your own representation.

    public ChargeResult charge(String key, long amount) {
        // TODO
        return new ChargeResult(0, 0, false);
    }

    public long totalCharged() {
        // TODO
        return 0;
    }
}
