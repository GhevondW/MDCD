"""Day 2 -- API Guarantees: Idempotency.

Every charge request comes with a key. The same key must never be
charged twice.

charge(key, amount):
  - If this key was never seen before: charge it. Add `amount` to the
    total and return ChargeResult(amount, running_total, was_new=True).
  - If this key was already charged: do NOT charge again. Return the
    amount of the FIRST charge and the current running total, with
    was_new=False. total_charged() must stay the same.

total_charged() returns the sum of all real (first-time) charges.
"""

from dataclasses import dataclass


@dataclass
class ChargeResult:
    amount: int
    running_total: int
    was_new: bool


class PaymentProcessor:
    def __init__(self) -> None:
        # TODO: choose your own representation.
        pass

    def charge(self, key: str, amount: int) -> ChargeResult:
        # TODO
        return ChargeResult(amount=0, running_total=0, was_new=False)

    def total_charged(self) -> int:
        # TODO
        return 0
