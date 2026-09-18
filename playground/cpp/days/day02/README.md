# Day 2 — Practical Tasks

Implement the `TODO`s in the headers under `include/day02/`; don't change
the test files. Problems 1–3 are the core set; 4–6 are challenge problems.

Run just this day's tests from `playground/cpp/build`:
```
ctest -C Release -R "^day02_"
```

## Core

### 1. Bounded Stack — Contracts & Invariants
`include/day02/bounded_stack.h`

`BoundedStack` must enforce its own contract instead of silently corrupting
state:
- `push(v)` — precondition: not full. If already full, throw
  `std::runtime_error` instead of pushing.
- `pop()` / `top()` — precondition: not empty. Same idea.
- Invariant to hold between calls: `0 <= size() <= capacity`.

### 2. Idempotent Payments — API Guarantees: Idempotency
`include/day02/idempotent_payments.h`

The same idempotency key must never be charged twice. `charge()` on a
repeat key must return the *original* charge's amount and the *current*
running total (unchanged) with `wasNew = false` — not charge again.

### 3. Atomic Record Update — API Guarantees: Atomicity
`include/day02/atomic_record.h`

A record has two fields, `a` and `b`. `set()` must be all-or-nothing: if
`a + b < 0`, the update is rejected *in full* — neither field may change,
whether the key already existed or not.

## Challenge

### 4. Transactional Key-Value Store — Atomicity
`include/day02/txn_kv.h`

A string key-value store with **nested** transactions: `begin()`,
`commit()`, `rollback()`. `commit()` merges the innermost transaction's
writes *and* deletes into its parent transaction (or into the store if it
has no parent); `rollback()` discards them. Reads always see the innermost
state. The full semantics are in the header.

### 5. Idempotency with Expiry & Conflict
`include/day02/expiring_payments.h`

Idempotency keys that expire on a logical clock: within the TTL a repeat
of a key is a *replay* (same amount) or a *conflict* (different amount) —
neither charges; after the TTL the key is forgotten and the next charge is
new. The exact rules are in the header.

### 6. Interval Booker — Invariants
`include/day02/interval_booker.h`

A booking calendar over half-open intervals `[start, end)` whose invariant
is that recorded bookings never overlap, plus `firstFree(from, duration)` —
the earliest free slot of a given length at or after `from`. Contracts and
semantics are in the header.
