"""Day 3 -- Critical sections: check-then-act, and two locks at once.

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
  transfer(to, amount)
                     move `amount` from this account to `to`, as ONE
                     step: no thread may ever see the money gone from one
                     account but not yet in the other. amount must be > 0
                     and `to` must be a different account; if not, raise
                     ValueError. If this balance < amount, return False
                     and change nothing.
  Account.total(a, b)
                     a's balance plus b's balance, read as ONE step: every
                     transfer between a and b happens either completely
                     before it or completely after it. a and b must be
                     different accounts; if not, raise ValueError.

Everything except the constructor may be called by many threads at the
same time.

    if self._balance >= amount:
        self._balance -= amount

is a check and an act -- another thread can slip in between them.

transfer and total need two accounts locked at once. Two threads that
lock the same two accounts in opposite orders wait for each other
forever -- the slides' "Two transfers, two locks". Hint: if every thread
locks any two accounts in the same order (say, by a number each account
gets when it is created), that cannot happen. And a method that already
holds a threading.Lock must not call another method that takes the same
lock: a Lock locked twice by the same thread waits for itself forever.
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

    def transfer(self, to: "Account", amount: int) -> bool:
        # TODO
        return False

    @staticmethod
    def total(a: "Account", b: "Account") -> int:
        # TODO
        return 0
