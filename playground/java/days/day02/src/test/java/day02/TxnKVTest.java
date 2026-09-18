package day02;

import org.junit.jupiter.api.Test;

import java.util.Optional;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertTrue;

class TxnKVTest {

    @Test
    void setThenGet() {
        TxnKV kv = new TxnKV();
        kv.set("a", "1");
        assertEquals(Optional.of("1"), kv.get("a"));
    }

    @Test
    void getMissingKeyReturnsNothing() {
        TxnKV kv = new TxnKV();
        assertTrue(kv.get("nope").isEmpty());
    }

    @Test
    void overwriteReplacesValue() {
        TxnKV kv = new TxnKV();
        kv.set("a", "1");
        kv.set("a", "2");
        assertEquals(Optional.of("2"), kv.get("a"));
    }

    @Test
    void removeReturnsTrueOnlyIfKeyWasVisible() {
        TxnKV kv = new TxnKV();
        kv.set("a", "1");
        assertTrue(kv.remove("a"));
        assertTrue(kv.get("a").isEmpty());
        assertFalse(kv.remove("a"));
        assertFalse(kv.remove("never-set"));
    }

    @Test
    void commitWithoutOpenTransactionFails() {
        TxnKV kv = new TxnKV();
        assertFalse(kv.commit());
    }

    @Test
    void rollbackWithoutOpenTransactionFails() {
        TxnKV kv = new TxnKV();
        assertFalse(kv.rollback());
    }

    @Test
    void setInsideTransactionIsVisibleBeforeCommit() {
        TxnKV kv = new TxnKV();
        kv.begin();
        kv.set("a", "1");
        assertEquals(Optional.of("1"), kv.get("a"));
    }

    @Test
    void rollbackDiscardsEveryChangeOfTheTransaction() {
        TxnKV kv = new TxnKV();
        kv.set("a", "old");
        kv.begin();
        kv.set("a", "new");
        kv.set("b", "1");
        assertTrue(kv.rollback());
        assertEquals(Optional.of("old"), kv.get("a"));
        assertTrue(kv.get("b").isEmpty());
    }

    @Test
    void commitMakesChangesPermanent() {
        TxnKV kv = new TxnKV();
        kv.begin();
        kv.set("a", "1");
        assertTrue(kv.commit());
        assertEquals(Optional.of("1"), kv.get("a"));
        assertFalse(kv.rollback()); // nothing left open to roll back
        assertEquals(Optional.of("1"), kv.get("a"));
    }

    @Test
    void removeInsideTransactionHidesOuterValue() {
        TxnKV kv = new TxnKV();
        kv.set("a", "base");
        kv.begin();
        assertTrue(kv.remove("a"));
        assertTrue(kv.get("a").isEmpty()); // gone as seen from inside
    }

    @Test
    void rollbackRestoresValueRemovedInsideTransaction() {
        TxnKV kv = new TxnKV();
        kv.set("a", "base");
        kv.begin();
        kv.remove("a");
        kv.rollback();
        assertEquals(Optional.of("base"), kv.get("a"));
    }

    @Test
    void committedRemoveDeletesFromTheStore() {
        TxnKV kv = new TxnKV();
        kv.set("a", "base");
        kv.begin();
        kv.remove("a");
        assertTrue(kv.commit());
        assertTrue(kv.get("a").isEmpty());
    }

    @Test
    void nestedInnerRollbackKeepsOuterChanges() {
        TxnKV kv = new TxnKV();
        kv.begin();
        kv.set("outer", "1");
        kv.begin();
        kv.set("inner", "2");
        assertTrue(kv.rollback()); // discards only the inner transaction
        assertTrue(kv.get("inner").isEmpty());
        assertEquals(Optional.of("1"), kv.get("outer"));
    }

    @Test
    void innerCommitMergesIntoParentNotIntoStore() {
        TxnKV kv = new TxnKV();
        kv.begin();
        kv.begin();
        kv.set("a", "1");
        assertTrue(kv.commit()); // inner commit: `a` now belongs to the OUTER txn
        assertEquals(Optional.of("1"), kv.get("a"));
        assertTrue(kv.rollback()); // outer rollback must take `a` down with it
        assertTrue(kv.get("a").isEmpty());
    }

    @Test
    void innerCommittedRemoveStaysInsideParent() {
        TxnKV kv = new TxnKV();
        kv.set("a", "base");
        kv.begin();
        kv.begin();
        kv.remove("a");
        assertTrue(kv.commit()); // remove merged into the outer txn
        assertTrue(kv.get("a").isEmpty());
        assertTrue(kv.rollback()); // outer rollback cancels the remove too
        assertEquals(Optional.of("base"), kv.get("a"));
    }

    @Test
    void setAfterRemoveInSameTransactionResurrectsKey() {
        TxnKV kv = new TxnKV();
        kv.set("a", "base");
        kv.begin();
        kv.remove("a");
        kv.set("a", "reborn");
        assertEquals(Optional.of("reborn"), kv.get("a"));
        kv.commit();
        assertEquals(Optional.of("reborn"), kv.get("a"));
    }

    @Test
    void valueIsVisibleThroughMultipleUntouchedLayers() {
        TxnKV kv = new TxnKV();
        kv.set("a", "base");
        kv.begin();
        kv.begin();
        kv.begin();
        assertEquals(Optional.of("base"), kv.get("a"));
    }
}
