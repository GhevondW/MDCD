# Day 2 — Practical Tasks (Go)

Implement the `TODO`s in the `.go` files below; don't change the
`_test.go` files. Problems 1–3 are the core set; 4–6 are challenge
problems. Each file's doc comment starts with the full problem statement.

Run just this day's tests from `playground/go`:
```
go test ./days/day02/...
```

## Core

### 1. Bounded Stack — Contracts & Invariants
`bounded_stack.go`

A stack of ints with a fixed capacity. `Push` on a full stack must return
`ErrFull`, `Pop`/`Top` on an empty stack must return `ErrEmpty` — instead
of changing anything — and `0 <= Size() <= capacity` must hold at all
times.

### 2. Idempotent Payments — API Guarantees: Idempotency
`idempotent_payments.go`

Every charge comes with a key, and the same key must never be charged
twice. A repeated key returns the first charge's amount and the unchanged
running total.

### 3. Atomic Record Update — API Guarantees: Atomicity
`atomic_record.go`

A record has two fields, `A` and `B`, and is valid only when `a + b >= 0`.
An invalid update must be rejected as a whole: neither field may change,
even if the key already has a record.

## Challenge

### 4. Transactional Key-Value Store — Atomicity
`txn_kv.go`

A key-value store where a transaction can be opened inside another
transaction. `Commit()` moves the innermost transaction's writes and
deletes one level out — not straight into the store; `Rollback()` throws
them away. The full semantics are in the file's doc comment.

### 5. Idempotency with Expiry & Conflict
`expiring_payments.go`

Idempotency keys that are only remembered for a limited time, on a clock
the caller passes in. While a key is remembered, the same amount is a
*replay* and a different amount is a *conflict* — neither charges
anything. When the time runs out, the key is forgotten. The full
semantics are in the file's doc comment.

### 6. Interval Booker — Invariants
`interval_booker.go`

A booking calendar over half-open intervals `[start, end)`. Bookings must
never overlap. Also implement `FirstFree(from, duration)`: the earliest
`start >= from` where a booking of the given length fits. The full
semantics are in the file's doc comment.
