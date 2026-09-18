"""Day 2 -- Challenge: Transactional Key-Value Store.

A key-value store (str -> str) with transactions. A transaction can be
opened inside another transaction (they nest).

  set(key, value)   store a value. If a transaction is open, the change
                    belongs to that transaction; otherwise it is applied
                    directly to the store.
  get(key)          the value visible right now, or None if the key does
                    not exist. Changes made inside an open transaction
                    are visible immediately.
  remove(key)       delete the key. Returns True if the key was visible
                    before the call. A delete made inside a transaction
                    also belongs to that transaction.
  begin()           open a new transaction.
  commit()          take everything the innermost transaction changed --
                    writes AND deletes -- and move it into the
                    transaction one level out (or into the store itself
                    if there is no outer transaction). Returns False if
                    no transaction is open.
  rollback()        throw away everything the innermost transaction
                    changed. Returns False if no transaction is open.
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
