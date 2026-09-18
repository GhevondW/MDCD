# Day 2 — Practical Tasks (Python)

Implement the `TODO`s in the `.py` files below; don't change the
`test_*.py` files. Problems 1–3 are the core set; 4–6 are challenge
problems.

Run just this day's tests from `playground/python`:
```
python3 -m unittest discover -s days/day02 -t . -p "test_*.py" -v
```
(`-t .` matters — without it, the relative imports in the test files fail.)

## Core

### 1. Bounded Stack — Contracts & Invariants
`bounded_stack.py`

`BoundedStack` must enforce its own contract instead of silently
corrupting state:
- `push(v)` — precondition: not full. If already full, raise
  `IndexError` instead of pushing.
- `pop()` / `top()` — precondition: not empty. Same idea.
- Invariant to hold between calls: `0 <= size() <= capacity`.

### 2. Idempotent Payments — API Guarantees: Idempotency
`idempotent_payments.py`

The same idempotency key must never be charged twice. `charge()` on a
repeat key must return the *original* charge's amount and the *current*
running total (unchanged) with `was_new = False` — not charge again.

### 3. Atomic Record Update — API Guarantees: Atomicity
`atomic_record.py`

A record has two fields, `a` and `b`. `set()` must be all-or-nothing: if
`a + b < 0`, the update is rejected *in full* — neither field may change,
whether the key already existed or not.

## Challenge

### 4. Transactional Key-Value Store — Atomicity
`txn_kv.py`

A string key-value store with **nested** transactions: `begin()`,
`commit()`, `rollback()`. `commit()` merges the innermost transaction's
writes *and* deletes into its parent transaction (or into the store if it
has no parent); `rollback()` discards them. Reads always see the innermost
state. The full semantics are in the module docstring.

### 5. Idempotency with Expiry & Conflict
`expiring_payments.py`

Idempotency keys that expire on a logical clock: within the TTL a repeat
of a key is a *replay* (same amount) or a *conflict* (different amount) —
neither charges; after the TTL the key is forgotten and the next charge is
new. The full semantics are in the module docstring.

### 6. Interval Booker — Invariants
`interval_booker.py`

A booking calendar over half-open intervals `[start, end)` whose invariant
is that recorded bookings never overlap, plus `first_free(from_, duration)` —
the earliest free slot of a given length at or after `from_`. The full
semantics are in the module docstring.
