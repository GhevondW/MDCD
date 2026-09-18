package day02;

import java.util.Optional;

/**
 * Day 2 -- Challenge: Transactional Key-Value Store.
 *
 * In-memory string-to-string store with NESTED transactions.
 *
 *   set(key, value)   store a value (in the innermost open transaction,
 *                     or directly if none is open)
 *   get(key)          value currently visible for `key`, or Optional.empty()
 *   remove(key)       delete `key`; returns whether it was visible
 *   begin()           open a new transaction (transactions nest)
 *   commit()          merge the innermost transaction's changes -- writes
 *                     AND deletes -- into its parent transaction (or into
 *                     the store if it has no parent). false if no
 *                     transaction is open.
 *   rollback()        discard the innermost transaction's changes.
 *                     false if no transaction is open.
 *
 * Reads always see the innermost state: a value set (or a key removed)
 * inside an open transaction is visible to get() immediately.
 */
public class TxnKV {
    // TODO: choose your own representation.

    public void set(String key, String value) {
        // TODO
    }

    public Optional<String> get(String key) {
        // TODO
        return Optional.empty();
    }

    public boolean remove(String key) {
        // TODO
        return false;
    }

    public void begin() {
        // TODO
    }

    public boolean commit() {
        // TODO
        return false;
    }

    public boolean rollback() {
        // TODO
        return false;
    }
}
