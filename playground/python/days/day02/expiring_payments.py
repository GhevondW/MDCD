"""Day 2 -- Challenge: Idempotency with Expiry & Conflict.

Like the idempotent payment processor, but keys are only remembered for
a limited time. There is no real clock: every call receives the current
time `now` as a plain number, and the numbers never go down from call
to call. A key recorded at time T is remembered until time T + ttl;
from now == T + ttl on, it is expired (forgotten).

charge(key, amount, now) -- three cases:
  - The key is unknown, or its record has expired: this is a NEW charge.
    Add `amount` to the total and remember (key, amount) until now + ttl.
  - The key is still remembered and `amount` equals the recorded amount:
    this is a REPLAY (the caller sent the same request twice). Return
    the recorded amount and the current total. Change NOTHING: not the
    total, not the recorded amount, not the expiry time.
  - The key is still remembered but `amount` is different: this is a
    CONFLICT (the same key was used for a different request -- a
    mistake). Return the recorded amount and the current total. Change
    nothing.

total_charged()   sum of all NEW charges.
active_keys(now)  how many recorded keys are not yet expired at `now`.
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
        # TODO: choose your own representation.
        pass

    def charge(self, key: str, amount: int, now: int) -> ChargeOutcome:
        # TODO
        return ChargeOutcome(Status.NEW, 0, 0)

    def total_charged(self) -> int:
        # TODO
        return 0

    def active_keys(self, now: int) -> int:
        # TODO
        return 0
