# Day 2 — Practical Tasks

Implement the `TODO`s in the headers under `include/day02/`; don't change
the test files. Problems 1–3 are the core set; 4–6 are challenge problems.
Each header starts with the full problem statement.

Run just this day's tests from `playground/cpp/build`:
```
ctest -C Release -R "^day02_"
```

## Core

### 1. Bounded Stack — Contracts & Invariants
`include/day02/bounded_stack.h`

A stack of ints with a fixed capacity. `push` on a full stack and
`pop`/`top` on an empty stack must throw `std::runtime_error` instead of
changing anything, and `0 <= size() <= capacity` must hold at all times.

### 2. Idempotent Payments — API Guarantees: Idempotency
`include/day02/idempotent_payments.h`

Every charge comes with a key, and the same key must never be charged
twice. A repeated key returns the first charge's amount and the unchanged
running total.

### 3. Atomic Record Update — API Guarantees: Atomicity
`include/day02/atomic_record.h`

A record has two fields, `a` and `b`, and is valid only when `a + b >= 0`.
An invalid update must be rejected as a whole: neither field may change,
even if the key already has a record.

## Challenge

### 4. Transactional Key-Value Store — Atomicity
`include/day02/txn_kv.h`

A key-value store where a transaction can be opened inside another
transaction. `commit()` moves the innermost transaction's writes and
deletes one level out — not straight into the store; `rollback()` throws
them away. The full semantics are in the header.

### 5. Idempotency with Expiry & Conflict
`include/day02/expiring_payments.h`

Idempotency keys that are only remembered for a limited time, on a clock
the caller passes in. While a key is remembered, the same amount is a
*replay* and a different amount is a *conflict* — neither charges
anything. When the time runs out, the key is forgotten. The full
semantics are in the header.

### 6. Interval Booker — Invariants
`include/day02/interval_booker.h`

A booking calendar over half-open intervals `[start, end)`. Bookings must
never overlap. Also implement `firstFree(from, duration)`: the earliest
`start >= from` where a booking of the given length fits. The full
semantics are in the header.
