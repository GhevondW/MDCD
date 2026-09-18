# Day 2 — Practical Tasks (Java)

Implement the `TODO`s in `src/main/java/day02/`; don't change
`src/test/`. Problems 1–3 are the core set; 4–6 are challenge problems.
Each class's Javadoc starts with the full problem statement.

Run just this day's tests from `playground/java`:
```
mvn -pl days/day02 test
```

## Core

### 1. Bounded Stack — Contracts & Invariants
`src/main/java/day02/BoundedStack.java`

A stack of ints with a fixed capacity. `push` on a full stack must throw
`IllegalStateException`, `pop`/`top` on an empty stack must throw
`NoSuchElementException` — instead of changing anything — and
`0 <= size() <= capacity` must hold at all times.

### 2. Idempotent Payments — API Guarantees: Idempotency
`src/main/java/day02/PaymentProcessor.java`

Every charge comes with a key, and the same key must never be charged
twice. A repeated key returns the first charge's amount and the unchanged
running total.

### 3. Atomic Record Update — API Guarantees: Atomicity
`src/main/java/day02/RecordStore.java`

A record has two fields, `a` and `b`, and is valid only when `a + b >= 0`.
An invalid update must be rejected as a whole: neither field may change,
even if the key already has a record.

## Challenge

### 4. Transactional Key-Value Store — Atomicity
`src/main/java/day02/TxnKV.java`

A key-value store where a transaction can be opened inside another
transaction. `commit()` moves the innermost transaction's writes and
deletes one level out — not straight into the store; `rollback()` throws
them away. The full semantics are in the class Javadoc.

### 5. Idempotency with Expiry & Conflict
`src/main/java/day02/ExpiringPaymentProcessor.java`

Idempotency keys that are only remembered for a limited time, on a clock
the caller passes in. While a key is remembered, the same amount is a
*replay* and a different amount is a *conflict* — neither charges
anything. When the time runs out, the key is forgotten. The full
semantics are in the class Javadoc.

### 6. Interval Booker — Invariants
`src/main/java/day02/IntervalBooker.java`

A booking calendar over half-open intervals `[start, end)`. Bookings must
never overlap. Also implement `firstFree(from, duration)`: the earliest
`start >= from` where a booking of the given length fits. The full
semantics are in the class Javadoc.
