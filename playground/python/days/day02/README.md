# Day 2 — Practical Tasks (Python)

Implement the `TODO`s in the `.py` files below; don't change the
`test_*.py` files. Problems 1–3 are the core set; 4–6 are challenge
problems. Each module's docstring starts with the full problem statement.

Run just this day's tests from `playground/python`:
```
python3 -m unittest discover -s days/day02 -t . -p "test_*.py" -v
```
(`-t .` matters — without it, the relative imports in the test files fail.)

## Core

### 1. Bounded Stack — Contracts & Invariants
`bounded_stack.py`

A stack of ints with a fixed capacity. `push` on a full stack and
`pop`/`top` on an empty stack must raise `IndexError` instead of changing
anything, and `0 <= size() <= capacity` must hold at all times.

### 2. Idempotent Payments — API Guarantees: Idempotency
`idempotent_payments.py`

Every charge comes with a key, and the same key must never be charged
twice. A repeated key returns the first charge's amount and the unchanged
running total.

### 3. Atomic Record Update — API Guarantees: Atomicity
`atomic_record.py`

A record has two fields, `a` and `b`, and is valid only when `a + b >= 0`.
An invalid update must be rejected as a whole: neither field may change,
even if the key already has a record.

## Challenge

### 4. Transactional Key-Value Store — Atomicity
`txn_kv.py`

A key-value store where a transaction can be opened inside another
transaction. `commit()` moves the innermost transaction's writes and
deletes one level out — not straight into the store; `rollback()` throws
them away. The full semantics are in the module docstring.

### 5. Idempotency with Expiry & Conflict
`expiring_payments.py`

Idempotency keys that are only remembered for a limited time, on a clock
the caller passes in. While a key is remembered, the same amount is a
*replay* and a different amount is a *conflict* — neither charges
anything. When the time runs out, the key is forgotten. The full
semantics are in the module docstring.

### 6. Interval Booker — Invariants
`interval_booker.py`

A booking calendar over half-open intervals `[start, end)`. Bookings must
never overlap. Also implement `first_free(from_, duration)`: the earliest
`start >= from_` where a booking of the given length fits. The full
semantics are in the module docstring.
