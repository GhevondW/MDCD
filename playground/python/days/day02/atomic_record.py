"""Day 2 -- API Guarantees: Atomicity.

A record has two fields, a and b. set() must be all-or-nothing: if
a + b < 0 the update is invalid and must be REJECTED IN FULL -- neither
field may change, whether the key already existed or not. Applying one
field before validating the other is exactly the bug this catches.
"""

from dataclasses import dataclass
from typing import Optional


@dataclass
class Record:
    a: int
    b: int


class RecordStore:
    def __init__(self) -> None:
        self._store: dict[str, Record] = {}

    def set(self, key: str, a: int, b: int) -> bool:
        # TODO: if a + b >= 0, commit BOTH fields together and return True.
        # Otherwise leave the store completely unchanged and return False.
        return False

    def get(self, key: str) -> Optional[Record]:
        # TODO: return the record if `key` has been successfully set, or
        # None otherwise.
        return None
