#pragma once
// Day 2 -- Challenge: Transactional Key-Value Store
//
// In-memory string->string store with NESTED transactions.
//
//   set(key, value)   store a value (in the innermost open transaction,
//                     or directly if none is open)
//   get(key)          value currently visible for `key`, or std::nullopt
//   erase(key)        delete `key`; returns whether it was visible
//   begin()           open a new transaction (transactions nest)
//   commit()          merge the innermost transaction's changes -- writes
//                     AND deletes -- into its parent transaction (or into
//                     the store if it has no parent). false if no
//                     transaction is open.
//   rollback()        discard the innermost transaction's changes.
//                     false if no transaction is open.
//
// Reads always see the innermost state: a value set (or a key erased)
// inside an open transaction is visible to get() immediately.

#include <optional>
#include <string>

class TxnKV {
public:
    void set(const std::string& key, const std::string& value) {
        (void)key;
        (void)value;
        // TODO
    }

    std::optional<std::string> get(const std::string& key) const {
        (void)key;
        // TODO
        return std::nullopt;
    }

    bool erase(const std::string& key) {
        (void)key;
        // TODO
        return false;
    }

    void begin() {
        // TODO
    }

    bool commit() {
        // TODO
        return false;
    }

    bool rollback() {
        // TODO
        return false;
    }

private:
    // TODO: choose your own representation.
};
