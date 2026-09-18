package day02;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertTrue;

class PaymentProcessorTest {

    @Test
    void firstChargeIsNew() {
        PaymentProcessor p = new PaymentProcessor();
        var r = p.charge("orderA", 100);
        assertTrue(r.wasNew());
        assertEquals(100, r.amount());
        assertEquals(100, r.runningTotal());
        assertEquals(100, p.totalCharged());
    }

    @Test
    void repeatKeyIsNotCharged() {
        PaymentProcessor p = new PaymentProcessor();
        p.charge("orderA", 100);
        var r = p.charge("orderA", 100);
        assertFalse(r.wasNew());
        assertEquals(100, r.amount());
        assertEquals(100, r.runningTotal()); // unchanged -- not double-charged
        assertEquals(100, p.totalCharged());
    }

    @Test
    void differentKeysAccumulate() {
        PaymentProcessor p = new PaymentProcessor();
        p.charge("orderA", 100);
        p.charge("orderB", 250);
        assertEquals(350, p.totalCharged());
    }

    @Test
    void repeatIgnoresNewAmountArgument() {
        // A caller retrying a request might resend the same amount; a repeat
        // must still report the ORIGINAL amount, never a re-derived one.
        PaymentProcessor p = new PaymentProcessor();
        p.charge("x", 5);
        var r = p.charge("x", 5);
        assertEquals(5, r.amount());
        assertEquals(5, p.totalCharged());
    }

    @Test
    void manyRepeatsStillChargeOnce() {
        PaymentProcessor p = new PaymentProcessor();
        p.charge("x", 5);
        p.charge("x", 5);
        p.charge("x", 5);
        p.charge("y", 10);
        assertEquals(15, p.totalCharged());
    }

    @Test
    void repeatIsIdempotentEvenWithOtherChargesBetween() {
        PaymentProcessor p = new PaymentProcessor();
        p.charge("A", 100);
        p.charge("B", 50);
        var r = p.charge("A", 100); // repeat of A, not adjacent to the first
        assertFalse(r.wasNew());
        assertEquals(100, r.amount());
        assertEquals(150, r.runningTotal()); // reflects total so far
        p.charge("C", 25);
        assertEquals(175, p.totalCharged()); // A never double-counted
    }

    @Test
    void zeroAmountChargeIsStillTrackedByKey() {
        PaymentProcessor p = new PaymentProcessor();
        var first = p.charge("free-trial", 0);
        assertTrue(first.wasNew());
        assertEquals(0, first.amount());
        var repeat = p.charge("free-trial", 0);
        assertFalse(repeat.wasNew());
        assertEquals(0, p.totalCharged());
    }
}
