"""Day 2 -- Challenge: Idempotency with Expiry & Conflict.

A payment processor whose idempotency keys EXPIRE. Time is a logical
clock: every call takes `now` (non-decreasing across calls); a key
recorded at time T with time-to-live `ttl` is expired once now >= T + ttl.

charge(key, amount, now):
  - key unknown, or its record expired  -> charge: add `amount` to the
    total, record (key, amount) with expiry now + ttl, return NEW.
  - key active, same amount             -> return REPLAY with the
    originally recorded amount and the current total. Nothing changes --
    not the total, not the recorded amount, not the expiry.
  - key active, different amount        -> return CONFLICT with the
    recorded amount and the current total. Nothing changes.

total_charged()   sum of all charges that returned NEW.
active_keys(now)  number of recorded keys not yet expired at `now`.
"""

import enum
from dataclasses import dataclass


class Status(enum.Enum):
    NEW = enum.auto()
    REPLAY = enum.auto()
    CONFLICT = enum.auto()


@dataclass
class ChargeOutcome:
    status: Status
    amount: int  # NEW: charged; REPLAY/CONFLICT: recorded
    running_total: int


class ExpiringPaymentProcessor:
    def __init__(self, ttl: int) -> None:
        self._ttl = ttl
        # TODO: choose your own representation.

    def charge(self, key: str, amount: int, now: int) -> ChargeOutcome:
        # TODO
        return ChargeOutcome(Status.NEW, 0, 0)

    def total_charged(self) -> int:
        # TODO
        return 0

    def active_keys(self, now: int) -> int:
        # TODO
        return 0
