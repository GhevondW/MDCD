"""Day 2 -- API Guarantees: Idempotency.

The SAME idempotency key must never be charged twice: a repeat of a key
you've already seen must return the result of the FIRST charge -- same
amount, unchanged running total -- not charge again.
"""

from dataclasses import dataclass


@dataclass
class ChargeResult:
    amount: int
    running_total: int
    was_new: bool


class PaymentProcessor:
    def __init__(self) -> None:
        self._seen: dict[str, int] = {}
        self._total = 0

    def charge(self, key: str, amount: int) -> ChargeResult:
        # TODO: if `key` was already charged, return the ORIGINAL amount and
        # the CURRENT running total, with was_new=False -- and don't change
        # total_charged(). Otherwise remember it, add it to the total, and
        # return it with was_new=True.
        return ChargeResult(amount=0, running_total=0, was_new=False)

    def total_charged(self) -> int:
        # TODO
        return 0
