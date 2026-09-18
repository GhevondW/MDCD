package day02;

import java.util.HashMap;
import java.util.Map;

/**
 * Day 2 -- API Guarantees: Idempotency.
 *
 * The SAME idempotency key must never be charged twice: a repeat of a key
 * you've already seen must return the result of the FIRST charge -- same
 * amount, unchanged running total -- not charge again.
 */
public class PaymentProcessor {
    public record ChargeResult(long amount, long runningTotal, boolean wasNew) {}

    private final Map<String, Long> seen = new HashMap<>();
    private long total = 0;

    // TODO: if `key` was already charged, return the ORIGINAL amount and
    // the CURRENT running total, with wasNew=false -- and don't change
    // totalCharged(). Otherwise remember it, add it to the total, and
    // return it with wasNew=true.
    public ChargeResult charge(String key, long amount) {
        return new ChargeResult(0, 0, false);
    }

    public long totalCharged() {
        // TODO
        return 0;
    }
}
