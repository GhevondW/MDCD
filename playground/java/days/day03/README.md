# Day 3 — Practical Tasks (Java)

Implement the `TODO`s in `src/main/java/day03/`; don't change
`src/test/`. Problems 1–3 are the core set; 4–6 are challenge problems.
Each class's Javadoc starts with the full problem statement.

Every class here is used by many threads at the same time. The tests start
real threads and hammer your code. A test that passes once doesn't prove
the code is race-free; a test that fails even once proves it isn't.

Run just this day's tests from `playground/java`:
```
mvn -pl days/day03 test
```

Each test has a 60-second limit, so a deadlock fails the test instead of
hanging it.

Java has no built-in race detector (like C++'s ThreadSanitizer) for plain
tests like these. So a green run is only evidence, not proof: run the
tests several times, and check your reasoning about every shared field. To
test concurrent code seriously, OpenJDK has a dedicated tool,
[jcstress](https://github.com/openjdk/jcstress).

## Core

### 1. Ticket Dispenser — Races
`src/main/java/day03/TicketDispenser.java`

The lecture's broken counter. `next()` hands out 1, 2, 3, ... — never the
same number twice and never skipping one, no matter how many threads call
it at once.

### 2. Bank Account — Check-then-act, two locks at once
`src/main/java/day03/Account.java`

`withdraw` succeeds only if the balance covers it, so the balance never
goes below zero. The check and the subtraction must happen as one step.
`transfer` moves money between two accounts and `total` reads two
balances, each as one step — so both lock two accounts at once, and two
opposite transfers must never deadlock.

### 3. Bounded Stack, shared — Invariants and API races
`src/main/java/day03/BoundedStack.java`

Day 2's bounded stack with the same contract, now used by many threads
at once. No value may be lost, returned twice, or made up. It also gets
`tryPush` and `tryPop`, which check and act in one call — the fix for
the API race in `if (!s.isEmpty()) s.pop();`.

## Challenge

### 4. Compute Once per Key — Singleton, fixed
`src/main/java/day03/OnceCache.java`

A cache that computes each key's value at most once, even when many
threads ask for the same missing key at the same moment. Different keys
must be able to compute at the same time, so one big lock won't do. The
full semantics are in the class Javadoc.

### 5. Thread-Safe Event Bus — Observer, fixed
`src/main/java/day03/EventBus.java`

Day 1's publisher, safe for many threads and for callbacks that
subscribe, unsubscribe or publish from inside a callback. Each publish
delivers to a snapshot of the subscribers, and no lock may be held while
a callback runs. The full semantics are in the class Javadoc.

### 6. Interval Booker, shared — Check-then-act
`src/main/java/day03/ConcurrentBooker.java`

Day 2's interval booker for many threads. `bookFirstFree` finds and
books a slot in one step, and `bookAll` books a whole group of intervals
or none of them. The full semantics are in the class Javadoc.
