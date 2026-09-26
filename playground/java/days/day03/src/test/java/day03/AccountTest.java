package day03;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.Timeout;

import java.util.concurrent.atomic.AtomicBoolean;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.concurrent.atomic.AtomicLong;

import static day03.RunTogether.runTogether;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

@Timeout(value = 60, threadMode = Timeout.ThreadMode.SEPARATE_THREAD)
class AccountTest {

    @Test
    void startsWithInitialBalance() {
        Account a = new Account(100);
        assertEquals(100, a.balance());
    }

    @Test
    void negativeInitialBalanceThrows() {
        assertThrows(IllegalArgumentException.class, () -> new Account(-1));
    }

    @Test
    void depositAddsAndWithdrawSubtracts() {
        Account a = new Account(0);
        a.deposit(50);
        assertTrue(a.withdraw(20));
        assertEquals(30, a.balance());
    }

    @Test
    void withdrawExactBalanceSucceeds() {
        Account a = new Account(40);
        assertTrue(a.withdraw(40));
        assertEquals(0, a.balance());
    }

    @Test
    void withdrawMoreThanBalanceFailsAndChangesNothing() {
        Account a = new Account(30);
        assertFalse(a.withdraw(31));
        assertEquals(30, a.balance());
    }

    @Test
    void nonPositiveAmountsThrowAndChangeNothing() {
        Account a = new Account(10);
        assertThrows(IllegalArgumentException.class, () -> a.deposit(0));
        assertThrows(IllegalArgumentException.class, () -> a.deposit(-5));
        assertThrows(IllegalArgumentException.class, () -> a.withdraw(0));
        assertThrows(IllegalArgumentException.class, () -> a.withdraw(-5));
        assertEquals(10, a.balance());
    }

    @Test
    void concurrentWithdrawalsNeverOverdraw() throws Exception {
        // 8 threads try to take 40,000 in total from an account holding
        // 10,000. Exactly 10,000 withdrawals of 1 can succeed. The race only
        // matters at the moment the money runs out, so play it 20 times.
        int threads = 8;
        int triesPerThread = 5000;
        for (int round = 0; round < 20; round++) {
            Account a = new Account(10_000);
            AtomicLong succeeded = new AtomicLong();
            runTogether(threads, t -> {
                for (int i = 0; i < triesPerThread; i++) {
                    if (a.withdraw(1)) succeeded.incrementAndGet();
                }
            });

            int r = round;
            assertEquals(10_000, succeeded.get(), () -> "round " + r + ": withdrawals that succeeded");
            assertEquals(0, a.balance(), () -> "round " + r + ": balance");
        }
    }

    @Test
    void concurrentDepositsAreNeverLost() throws Exception {
        int threads = 8;
        int perThread = 20_000;
        Account a = new Account(0);
        runTogether(threads, t -> {
            for (int i = 0; i < perThread; i++) a.deposit(1);
        });

        assertEquals(threads * perThread, a.balance(), "a deposit was lost");
    }

    @Test
    void concurrentDepositsAndWithdrawalsBalanceOut() throws Exception {
        // 4 threads deposit 3 at a time, 4 threads withdraw 2 at a time, and
        // one watcher checks that the balance is never seen below zero.
        int perThread = 20_000;
        Account a = new Account(0);
        AtomicLong withdrawn = new AtomicLong();
        AtomicBoolean sawNegative = new AtomicBoolean(false);
        AtomicInteger workersLeft = new AtomicInteger(8);
        runTogether(9, t -> {
            if (t == 8) { // the watcher
                while (workersLeft.get() > 0) {
                    if (a.balance() < 0) sawNegative.set(true);
                }
                return;
            }
            try {
                for (int i = 0; i < perThread; i++) {
                    if (t < 4) {
                        a.deposit(3);
                    } else if (a.withdraw(2)) {
                        withdrawn.addAndGet(2);
                    }
                }
            } finally {
                workersLeft.decrementAndGet(); // even if this worker failed
            }
        });

        assertFalse(sawNegative.get(), "the balance was seen below zero");
        assertEquals(4L * perThread * 3 - withdrawn.get(), a.balance());
    }
}
