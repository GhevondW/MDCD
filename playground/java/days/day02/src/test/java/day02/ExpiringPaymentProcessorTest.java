package day02;

import org.junit.jupiter.api.Test;

import day02.ExpiringPaymentProcessor.Status;

import static org.junit.jupiter.api.Assertions.assertEquals;

class ExpiringPaymentProcessorTest {

    @Test
    void firstChargeIsNew() {
        ExpiringPaymentProcessor p = new ExpiringPaymentProcessor(10);
        var r = p.charge("k", 100, 0);
        assertEquals(Status.NEW, r.status());
        assertEquals(100, r.amount());
        assertEquals(100, r.runningTotal());
        assertEquals(100, p.totalCharged());
    }

    @Test
    void sameAmountWithinTtlIsReplay() {
        ExpiringPaymentProcessor p = new ExpiringPaymentProcessor(10);
        p.charge("k", 100, 0);
        var r = p.charge("k", 100, 5);
        assertEquals(Status.REPLAY, r.status());
        assertEquals(100, r.amount());
        assertEquals(100, r.runningTotal());
        assertEquals(100, p.totalCharged());
    }

    @Test
    void differentAmountWithinTtlIsConflict() {
        ExpiringPaymentProcessor p = new ExpiringPaymentProcessor(10);
        p.charge("k", 100, 0);
        var r = p.charge("k", 250, 5);
        assertEquals(Status.CONFLICT, r.status());
        assertEquals(100, r.amount());          // the RECORDED amount, not 250
        assertEquals(100, r.runningTotal());
        assertEquals(100, p.totalCharged());    // nothing was charged
    }

    @Test
    void expiredKeyIsChargedAgain() {
        ExpiringPaymentProcessor p = new ExpiringPaymentProcessor(10);
        p.charge("k", 100, 0);
        var r = p.charge("k", 100, 20);
        assertEquals(Status.NEW, r.status());
        assertEquals(200, p.totalCharged());
    }

    @Test
    void expiryBoundaryIsInclusive() {
        // Recorded at 0 with ttl 10 -> expired exactly at now == 10.
        ExpiringPaymentProcessor p = new ExpiringPaymentProcessor(10);
        p.charge("k", 100, 0);
        assertEquals(Status.REPLAY, p.charge("k", 100, 9).status());
        assertEquals(Status.NEW, p.charge("k", 100, 10).status());
        assertEquals(200, p.totalCharged());
    }

    @Test
    void replayDoesNotExtendExpiry() {
        ExpiringPaymentProcessor p = new ExpiringPaymentProcessor(10);
        p.charge("k", 100, 0);
        assertEquals(Status.REPLAY, p.charge("k", 100, 9).status());
        // If the replay at 9 had refreshed the record, the key would still be
        // active at 10. It must not be.
        assertEquals(Status.NEW, p.charge("k", 100, 10).status());
    }

    @Test
    void conflictChangesNothing() {
        ExpiringPaymentProcessor p = new ExpiringPaymentProcessor(10);
        p.charge("k", 100, 0);
        p.charge("k", 250, 5);                        // conflict
        var r = p.charge("k", 100, 9);                // original amount again
        assertEquals(Status.REPLAY, r.status());      // record still 100
        assertEquals(100, r.amount());
        assertEquals(100, p.totalCharged());
        assertEquals(Status.NEW, p.charge("k", 100, 10).status()); // expiry not shifted
    }

    @Test
    void afterExpiryDifferentAmountIsNewNotConflict() {
        ExpiringPaymentProcessor p = new ExpiringPaymentProcessor(10);
        p.charge("k", 100, 0);
        var r = p.charge("k", 250, 15);
        assertEquals(Status.NEW, r.status());
        assertEquals(250, r.amount());
        assertEquals(350, p.totalCharged());
        // The record was replaced: 250 is now the recorded amount.
        assertEquals(Status.REPLAY, p.charge("k", 250, 16).status());
        assertEquals(Status.CONFLICT, p.charge("k", 100, 16).status());
    }

    @Test
    void keysAreIndependent() {
        ExpiringPaymentProcessor p = new ExpiringPaymentProcessor(10);
        p.charge("a", 100, 0);
        p.charge("b", 50, 8);
        assertEquals(150, p.totalCharged());
        assertEquals(Status.NEW, p.charge("a", 100, 12).status());    // a expired
        assertEquals(Status.REPLAY, p.charge("b", 50, 12).status());  // b still active
    }

    @Test
    void activeKeysCountsOnlyUnexpired() {
        ExpiringPaymentProcessor p = new ExpiringPaymentProcessor(10);
        p.charge("a", 1, 0);   // expires at 10
        p.charge("b", 2, 5);   // expires at 15
        p.charge("c", 3, 8);   // expires at 18
        assertEquals(3, p.activeKeys(0));
        assertEquals(2, p.activeKeys(12));
        assertEquals(1, p.activeKeys(15));
        assertEquals(0, p.activeKeys(18));
    }

    @Test
    void zeroTtlMeansEveryChargeIsNew() {
        ExpiringPaymentProcessor p = new ExpiringPaymentProcessor(0);
        assertEquals(Status.NEW, p.charge("k", 100, 5).status());
        assertEquals(Status.NEW, p.charge("k", 100, 5).status());
        assertEquals(200, p.totalCharged());
        assertEquals(0, p.activeKeys(5));
    }

    @Test
    void replayReportsCurrentTotalNotTotalAtFirstCharge() {
        ExpiringPaymentProcessor p = new ExpiringPaymentProcessor(100);
        p.charge("a", 100, 0);
        p.charge("b", 50, 1);
        var r = p.charge("a", 100, 2);
        assertEquals(Status.REPLAY, r.status());
        assertEquals(150, r.runningTotal());
    }

    @Test
    void totalIsSumOfNewChargesOnly() {
        ExpiringPaymentProcessor p = new ExpiringPaymentProcessor(10);
        p.charge("a", 100, 0);   // new
        p.charge("a", 100, 1);   // replay
        p.charge("a", 999, 2);   // conflict
        p.charge("b", 50, 3);    // new
        p.charge("a", 25, 12);   // expired -> new
        assertEquals(175, p.totalCharged());
    }
}
