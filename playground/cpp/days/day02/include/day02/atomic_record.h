#pragma once
// Day 2 -- API Guarantees: Atomicity
//
// A record has two fields, a and b. A record is valid only when
// a + b >= 0.
//
// set(key, a, b) must be all-or-nothing:
//   - If a + b >= 0: store both fields together and return true.
//   - If a + b < 0: the update is invalid. Reject the WHOLE update -- do
//     not change either field -- and return false. This also holds when
//     the key already has a record: the old record must stay exactly as
//     it was.
//
// get(key) returns the stored record, or std::nullopt if this key was
// never successfully set.

#include <optional>
#include <string>

struct Record {
    long long a = 0;
    long long b = 0;
};

class RecordStore {
public:
    bool set(const std::string& key, long long a, long long b) {
        (void)key;
        (void)a;
        (void)b;
        // TODO
        return false;
    }

    std::optional<Record> get(const std::string& key) const {
        (void)key;
        // TODO
        return std::nullopt;
    }

private:
    // TODO: choose your own representation.
};
