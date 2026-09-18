package day02;

import java.util.Optional;

/**
 * Day 2 -- API Guarantees: Atomicity.
 *
 * A record has two fields, a and b. A record is valid only when
 * a + b >= 0.
 *
 * set(key, a, b) must be all-or-nothing:
 *   - If a + b >= 0: store both fields together and return true.
 *   - If a + b < 0: the update is invalid. Reject the WHOLE update -- do
 *     not change either field -- and return false. This also holds when
 *     the key already has a record: the old record must stay exactly as
 *     it was.
 *
 * get(key) returns the stored record, or Optional.empty() if this key
 * was never successfully set.
 */
public class RecordStore {
    public record Fields(long a, long b) {}

    // TODO: choose your own representation.

    public boolean set(String key, long a, long b) {
        // TODO
        return false;
    }

    public Optional<Fields> get(String key) {
        // TODO
        return Optional.empty();
    }
}
