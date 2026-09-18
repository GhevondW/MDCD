package day02;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertTrue;

class RecordStoreTest {

    @Test
    void getOnUnknownKeyIsEmpty() {
        RecordStore store = new RecordStore();
        assertTrue(store.get("ghost").isEmpty());
    }

    @Test
    void validSetCommitsBothFields() {
        RecordStore store = new RecordStore();
        assertTrue(store.set("acct1", 10, 5));
        var r = store.get("acct1");
        assertTrue(r.isPresent());
        assertEquals(10, r.get().a());
        assertEquals(5, r.get().b());
    }

    @Test
    void invalidSetIsRejected() {
        RecordStore store = new RecordStore();
        assertFalse(store.set("acct2", -5, -5));
        assertTrue(store.get("acct2").isEmpty());
    }

    @Test
    void rejectedSetLeavesExistingRecordUntouched() {
        // The core atomicity check: a rejected update on an EXISTING key
        // must not partially apply -- field a must not change even though
        // it would have been "valid" applied in isolation.
        RecordStore store = new RecordStore();
        store.set("acct1", 10, 5);
        assertFalse(store.set("acct1", -20, 3)); // -20 + 3 < 0 -> reject
        var r = store.get("acct1");
        assertTrue(r.isPresent());
        assertEquals(10, r.get().a()); // unchanged, not -20
        assertEquals(5, r.get().b()); // unchanged
    }

    @Test
    void zeroSumIsValid() {
        RecordStore store = new RecordStore();
        assertTrue(store.set("a", 0, 0));
        var r = store.get("a");
        assertTrue(r.isPresent());
        assertEquals(0, r.get().a());
        assertEquals(0, r.get().b());
    }

    @Test
    void negativeFieldWithCompensatingPositiveIsValid() {
        // The rule is a + b >= 0, NOT "both fields individually
        // non-negative" -- a single negative field with a large enough
        // positive one is fine.
        RecordStore store = new RecordStore();
        assertTrue(store.set("acct3", -5, 10));
        var r = store.get("acct3");
        assertTrue(r.isPresent());
        assertEquals(-5, r.get().a());
        assertEquals(10, r.get().b());
    }

    @Test
    void validSetOverwritesPreviousValidRecord() {
        RecordStore store = new RecordStore();
        store.set("acct1", 10, 5);
        assertTrue(store.set("acct1", 1, 2)); // second valid set replaces it
        var r = store.get("acct1");
        assertEquals(1, r.get().a());
        assertEquals(2, r.get().b());
    }

    @Test
    void keysAreIndependent() {
        RecordStore store = new RecordStore();
        store.set("acct1", 10, 5);
        store.set("acct2", 1, 1);
        assertFalse(store.set("acct2", -100, 0)); // reject on acct2 only
        var r1 = store.get("acct1");
        var r2 = store.get("acct2");
        assertEquals(10, r1.get().a());
        assertEquals(5, r1.get().b());
        assertEquals(1, r2.get().a()); // acct2's rejection didn't touch acct1
        assertEquals(1, r2.get().b());
    }
}
