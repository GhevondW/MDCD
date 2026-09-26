package day03;

/**
 * Day 3 -- Critical sections: check-then-act, and two locks at once.
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
 *   transfer(to, amount)
 *                      move `amount` from this account to `to`, as ONE
 *                      step: no thread may ever see the money gone from
 *                      one account but not yet in the other. amount must
 *                      be > 0 and `to` must be a different account; if
 *                      not, throw IllegalArgumentException. If this
 *                      balance < amount, return false and change nothing.
 *   total(a, b)        (static) a's balance plus b's balance, read as ONE
 *                      step: every transfer between a and b happens either
 *                      completely before it or completely after it. a and
 *                      b must be different accounts; if not, throw
 *                      IllegalArgumentException.
 *
 * Everything except the constructor may be called by many threads at the
 * same time.
 *
 *     if (balance >= amount) balance -= amount;
 *
 * is a check and an act -- another thread can slip in between them.
 *
 * transfer and total need two accounts locked at once. Two threads that
 * lock the same two accounts in opposite orders wait for each other
 * forever -- the slides' "Two transfers, two locks". Hint: if every
 * thread locks any two accounts in the same order (say, by a number each
 * account gets when it is created), that cannot happen.
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

    public boolean transfer(Account to, long amount) {
        // TODO
        return false;
    }

    public static long total(Account a, Account b) {
        // TODO
        return 0;
    }
}
