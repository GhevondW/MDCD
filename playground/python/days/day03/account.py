"""Day 3 -- Critical sections: check-then-act.

A bank account shared by many threads -- think of one card used by
several people at the same moment. The balance must never go below
zero.

  Account(initial)   start with `initial`. initial must be >= 0; if
                     not, raise ValueError.
  deposit(amount)    add `amount`. amount must be > 0; if not, raise
                     ValueError and change nothing.
  withdraw(amount)   amount must be > 0; if not, raise ValueError and
                     change nothing.
                     If balance >= amount: subtract it, return True.
                     Otherwise return False and change nothing.
  balance()          the current balance.

All four may be called by many threads at the same time.

    if self._balance >= amount:
        self._balance -= amount

is a check and an act -- another thread can slip in between them.
"""


class Account:
    def __init__(self, initial: int) -> None:
        # TODO: choose your own representation.
        pass

    def deposit(self, amount: int) -> None:
        # TODO
        pass

    def withdraw(self, amount: int) -> bool:
        # TODO
        return False

    def balance(self) -> int:
        # TODO
        return 0
