#pragma once
// Day 2 -- API Guarantees: Atomicity
//
// A record has two fields, a and b. set() must be all-or-nothing: if
// a + b < 0 the update is invalid and must be REJECTED IN FULL -- neither
// field may change, whether the key already existed or not. Applying one
// field before validating the other is exactly the bug this catches.

#include <optional>
#include <string>
#include <unordered_map>

struct Record {
    long long a = 0;
    long long b = 0;
};

class RecordStore {
public:
    // TODO: if a + b >= 0, commit BOTH fields together and return true.
    // Otherwise leave the store completely unchanged and return false.
    bool set(const std::string& key, long long a, long long b) {
        (void)key;
        (void)a;
        (void)b;
        return false;
    }

    // TODO: return the record if `key` has been successfully set, or
    // std::nullopt otherwise.
    std::optional<Record> get(const std::string& key) const {
        (void)key;
        return std::nullopt;
    }

private:
    std::unordered_map<std::string, Record> store_;
};
