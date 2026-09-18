package day02;

import java.util.HashMap;
import java.util.Map;
import java.util.Optional;

/**
 * Day 2 -- API Guarantees: Atomicity.
 *
 * A record has two fields, a and b. set() must be all-or-nothing: if
 * a + b < 0 the update is invalid and must be REJECTED IN FULL -- neither
 * field may change, whether the key already existed or not. Applying one
 * field before validating the other is exactly the bug this catches.
 */
public class RecordStore {
    public record Fields(long a, long b) {}

    private final Map<String, Fields> store = new HashMap<>();

    // TODO: if a + b >= 0, commit BOTH fields together and return true.
    // Otherwise leave the store completely unchanged and return false.
    public boolean set(String key, long a, long b) {
        return false;
    }

    // TODO: return the record if `key` has been successfully set, or
    // Optional.empty() otherwise.
    public Optional<Fields> get(String key) {
        return Optional.empty();
    }
}
