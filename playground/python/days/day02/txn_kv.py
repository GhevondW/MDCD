"""Day 2 -- Challenge: Transactional Key-Value Store.

In-memory string->string store with NESTED transactions.

  set(key, value)   store a value (in the innermost open transaction,
                    or directly if none is open)
  get(key)          value currently visible for `key`, or None
  remove(key)       delete `key`; returns whether it was visible
  begin()           open a new transaction (transactions nest)
  commit()          merge the innermost transaction's changes -- writes
                    AND deletes -- into its parent transaction (or into
                    the store if it has no parent). False if no
                    transaction is open.
  rollback()        discard the innermost transaction's changes.
                    False if no transaction is open.

Reads always see the innermost state: a value set (or a key removed)
inside an open transaction is visible to get() immediately.
"""

from typing import Optional


class TxnKV:
    def __init__(self) -> None:
        # TODO: choose your own representation.
        pass

    def set(self, key: str, value: str) -> None:
        # TODO
        pass

    def get(self, key: str) -> Optional[str]:
        # TODO
        return None

    def remove(self, key: str) -> bool:
        # TODO
        return False

    def begin(self) -> None:
        # TODO
        pass

    def commit(self) -> bool:
        # TODO
        return False

    def rollback(self) -> bool:
        # TODO
        return False
