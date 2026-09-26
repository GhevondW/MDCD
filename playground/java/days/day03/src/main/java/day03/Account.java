package day03;

/**
 * Day 3 -- Critical sections: check-then-act.
 *
 * A bank account shared by many threads -- think of one card used by
 * several people at the same moment. The balance must never go below
 * zero.
 *
 *   Account(initial)   start with `initial`. initial must be >= 0; if
 *                      not, throw IllegalArgumentException.
 *   deposit(amount)    add `amount`. amount must be > 0; if not, throw
 *                      IllegalArgumentException and change nothing.
 *   withdraw(amount)   amount must be > 0; if not, throw
 *                      IllegalArgumentException and change nothing.
 *                      If balance >= amount: subtract it, return true.
 *                      Otherwise return false and change nothing.
 *   balance()          the current balance.
 *
 * deposit, withdraw and balance may be called by many threads at the
 * same time.
 *
 *     if (balance >= amount) balance -= amount;
 *
 * is a check and an act -- another thread can slip in between them.
 */
public class Account {
    // TODO: choose your own representation.

    public Account(long initial) {
        // TODO
    }

    public void deposit(long amount) {
        // TODO
    }

    public boolean withdraw(long amount) {
        // TODO
        return false;
    }

    public long balance() {
        // TODO
        return 0;
    }
}
