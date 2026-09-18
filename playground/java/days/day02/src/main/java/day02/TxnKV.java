package day02;

import java.util.Optional;

/**
 * Day 2 -- Challenge: Transactional Key-Value Store.
 *
 * A key-value store (String -> String) with transactions. A transaction
 * can be opened inside another transaction (they nest).
 *
 *   set(key, value)   store a value. If a transaction is open, the change
 *                     belongs to that transaction; otherwise it is applied
 *                     directly to the store.
 *   get(key)          the value visible right now, or Optional.empty() if
 *                     the key does not exist. Changes made inside an open
 *                     transaction are visible immediately.
 *   remove(key)       delete the key. Returns true if the key was visible
 *                     before the call. A delete made inside a transaction
 *                     also belongs to that transaction.
 *   begin()           open a new transaction.
 *   commit()          take everything the innermost transaction changed --
 *                     writes AND deletes -- and move it into the
 *                     transaction one level out (or into the store itself
 *                     if there is no outer transaction). Returns false if
 *                     no transaction is open.
 *   rollback()        throw away everything the innermost transaction
 *                     changed. Returns false if no transaction is open.
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
