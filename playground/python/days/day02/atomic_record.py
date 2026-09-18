"""Day 2 -- API Guarantees: Atomicity.

A record has two fields, a and b. A record is valid only when
a + b >= 0.

set(key, a, b) must be all-or-nothing:
  - If a + b >= 0: store both fields together and return True.
  - If a + b < 0: the update is invalid. Reject the WHOLE update -- do
    not change either field -- and return False. This also holds when
    the key already has a record: the old record must stay exactly as
    it was.

get(key) returns the stored record, or None if this key was never
successfully set.
"""

from dataclasses import dataclass
from typing import Optional


@dataclass
class Record:
    a: int
    b: int


class RecordStore:
    def __init__(self) -> None:
        # TODO: choose your own representation.
        pass

    def set(self, key: str, a: int, b: int) -> bool:
        # TODO
        return False

    def get(self, key: str) -> Optional[Record]:
        # TODO
        return None
