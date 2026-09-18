#pragma once
// Day 2 -- Challenge: Transactional Key-Value Store
//
// A key-value store (string -> string) with transactions. A transaction
// can be opened inside another transaction (they nest).
//
//   set(key, value)   store a value. If a transaction is open, the change
//                     belongs to that transaction; otherwise it is applied
//                     directly to the store.
//   get(key)          the value visible right now, or std::nullopt if the
//                     key does not exist. Changes made inside an open
//                     transaction are visible immediately.
//   erase(key)        delete the key. Returns true if the key was visible
//                     before the call. A delete made inside a transaction
//                     also belongs to that transaction.
//   begin()           open a new transaction.
//   commit()          take everything the innermost transaction changed --
//                     writes AND deletes -- and move it into the
//                     transaction one level out (or into the store itself
//                     if there is no outer transaction). Returns false if
//                     no transaction is open.
//   rollback()        throw away everything the innermost transaction
//                     changed. Returns false if no transaction is open.

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
