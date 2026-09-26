"""Day 3 -- Races: the broken counter.

A ticket machine at a bank: every customer who walks in takes the next
number. Many threads take tickets at the same time.

  next()    hand out the next ticket: 1 for the first call, then 2,
            3, ... No two calls may ever get the same number, and no
            number may be skipped.
  issued()  how many tickets have been handed out so far.

Both may be called by many threads at the same time.

Remember the lecture: count += 1 is load -> add -> store, three steps
another thread can slip between. next() has exactly that shape.
"""


class TicketDispenser:
    def __init__(self) -> None:
        # TODO: choose your own representation.
        pass

    def next(self) -> int:
        # TODO
        return 0

    def issued(self) -> int:
        # TODO
        return 0
