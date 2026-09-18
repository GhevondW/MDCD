# Day 2 — Practical Tasks (Go)

Implement the `TODO`s in the `.go` files below; don't change the
`_test.go` files. Problems 1–3 are the core set; 4–6 are challenge
problems.

Run just this day's tests from `playground/go`:
```
go test ./days/day02/...
```

## Core

### 1. Bounded Stack — Contracts & Invariants
`bounded_stack.go`

`BoundedStack` must enforce its own contract instead of silently
corrupting state:
- `Push(v)` — precondition: not full. If already full, return `ErrFull`
  instead of pushing.
- `Pop()` / `Top()` — precondition: not empty. Return `ErrEmpty` instead.
- Invariant to hold between calls: `0 <= Size() <= capacity`.

### 2. Idempotent Payments — API Guarantees: Idempotency
`idempotent_payments.go`

The same idempotency key must never be charged twice. `Charge()` on a
repeat key must return the *original* charge's amount and the *current*
running total (unchanged) with `WasNew = false` — not charge again.

### 3. Atomic Record Update — API Guarantees: Atomicity
`atomic_record.go`

A record has two fields, `A` and `B`. `Set()` must be all-or-nothing: if
`a + b < 0`, the update is rejected *in full* — neither field may change,
whether the key already existed or not.

## Challenge

### 4. Transactional Key-Value Store — Atomicity
`txn_kv.go`

A string key-value store with **nested** transactions: `Begin()`,
`Commit()`, `Rollback()`. `Commit()` merges the innermost transaction's
writes *and* deletes into its parent transaction (or into the store if it
has no parent); `Rollback()` discards them. Reads always see the innermost
state. The full semantics are in the file's doc comment.

### 5. Idempotency with Expiry & Conflict
`expiring_payments.go`

Idempotency keys that expire on a logical clock: within the TTL a repeat
of a key is a *replay* (same amount) or a *conflict* (different amount) —
neither charges; after the TTL the key is forgotten and the next charge is
new. The full semantics are in the file's doc comment.

### 6. Interval Booker — Invariants
`interval_booker.go`

A booking calendar over half-open intervals `[start, end)` whose invariant
is that recorded bookings never overlap, plus `FirstFree(from, duration)` —
the earliest free slot of a given length at or after `from`. The full
semantics are in the file's doc comment.
