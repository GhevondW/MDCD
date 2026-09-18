# Day 2 — Practical Tasks (Java)

Implement the `TODO`s in `src/main/java/day02/`; don't change
`src/test/`. Problems 1–3 are the core set; 4–6 are challenge problems.

Run just this day's tests from `playground/java`:
```
mvn -pl days/day02 test
```

## Core

### 1. Bounded Stack — Contracts & Invariants
`src/main/java/day02/BoundedStack.java`

`BoundedStack` must enforce its own contract instead of silently
corrupting state:
- `push(v)` — precondition: not full. If already full, throw
  `IllegalStateException` instead of pushing.
- `pop()` / `top()` — precondition: not empty. Throw
  `NoSuchElementException` instead.
- Invariant to hold between calls: `0 <= size() <= capacity`.

### 2. Idempotent Payments — API Guarantees: Idempotency
`src/main/java/day02/PaymentProcessor.java`

The same idempotency key must never be charged twice. `charge()` on a
repeat key must return the *original* charge's amount and the *current*
running total (unchanged) with `wasNew = false` — not charge again.

### 3. Atomic Record Update — API Guarantees: Atomicity
`src/main/java/day02/RecordStore.java`

A record has two fields, `a` and `b`. `set()` must be all-or-nothing: if
`a + b < 0`, the update is rejected *in full* — neither field may change,
whether the key already existed or not.

## Challenge

### 4. Transactional Key-Value Store — Atomicity
`src/main/java/day02/TxnKV.java`

A string key-value store with **nested** transactions: `begin()`,
`commit()`, `rollback()`. `commit()` merges the innermost transaction's
writes *and* deletes into its parent transaction (or into the store if it
has no parent); `rollback()` discards them. Reads always see the innermost
state. The full semantics are in the class Javadoc.

### 5. Idempotency with Expiry & Conflict
`src/main/java/day02/ExpiringPaymentProcessor.java`

Idempotency keys that expire on a logical clock: within the TTL a repeat
of a key is a *replay* (same amount) or a *conflict* (different amount) —
neither charges; after the TTL the key is forgotten and the next charge is
new. The full semantics are in the class Javadoc.

### 6. Interval Booker — Invariants
`src/main/java/day02/IntervalBooker.java`

A booking calendar over half-open intervals `[start, end)` whose invariant
is that recorded bookings never overlap, plus `firstFree(from, duration)` —
the earliest free slot of a given length at or after `from`. The full
semantics are in the class Javadoc.
